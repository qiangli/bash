package interp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/qiangli/bash/expand"
	"github.com/qiangli/bash/syntax"
)

// TestHandlerKindValues pins the explicit public-to-private kind mapping:
// Exec starts at 1 and each constant lands on its private enum peer.
func TestHandlerKindValues(t *testing.T) {
	t.Parallel()

	if HandlerKindExec != 1 {
		t.Fatalf("HandlerKindExec = %d, want 1", int(HandlerKindExec))
	}
	cases := []struct {
		public HandlerKind
		want   handlerKind
	}{
		{HandlerKindExec, handlerKindExec},
		{HandlerKindCall, handlerKindCall},
		{HandlerKindOpen, handlerKindOpen},
		{HandlerKindReadDir, handlerKindReadDir},
	}
	for _, tc := range cases {
		got := AdaptHandlerContext(HandlerContext{}, tc.public, HandlerHooks{})
		if got.kind != tc.want {
			t.Errorf("AdaptHandlerContext kind %d mapped to %d, want %d",
				int(tc.public), int(got.kind), int(tc.want))
		}
	}
	if got := AdaptHandlerContext(HandlerContext{}, HandlerKind(99), HandlerHooks{}); got.kind != handlerKind(0) {
		t.Errorf("out-of-range kind mapped to %d, want invalid 0", int(got.kind))
	}
}

// TestAdaptHandlerHooksForwarding verifies each of the six hooks is
// consulted with the caller's arguments.
func TestAdaptHandlerHooksForwarding(t *testing.T) {
	t.Parallel()

	var gotSignal, gotDescribe string
	var gotBuiltinArgs []string
	ctx := context.Background()
	hc := AdaptHandlerContext(HandlerContext{}, HandlerKindExec, HandlerHooks{
		DryRun:        func() bool { return true },
		Async:         func() bool { return true },
		Umask:         func() os.FileMode { return 0o027 },
		SignalIgnored: func(s string) bool { gotSignal = s; return s == "PIPE" },
		Builtin:       func(_ context.Context, args []string) error { gotBuiltinArgs = args; return nil },
		DescribeValue: func(s string) (ValueDescription, bool) {
			gotDescribe = s
			return ValueDescription{Name: s, Type: "string", Kind: "scalar", Value: "v"}, true
		},
	})

	if !hc.DryRun() {
		t.Error("DryRun did not forward hook result")
	}
	if !hc.Async() {
		t.Error("Async did not forward hook result")
	}
	if got := hc.Umask(); got != 0o027 {
		t.Errorf("Umask = %v, want 027", got)
	}
	if !hc.SignalIgnored("PIPE") || gotSignal != "PIPE" {
		t.Errorf("SignalIgnored did not forward argument, got %q", gotSignal)
	}
	if err := hc.Builtin(ctx, []string{"echo", "hi"}); err != nil {
		t.Errorf("Builtin returned error: %v", err)
	} else if len(gotBuiltinArgs) != 2 || gotBuiltinArgs[0] != "echo" {
		t.Errorf("Builtin did not forward args, got %q", gotBuiltinArgs)
	}
	d, ok := hc.DescribeValue("probe")
	if !ok || gotDescribe != "probe" || d.Value != "v" {
		t.Errorf("DescribeValue did not forward, got %+v,%v arg %q", d, ok, gotDescribe)
	}
}

// TestAdaptHandlerHooksLive verifies adapted views consult the hooks on
// every call rather than snapshotting results at adapt time.
func TestAdaptHandlerHooksLive(t *testing.T) {
	t.Parallel()

	var dryCalls, asyncCalls, umaskCalls, signalCalls, describeCalls int
	enabled := false
	hc := AdaptHandlerContext(HandlerContext{}, HandlerKindExec, HandlerHooks{
		DryRun: func() bool { dryCalls++; return enabled },
		Async:  func() bool { asyncCalls++; return !enabled },
		Umask: func() os.FileMode {
			umaskCalls++
			if enabled {
				return 0o077
			}
			return 0o022
		},
		SignalIgnored: func(string) bool { signalCalls++; return enabled },
		DescribeValue: func(string) (ValueDescription, bool) {
			describeCalls++
			return ValueDescription{}, enabled
		},
	})

	if hc.DryRun() || !hc.Async() {
		t.Fatal("hooks did not report initial state")
	}
	enabled = true
	if !hc.DryRun() || hc.Async() {
		t.Fatal("hooks did not observe changed state without re-adapting")
	}
	if got := hc.Umask(); got != 0o077 {
		t.Fatalf("Umask = %v, want 077 after change", got)
	}
	if !hc.SignalIgnored("PIPE") {
		t.Fatal("SignalIgnored did not observe changed state")
	}
	if _, ok := hc.DescribeValue("x"); !ok {
		t.Fatal("DescribeValue did not observe changed state")
	}
	if dryCalls != 2 || asyncCalls != 2 || umaskCalls != 1 || signalCalls != 1 || describeCalls != 1 {
		t.Errorf("hook call counts = %d/%d/%d/%d/%d, want 2/2/1/1/1",
			dryCalls, asyncCalls, umaskCalls, signalCalls, describeCalls)
	}
}

