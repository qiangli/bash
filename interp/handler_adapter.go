// Copyright (c) 2026, the bashy authors
// See LICENSE for licensing information

package interp

import (
	"context"
	"fmt"
	"os"
)

// HandlerKind names which handler a [HandlerContext] was built for. It maps
// explicitly onto the private handlerKind enum: Exec starts at 1, matching
// handlerKindExec, so an owned engine can label adapted contexts without
// reaching into private state.
type HandlerKind int

const (
	// HandlerKindExec labels contexts handed to an [ExecHandlerFunc].
	HandlerKindExec HandlerKind = iota + 1
	// HandlerKindCall labels contexts handed to a [CallHandlerFunc].
	HandlerKindCall
	// HandlerKindOpen labels contexts handed to an [OpenHandlerFunc].
	HandlerKindOpen
	// HandlerKindReadDir labels contexts handed to a [ReadDirHandlerFunc2].
	HandlerKindReadDir
)

// HandlerHooks supplies the live callbacks behind an adapted
// [HandlerContext] built by [AdaptHandlerContext]. Every hook is optional:
// a nil hook reports a safe neutral result (false, zero, or an error for
// Builtin) instead of consulting the Runner, so an owned engine gets fully
// deterministic views.
type HandlerHooks struct {
	// DryRun reports whether the emulated dryrun option is on.
	DryRun func() bool
	// Async reports whether the emulated command runs asynchronously.
	Async func() bool
	// Umask reports the emulated virtual file-creation mask.
	Umask func() os.FileMode
	// SignalIgnored reports whether the named signal is ignored.
	SignalIgnored func(string) bool
	// Builtin runs a builtin; only consulted for exec-kind contexts.
	Builtin func(context.Context, []string) error
	// DescribeValue reports the named session value, read-only.
	DescribeValue func(string) (ValueDescription, bool)
}

// WithHandlerContext returns a context carrying hc, retrievable via
// [HandlerCtx]. It lets an owned engine synthesize handler contexts without
// a Runner; ordinary Runner-built contexts keep travelling the same key.
func WithHandlerContext(ctx context.Context, hc HandlerContext) context.Context {
	return context.WithValue(ctx, handlerCtxKey{}, hc)
}

// AdaptHandlerContext returns a copy of hc labelled with kind and backed by
// hooks. The adapted view calls the copied hook functions on every access,
// so state captured by the callbacks (for example variables closed over by
// a closure) is observed live without re-adapting. The HandlerHooks value
// itself is copied at adapt time: reassigning fields on the caller's struct
// afterwards has no effect; only the function values held by the copy are
// consulted. Existing public fields are preserved untouched; only the
// private kind is relabelled and the hooks attached. An out-of-range kind
// maps to the invalid zero value, which keeps the exec-only
// [HandlerContext.Builtin] restriction closed.
func AdaptHandlerContext(hc HandlerContext, kind HandlerKind, hooks HandlerHooks) HandlerContext {
	switch kind {
	case HandlerKindExec:
		hc.kind = handlerKindExec
	case HandlerKindCall:
		hc.kind = handlerKindCall
	case HandlerKindOpen:
		hc.kind = handlerKindOpen
	case HandlerKindReadDir:
		hc.kind = handlerKindReadDir
	default:
		hc.kind = handlerKind(0)
	}
	hc.hooks = &hooks
	return hc
}

// errNoBuiltinHook is the safe neutral result of calling Builtin on an
// adapted exec-kind context whose hooks carry no Builtin function.
var errNoBuiltinHook = fmt.Errorf("HandlerContext.Builtin has no adapted builtin hook")
