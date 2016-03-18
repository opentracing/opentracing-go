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
	// Binary encodes the Span in a BinaryCarrier instance.
	//
	// The `carrier` for injection and extraction must be a `BinaryCarrier`
	// instance (not a pointer to a BinaryCarrier instance).
	Binary BuiltinFormat = iota

	// TextMap encodes the Span in a TextMapCarrier instance.
	//
	// The `carrier` for injection and extraction must be a `TextMapCarrier`
	// instance (not a pointer to a TextMapCarrier instance).
	TextMap

	// GoHTTPHeader encodes the Span into a Go http.Header instance (both the
	// tracer state and any baggage).
	//
	// The `carrier` for both injection and extraction must be an http.Header
	// instance (not a pointer to an http.Header instance).
	//
	// If there are entries in the http.Header map prior to a call to Inject(),
	// they are left alone (i.e., the map is not cleared). Similarly, in calls
	// to Join() it is fine (and expected in some cases) for the http.Header
	// map to contain other unrelated data (i.e., non-OpenTracing headers).
	GoHTTPHeader

	// SplitBinary is DEPRECATED
	SplitBinary
	// SplitText is DEPRECATED
	SplitText
)

// TextMapCarrier represents a Span for propagation as a key:value map of
// unicode strings.
//
// If there are entries in the TextMapCarrier prior to a call to Inject(), they
// are left alone (i.e., the map is not cleared). Similarly, in calls to Join()
// it is fine (and expected in some cases) for the TextMapCarrier to contain
// other unrelated data (e.g., arbitrary HTTP header pairs).
type TextMapCarrier map[string]string

// BinaryCarrier represents a Span for propagation as an opaque byte array.
//
// It is fine to pass `nil` to Inject(); in that case, it will allocate a
// []byte for the caller.
type BinaryCarrier *[]byte

// SplitTextCarrier is DEPRECATED
type SplitTextCarrier struct {
	TracerState map[string]string
	Baggage     map[string]string
}

// NewSplitTextCarrier is DEPRECATED
func NewSplitTextCarrier() *SplitTextCarrier {
	return &SplitTextCarrier{}
}

// SplitBinaryCarrier is DEPRECATED
type SplitBinaryCarrier struct {
	TracerState []byte
	Baggage     []byte
}

// NewSplitBinaryCarrier is DEPRECATED
func NewSplitBinaryCarrier() *SplitBinaryCarrier {
	return &SplitBinaryCarrier{}
}
