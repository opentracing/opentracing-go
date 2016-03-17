package opentracing

import (
	"errors"
	"io"
)

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

// TextMapCarrier represents a Span for propagation within a key:value map of
// unicode strings.
type TextMapCarrier interface {
	// Add a key:value pair to the carrier.`key` should be prefixed in a way
	// that protects against collisions with things like HTTP headers (which
	// may share space with Tracer data in the TextMapCarrier).
	Add(key, val string)

	// GetAll returns all contents of the text map via repeated calls to the
	// `handler` function. If any call to `handler` returns a non-nil error,
	// GetAll terminates and returns that error.
	//
	// NOTE: `handler` may be invoked for key:value combinations that were
	// *not* added via `Set` (in this process or otherwise). As such,
	// implementations MUST check that `key` is something they care about.
	//
	// NOTE: A single `key` may appear in multiple calls to `handler` for a
	// single `GetAll` invocation.
	GetAll(handler func(key, val string) error) error
}

// BinaryCarrier represents a Span for propagation as an opaque byte stream.
//
// The io.Writer is intended for use with Tracer.Inject() and the io.Reader is
// intended for use with Tracer.Join().
type BinaryCarrier interface {
	io.ReadWriter
}

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
