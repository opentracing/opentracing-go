package opentracing

import (
	"errors"
	"net/http"
	"net/url"
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
	// Binary encodes the Span for propagation as opaque binary data.
	//
	// For Tracer.Inject(): the carrier must be an `io.Writer`.
	//
	// For Tracer.Join(): the carrier must be an `io.Reader`.
	Binary BuiltinFormat = iota

	// TextMap encodes the Span as key:value pairs.
	//
	// For Tracer.Inject(): the carrier must be a `TextMapWriter`.
	//
	// For Tracer.Join(): the carrier must be a `TextMapReader`.
	//
	// See HTTPHeaderTextMapCarrier for an implementation of both TextMapWriter
	// and TextMapReader that defers to an http.Header instance for storage.
	// For example, Inject():
	//
	//    carrier := HTTPHeaderTextMapCarrier(httpReq.Header)
	//    err := span.Tracer().Inject(span, TextMap, carrier)
	//
	// Or Join():
	//
	//    carrier := HTTPHeaderTextMapCarrier(httpReq.Header)
	//    span, err := tracer.Join("opName", TextMap, carrier)
	//
	TextMap
)

// TextMapWriter is the Inject() carrier for the TextMap builtin format. With
// it, the caller can encode a Span for propagation as entries in a multimap of
// unicode strings.
type TextMapWriter interface {
	// Set a key:value pair to the carrier. Multiple calls to Set() for the
	// same key leads to undefined behavior.
	//
	// NOTE: Since HTTP headers are a particularly important use case for the
	// TextMap carrier, `key` parameters identify their respective values in a
	// case-insensitive manner.
	//
	// NOTE: The backing store for the TextMapWriter may contain unrelated data
	// (e.g., arbitrary HTTP headers). As such, the TextMap writer and reader
	// should agree on a prefix or other convention to distinguish their
	// key:value pairs.
	Set(key, val string)
}

// TextMapReader is the Join() carrier for the TextMap builtin format. With it,
// the caller can decode a propagated Span as entries in a multimap of unicode
// strings.
type TextMapReader interface {
	// ForeachKey returns TextMap contents via repeated calls to the `handler`
	// function. If any call to `handler` returns a non-nil error, ForeachKey
	// terminates and returns that error.
	//
	// NOTE: A single `key` may appear in multiple calls to `handler` for a
	// single `ForeachKey` invocation.
	//
	// NOTE: The ForeachKey handler *may* be invoked for keys not set by any
	// TextMap writer (e.g., totally unrelated HTTP headers). As such, the
	// TextMap writer and reader should agree on a prefix or other convention
	// to distinguish their key:value pairs.
	//
	// The "foreach" callback pattern reduces unnecessary copying in some cases
	// and also allows implementations to hold locks while the map is read.
	ForeachKey(handler func(key, val string) error) error
}

// HTTPHeaderTextMapCarrier satisfies both TextMapWriter and TextMapReader.
//
type HTTPHeaderTextMapCarrier http.Header

// Set conforms to the TextMapWriter interface.
func (c HTTPHeaderTextMapCarrier) Set(key, val string) {
	h := http.Header(c)
	h.Add(key, url.QueryEscape(val))
}

// ForeachKey conforms to the TextMapReader interface.
func (c HTTPHeaderTextMapCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range c {
		for _, v := range vals {
			rawV, err := url.QueryUnescape(v)
			if err != nil {
				// We don't know if there was an error escaping an
				// OpenTracing-related header or something else; as such, we
				// continue rather than return the error.
				continue
			}
			if err = handler(k, rawV); err != nil {
				return err
			}
		}
	}
	return nil
}
