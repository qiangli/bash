package errdetail

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// carrierFunc lets tests build errors that satisfy the implicit carrier
// contract without importing any language worker.
type carrierFunc func() *ErrorDetail

func (f carrierFunc) Error() string { return "carrier" }

func (f carrierFunc) ForeignErrorDetail() *ErrorDetail { return f() }

func mustDecode(t *testing.T, frame string) ErrorDetail {
	t.Helper()
	var detail ErrorDetail
	if err := json.Unmarshal([]byte(frame), &detail); err != nil {
		t.Fatalf("UnmarshalJSON(%q) error = %v", frame, err)
	}
	return detail
}

func TestUnmarshalJSONPreservesWireContract(t *testing.T) {
	tests := []struct {
		name  string
		frame string
		want  ErrorDetail
	}{
		{
			name:  "structured",
			frame: `{"code":"ValueError","message":"boom","help":"trace"}`,
			want:  ErrorDetail{Code: "ValueError", Message: "boom", Help: "trace"},
		},
		{
			name:  "legacy string",
			frame: `"plain boom"`,
			want:  ErrorDetail{Message: "plain boom"},
		},
		{
			name:  "null",
			frame: `null`,
			want:  ErrorDetail{},
		},
		{
			name:  "empty object",
			frame: `{}`,
			want:  ErrorDetail{},
		},
		{
			name:  "nested cause",
			frame: `{"code":"RuntimeError","message":"outer","cause":{"code":"ValueError","message":"inner"}}`,
			want: ErrorDetail{Code: "RuntimeError", Message: "outer",
				Cause: &ErrorDetail{Code: "ValueError", Message: "inner"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustDecode(t, tt.frame); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("UnmarshalJSON(%q) = %#v, want %#v", tt.frame, got, tt.want)
			}
		})
	}
}

// TestUnmarshalJSONBlankFrameDirect exercises the preserved blank-input
// behavior by invoking UnmarshalJSON directly: encoding/json rejects a
// whitespace-only frame before the method ever runs, so going through
// json.Unmarshal cannot cover this decoder branch.
func TestUnmarshalJSONBlankFrameDirect(t *testing.T) {
	for _, frame := range []string{"", "   ", "\n\t "} {
		var detail ErrorDetail
		if err := detail.UnmarshalJSON([]byte(frame)); err != nil {
			t.Fatalf("UnmarshalJSON(%q) error = %v", frame, err)
		}
		if detail != (ErrorDetail{}) {
			t.Fatalf("UnmarshalJSON(%q) = %#v, want zero detail", frame, detail)
		}
	}
}

func TestUnmarshalJSONRejectsInvalid(t *testing.T) {
	var detail ErrorDetail
	if err := json.Unmarshal([]byte(`{"code":`), &detail); err == nil {
		t.Fatalf("UnmarshalJSON(invalid) = nil, want error")
	}
}

func TestMarshalKeepsFieldNames(t *testing.T) {
	out, err := json.Marshal(ErrorDetail{
		Code: "ValueError", Message: "boom", Help: "trace",
		Cause: &ErrorDetail{Code: "KeyError"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"code":"ValueError","message":"boom","help":"trace","cause":{"code":"KeyError"}}`
	if string(out) != want {
		t.Fatalf("Marshal = %s, want %s", out, want)
	}
	roundTrip := mustDecode(t, string(out))
	if roundTrip.Code != "ValueError" || roundTrip.Cause == nil || roundTrip.Cause.Code != "KeyError" {
		t.Fatalf("round trip = %#v", roundTrip)
	}
}

func TestCloneDeepAndIsolated(t *testing.T) {
	if Clone(nil) != nil {
		t.Fatalf("Clone(nil) != nil")
	}
	orig := &ErrorDetail{Code: "RuntimeError", Message: "outer", Help: "trace",
		Cause: &ErrorDetail{Code: "ValueError", Message: "inner"}}
	dup := Clone(orig)
	if !reflect.DeepEqual(dup, orig) {
		t.Fatalf("Clone = %#v, want %#v", dup, orig)
	}
	dup.Message = "changed"
	dup.Cause.Code = "changed"
	if orig.Message != "outer" || orig.Cause.Code != "ValueError" {
		t.Fatalf("Clone aliases original: %#v", orig)
	}
	// The copy must not share cause pointers with the source.
	if dup.Cause == orig.Cause {
		t.Fatalf("Clone shares cause pointer")
	}
	empty := Clone(&ErrorDetail{})
	if empty == nil || *empty != (ErrorDetail{}) {
		t.Fatalf("Clone(empty) = %#v", empty)
	}
}

func TestForeignErrorDetailCarrier(t *testing.T) {
	want := &ErrorDetail{Code: "ValueError", Message: "boom"}
	err := carrierFunc(func() *ErrorDetail { return Clone(want) })
	got, ok := ForeignErrorDetail(err)
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("ForeignErrorDetail(direct) = %#v, %t", got, ok)
	}
	wrapped := fmt.Errorf("island call: %w", fmt.Errorf("retry: %w", err))
	got, ok = ForeignErrorDetail(wrapped)
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("ForeignErrorDetail(wrapped) = %#v, %t", got, ok)
	}
	// The returned detail is a copy: mutating it leaves the error intact.
	got.Message = "mutated"
	if again, ok := ForeignErrorDetail(err); !ok || again.Message != "boom" {
		t.Fatalf("carrier not isolated: %#v, %t", again, ok)
	}
}

func TestForeignErrorDetailNilDetail(t *testing.T) {
	// A carrier reporting a nil detail must resolve to (nil, false) so
	// existing polyglot callers can dereference on ok.
	err := carrierFunc(func() *ErrorDetail { return nil })
	if got, ok := ForeignErrorDetail(err); ok || got != nil {
		t.Fatalf("ForeignErrorDetail(nil-detail) = %#v, %t; want (nil, false)", got, ok)
	}
	wrapped := fmt.Errorf("island call: %w", err)
	if got, ok := ForeignErrorDetail(wrapped); ok || got != nil {
		t.Fatalf("ForeignErrorDetail(wrapped nil-detail) = %#v, %t; want (nil, false)", got, ok)
	}
}

func TestForeignErrorDetailRejectsOrdinary(t *testing.T) {
	for name, err := range map[string]error{
		"nil":       nil,
		"plain":     errors.New("boom"),
		"wrapped":   fmt.Errorf("ctx: %w", errors.New("boom")),
		"transport": fmt.Errorf("invalid worker response: %s", "unexpected end"),
	} {
		if got, ok := ForeignErrorDetail(err); ok || got != nil {
			t.Fatalf("ForeignErrorDetail(%s) = %#v, %t; want (nil, false)", name, got, ok)
		}
	}
}