// TestAdaptHandlerNilHooksNeutral verifies adapted views with nil hooks
// report safe neutral results instead of consulting the Runner. The base is
// a real Runner-backed exec context with observable non-neutral state (a
// umask set by the umask builtin, a known shell variable, and a working
// exec-kind Builtin), so an accidental fallback to the Runner would report
// non-neutral values and fail the test.
func TestAdaptHandlerNilHooksNeutral(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r, err := New(StdIO(nil, io.Discard, io.Discard))
	if err != nil {
		t.Fatal(err)
	}
	parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
	file, err := parser.Parse(strings.NewReader("umask 027; adapter_probe_set=present\n"), "adapter.sh")
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Run(ctx, file); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	base := HandlerCtx(r.handlerCtx(ctx, handlerKindExec, syntax.Pos{}))
	if got := base.Umask(); got != 0o027 {
		t.Fatalf("base Umask = %v, want 027 (Runner-backed sanity check)", got)
	}
	if d, ok := base.DescribeValue("adapter_probe_set"); !ok || d.Value != "present" {
		t.Fatalf("base DescribeValue = %+v,%v, want present,true (Runner-backed sanity check)", d, ok)
	}
	if err := base.Builtin(ctx, []string{"true"}); err != nil {
		t.Fatalf("base Builtin(true) err = %v, want nil (Runner-backed sanity check)", err)
	}
	hc := AdaptHandlerContext(base, HandlerKindExec, HandlerHooks{})
	if hc.DryRun() {
		t.Error("nil DryRun hook = true, want false")
	}
	if hc.Async() {
		t.Error("nil Async hook = true, want false")
	}
	if got := hc.Umask(); got != 0 {
		t.Errorf("nil Umask hook = %v, want 0", got)
	}
	if hc.SignalIgnored("PIPE") {
		t.Error("nil SignalIgnored hook = true, want false")
	}
	if err := hc.Builtin(context.Background(), []string{"true"}); !errors.Is(err, errNoBuiltinHook) {
		t.Errorf("nil Builtin hook err = %v, want errNoBuiltinHook", err)
	}
	if d, ok := hc.DescribeValue("adapter_probe_set"); ok || !reflect.DeepEqual(d, ValueDescription{}) {
		t.Errorf("nil DescribeValue hook = %+v,%v, want zero,false", d, ok)
	}
}

// TestAdaptHandlerBuiltinKindRestriction verifies the exec-only Builtin
// restriction survives adaptation: non-exec kinds reject even with a hook,
// and an out-of-range kind stays closed.
func TestAdaptHandlerBuiltinKindRestriction(t *testing.T) {
	t.Parallel()

	called := false
	hooks := HandlerHooks{
		Builtin: func(context.Context, []string) error { called = true; return nil },
	}
	for _, kind := range []HandlerKind{
		HandlerKindCall, HandlerKindOpen, HandlerKindReadDir, HandlerKind(99),
	} {
		called = false
		hc := AdaptHandlerContext(HandlerContext{}, kind, hooks)
		err := hc.Builtin(context.Background(), []string{"true"})
		if err == nil || !strings.Contains(err.Error(), "only be called via an ExecHandlerFunc") {
			t.Errorf("kind %d Builtin err = %v, want exec-only restriction", int(kind), err)
		}
		if called {
			t.Errorf("kind %d Builtin invoked the hook, want rejection", int(kind))
		}
	}
	hc := AdaptHandlerContext(HandlerContext{}, HandlerKindExec, hooks)
	if err := hc.Builtin(context.Background(), []string{"true"}); err != nil {
		t.Errorf("exec-kind Builtin err = %v, want nil", err)
	}
	if !called {
		t.Error("exec-kind Builtin did not invoke the hook")
	}
}

