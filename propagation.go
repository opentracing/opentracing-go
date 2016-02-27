package opentracing

import "errors"

///////////////////////////////////////////////////////////////////////////////
// CORE PROPAGATION INTERFACES:
///////////////////////////////////////////////////////////////////////////////

var (
	// ErrUnsupportedFormat occurs when the `format` passed to Tracer.Inject() or
	// Tracer.Join() is not recognized by the Tracer implementation.
	ErrUnsupportedFormat = errors.New("opentracing: Unknown or unsupported Inject/Join format")

	// ErrTraceNotFound occurs when the `carrier` passed to Tracer.Join() is
	// valid and uncorrupted but has insufficient information to join or resume
	// a trace.
	ErrTraceNotFound = errors.New("opentracing: Trace not found in Join carrier")

	// ErrInvalidSpan errors occur when Tracer.Inject() is asked to operate on
	// a Span which it is not prepared to handle (for example, since it was
	// created by a different tracer implementation).
	ErrInvalidSpan = errors.New("opentracing: Span type incompatible with tracer")

	// ErrInvalidCarrier errors occur when Tracer.Inject() or Tracer.Join()
	// implementations expect a different type of `carrier` than they are
	// given.
	ErrInvalidCarrier = errors.New("opentracing: Invalid Inject/Join carrier")

	// ErrTraceCorrupted occurs when the `carrier` passed to Tracer.Join() is
	// of the expected type but is corrupted.
	ErrTraceCorrupted = errors.New("opentracing: Trace data corrupted in Join carrier")
)

///////////////////////////////////////////////////////////////////////////////
// BUILTIN PROPAGATION FORMATS:
///////////////////////////////////////////////////////////////////////////////

// BuiltinFormat is used to demarcate the values within package `opentracing`
// that are intended for use with the Tracer.Inject() and Tracer.Join()
// methods.
type BuiltinFormat byte

const (
	// SplitBinary encodes the Span in a SplitBinaryCarrier instance.
	//
	// The `carrier` for injection and extraction must be a
	// `*SplitBinaryCarrier` instance.
	SplitBinary BuiltinFormat = iota

	// SplitText encodes the Span in a SplitTextCarrier instance.
	//
	// The `carrier` for injection and extraction must be a `*SplitTextCarrier`
	// instance.
	SplitText

	// GoHTTPHeader encodes the Span into a Go http.Header instance (both the
	// tracer state and any baggage).
	//
	// The `carrier` for both injection and extraction must be an http.Header
	// instance.
	GoHTTPHeader
)

// SplitTextCarrier breaks a propagated Span into two pieces.
//
// The Span is separated in this way for a variety of reasons; the most
// important is to give OpenTracing users a portable way to opt out of
// Baggage propagation entirely if they deem it a stability risk.
//
// It is legal to provide one or both maps as `nil`; they will be created
// as needed. If non-nil maps are provided, they will be used without
// clearing them out on injection.
type SplitTextCarrier struct {
	// TracerState is Tracer-specific context that must cross process
	// boundaries. For example, in Dapper this would include a trace_id, a
	// span_id, and a bitmask representing the sampling status for the given
	// trace.
	TracerState map[string]string

	// Any Baggage for the encoded Span (per Span.SetBaggageItem).
	Baggage map[string]string
}

// NewSplitTextCarrier creates a new SplitTextCarrier.
func NewSplitTextCarrier() *SplitTextCarrier {
	return &SplitTextCarrier{}
}

// SplitBinaryCarrier breaks a propagated Span into two pieces.
//
// The Span is separated in this way for a variety of reasons; the most
// important is to give OpenTracing users a portable way to opt out of
// Baggage propagation entirely if they deem it a stability risk.
//
// Both byte slices may be nil; on injection, what is provided will be cleared
// and the resulting capacity used.
type SplitBinaryCarrier struct {
	// TracerState is Tracer-specific context that must cross process
	// boundaries. For example, in Dapper this would include a trace_id, a
	// span_id, and a bitmask representing the sampling status for the given
	// trace.
	TracerState []byte

	// Any Baggage for the encoded Span (per Span.SetBaggageItem).
	Baggage []byte
}

// NewSplitBinaryCarrier creates a new SplitTextCarrier.
func NewSplitBinaryCarrier() *SplitBinaryCarrier {
	return &SplitBinaryCarrier{}
}
