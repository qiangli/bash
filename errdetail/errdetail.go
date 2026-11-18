// Package errdetail owns the neutral structured-error data contract shared
// by the Classic foundation and its consumers. It depends only on the Go
// standard library and never imports coreutils, Yoke, Bashsharp, or any
// language worker implementation: workers keep their private error types
// with their owner and expose their detail through the implicit carrier
// interface below, which ForeignErrorDetail resolves with errors.As.
//
// The JSON contract here is frozen: field names, tags, and UnmarshalJSON
// behavior (legacy plain-string and null frames) are preserved verbatim
// from the previous owner so existing wire frames keep decoding unchanged.
package errdetail

import (
	"bytes"
	"encoding/json"
	"errors"
)

// ErrorDetail is the neutral wire/detail shape for structured foreign worker
// failures. It is deliberately not an error type and it is not tied to any
// command-envelope package: callers that need a domain-specific envelope
// should map this DTO at their own boundary. Workers may still send the legacy
// plain string; decoding normalizes that form into Message while leaving Code
// and Help empty.
type ErrorDetail struct {
	Code    string       `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
	Help    string       `json:"help,omitempty"`
	Cause   *ErrorDetail `json:"cause,omitempty"`
}

// carrier is the implicit contract for errors that report a structured
// detail. It is matched structurally with errors.As, so a worker owner only
// has to implement this one method; no shared base type is required.
type carrier interface {
	ForeignErrorDetail() *ErrorDetail
}

// ForeignErrorDetail returns the structured detail carried by an error. The
// returned detail is the carrier's own deep copy, so callers may read or
// mutate it without affecting the error. It reports false when no error in
// err's chain implements the carrier contract, and also when the carrier
// reports a nil detail, so callers may dereference on ok.
func ForeignErrorDetail(err error) (*ErrorDetail, bool) {
	var carried carrier
	if !errors.As(err, &carried) {
		return nil, false
	}
	detail := carried.ForeignErrorDetail()
	if detail == nil {
		return nil, false
	}
	return detail, true
}

// Clone returns a deep copy of d, or nil when d is nil.
func Clone(d *ErrorDetail) *ErrorDetail {
	if d == nil {
		return nil
	}
	out := &ErrorDetail{Code: d.Code, Message: d.Message, Help: d.Help}
	out.Cause = Clone(d.Cause)
	return out
}

// UnmarshalJSON decodes structured frames, the legacy plain-string frame
// (normalized into Message), and null or blank frames (normalized into the
// zero detail).
func (d *ErrorDetail) UnmarshalJSON(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 || bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*d = ErrorDetail{}
		return nil
	}
	var legacy string
	if err := json.Unmarshal(data, &legacy); err == nil {
		*d = ErrorDetail{Message: legacy}
		return nil
	}
	type errorDetail ErrorDetail
	var structured errorDetail
	if err := json.Unmarshal(data, &structured); err != nil {
		return err
	}
	*d = ErrorDetail(structured)
	return nil
}