// TestWithHandlerContextRoundtrip verifies every preserved public field
// survives a store/retrieve cycle, including exec flags, stdio, env, and an
// invalid position.
func TestWithHandlerContextRoundtrip(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	base := HandlerContext{
		Agentic:      true,
		Env:          expand.ListEnviron("A=1", "B=two"),
		Dir:          "/probe",
		Pos:          syntax.Pos{},
		Stdin:        strings.NewReader("in"),
		Stdout:       &stdout,
		Stderr:       &stderr,
		ExecAs:       "argv0",
		ExecClearEnv: true,
		ExecReplace:  true,
	}
	got := HandlerCtx(WithHandlerContext(context.Background(), base))
	if got.Agentic != base.Agentic {
		t.Error("Agentic did not roundtrip")
	}
	if got.Env.Get("A").Str != "1" || got.Env.Get("B").Str != "two" {
		t.Error("Env did not roundtrip")
	}
	if got.Dir != base.Dir {
		t.Errorf("Dir = %q, want %q", got.Dir, base.Dir)
	}
	if got.Pos != (syntax.Pos{}) {
		t.Errorf("invalid Pos did not roundtrip, got %v", got.Pos)
	}
	if got.Stdin != base.Stdin {
		t.Error("Stdin did not roundtrip")
	}
	if got.Stdout != base.Stdout || got.Stderr != base.Stderr {
		t.Error("stdio did not roundtrip")
	}
	if got.ExecAs != "argv0" || !got.ExecClearEnv || !got.ExecReplace {
		t.Errorf("exec flags = %q/%v/%v, want argv0/true/true",
			got.ExecAs, got.ExecClearEnv, got.ExecReplace)
	}
	if got.hooks != nil {
		t.Error("unadapted base roundtripped as adapted")
	}

	// An adapted view keeps its live hooks across the roundtrip.
	adapted := AdaptHandlerContext(base, HandlerKindExec, HandlerHooks{
		DryRun: func() bool { return true },
	})
	if back := HandlerCtx(WithHandlerContext(context.Background(), adapted)); !back.DryRun() {
		t.Error("adapted hooks did not survive the roundtrip")
	}
}

// TestOrdinaryRunnerFallback verifies contexts built by the ordinary Runner
// path stay unadapted and behave as before, using literal expectations only.
func TestOrdinaryRunnerFallback(t *testing.T) {
	t.Parallel()

	var captured HandlerContext
	var buf bytes.Buffer
	r, err := New(
		Dir(t.TempDir()),
		StdIO(nil, &buf, &buf),
		CallHandler(func(ctx context.Context, args []string) ([]string, error) {
			captured = HandlerCtx(ctx)
			return args, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	parser := syntax.NewParser(syntax.Variant(syntax.LangBash))
	file, err := parser.Parse(strings.NewReader(":\n"), "adapter.sh")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.Run(ctx, file); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if captured.hooks != nil {
		t.Fatal("ordinary Runner context carries hooks")
	}
	if captured.kind != handlerKindCall {
		t.Errorf("ordinary call context kind = %d, want %d", int(captured.kind), int(handlerKindCall))
	}
	if captured.DryRun() {
		t.Error("ordinary DryRun = true, want false")
	}
	if captured.Async() {
		t.Error("ordinary Async = true, want false")
	}
	if captured.SignalIgnored("BOGUS-SIGNAL") {
		t.Error("ordinary SignalIgnored(unknown) = true, want false")
	}
	if captured.SignalIgnored("KILL") {
		t.Error("ordinary SignalIgnored(KILL) = true, want false")
	}
	if _, ok := captured.DescribeValue("definitely_unset_adapter_probe"); ok {
		t.Error("ordinary DescribeValue(unset) = true, want false")
	}
	if err := captured.Builtin(ctx, []string{"true"}); err == nil ||
		!strings.Contains(err.Error(), "only be called via an ExecHandlerFunc") {
		t.Errorf("ordinary call-kind Builtin err = %v, want exec-only restriction", err)
	}

	// The ordinary exec-kind Builtin path still reaches the Runner.
	execCtx := r.handlerCtx(context.Background(), handlerKindExec, syntax.Pos{})
	if err := HandlerCtx(execCtx).Builtin(ctx, []string{"true"}); err != nil {
		t.Errorf("ordinary exec-kind Builtin(true) err = %v, want nil", err)
	}
}

// TestZeroHandlerContextNeutral verifies a zero HandlerContext keeps the
// long-standing nil-runner neutral behavior.
func TestZeroHandlerContextNeutral(t *testing.T) {
	t.Parallel()

	var hc HandlerContext
	if hc.DryRun() || hc.Async() {
		t.Error("zero HandlerContext reports DryRun/Async, want false")
	}
	if got := hc.Umask(); got != 0 {
		t.Errorf("zero HandlerContext Umask = %v, want 0", got)
	}
	if hc.SignalIgnored("PIPE") {
		t.Error("zero HandlerContext SignalIgnored = true, want false")
	}
	if _, ok := hc.DescribeValue("x"); ok {
		t.Error("zero HandlerContext DescribeValue = true, want false")
	}
	if err := hc.Builtin(context.Background(), []string{"true"}); err == nil {
		t.Error("zero HandlerContext Builtin = nil, want exec-only restriction error")
	}
}
