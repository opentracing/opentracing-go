package opentracing

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
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

	// TextMap encodes the Span in a TextMapCarrier instance.
	//
	// For Tracer.Inject(): the carrier must be a `TextMapWriter`.
	//
	// For Tracer.Join(): the carrier must be a `TextMapReader`.
	//
	// See HTTPHeaderTextMapCarrier for an implementation of both TextMapWriter
	// and TextMapReader that defers to an http.Header instance for storage.
	// For example, Inject():
	//
	//    carrier := HTTPHeaderTextMapCarrier{
	//    	HeaderPrefix: "opentracing-",
	//    	Header:       h,
	//    }
	//    err := span.Tracer().Inject(span, TextMap, carrier)
	//
	// Or Join():
	//
	//    carrier := HTTPHeaderTextMapCarrier{
	//    	HeaderPrefix: "opentracing-",
	//    	Header:       h,
	//    }
	//    span, err := tracer.Join("opName", TextMap, carrier)
	//
	TextMap

	// SplitBinary is DEPRECATED
	SplitBinary
	// SplitText is DEPRECATED
	SplitText
)

// TextMapWriter is the Inject() carrier for the TextMap builtin format. With
// it, the caller can encode a Span for propagation as entries in a multimap of
// unicode strings.
type TextMapWriter interface {
	// Add a key:value pair to the carrier. Multiple values may be added for a
	// single (repeated) key.
	Add(key, val string)
}

// TextMapWriter is the Join() carrier for the TextMap builtin format. With it,
// the caller can decode a propagated Span as entries in a multimap of unicode
// strings.
type TextMapReader interface {
	// ReadAllEntries returns TextMap contents via repeated calls to the
	// `handler` function. If any call to `handler` returns a non-nil error,
	// ReadAllEntries terminates and returns that error.
	//
	// NOTE: A single `key` may appear in multiple calls to `handler` for a
	// single `ReadAllEntries` invocation.
	ReadAllEntries(handler func(key, val string) error) error
}

// HTTPHeaderTextMapCarrier satisfies both TextMapWriter and TextMapReader.
//
// NOTE: All `key` parameters to Add() and the ReadAllEntries() handler func
// are lowercased since http.Header doesn't respect character casing for keys.
type HTTPHeaderTextMapCarrier struct {
	// The prefix used to distinguish the TextMap entries within the
	// http.Header map.
	HeaderPrefix string

	http.Header
}

var ErrNoHeaderPrefix = errors.New("HTTPHeaderTextMapCarrier.HeaderPrefix is empty")

func (c HTTPHeaderTextMapCarrier) Add(key, val string) {
	c.Header.Add(strings.ToLower(c.HeaderPrefix+key), url.QueryEscape(val))
}
func (c HTTPHeaderTextMapCarrier) ReadAllEntries(handler func(key, val string) error) error {
	if len(c.HeaderPrefix) == 0 {
		return ErrNoHeaderPrefix
	}
	for k, vals := range c.Header {
		k = strings.ToLower(k)
		if !strings.HasPrefix(k, c.HeaderPrefix) {
			continue
		}
		kSuffix := k[len(c.HeaderPrefix):]
		for _, v := range vals {
			rawV, err := url.QueryUnescape(v)
			if err != nil {
				continue
			}
			if err = handler(kSuffix, rawV); err != nil {
				return err
			}
		}
	}
	return nil
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
