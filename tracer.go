package opentracing

import "time"

// Tracer is a simple, thin interface for Span creation.
//
// A straightforward implementation is available via the
// `opentracing/basictracer-go` package's `standardtracer.New()'.
type Tracer interface {
	// Create, start, and return a new Span with the given `operationName`, all
	// without specifying a parent Span that can be used to incorporate the
	// newly-returned Span into an existing trace. (I.e., the returned Span is
	// the "root" of its trace).
	//
	// Examples:
	//
	//     var tracer opentracing.Tracer = ...
	//
	//     sp := tracer.StartSpan("GetFeed")
	//
	//     sp := tracer.StartSpanWithOptions(opentracing.SpanOptions{
	//         OperationName: "LoggedHTTPRequest",
	//         Tags: opentracing.Tags{"user_agent", loggedReq.UserAgent},
	//         StartTime: loggedReq.Timestamp,
	//     })
	//
	StartSpan(operationName string) Span
	StartSpanWithOptions(opts StartSpanOptions) Span

	// Inject() takes the `sp` Span instance and represents it for propagation
	// within `carrier`. The actual type of `carrier` depends on the value of
	// `format`.
	//
	// OpenTracing defines a common set of `format` values (see BuiltinFormat),
	// and each has an expected carrier type.
	//
	// Other packages may declare their own `format` values, much like the keys
	// used by the `net.Context` package (see
	// https://godoc.org/golang.org/x/net/context#WithValue).
	//
	// Example usage (sans error handling):
	//
	//     carrier := opentracing.HTTPHeaderTextMapCarrier(httpReq.Header)
	//     tracer.Inject(
	//         span,
	//         opentracing.TextMap,
	//         carrier)
	//
	// NOTE: All opentracing.Tracer implementations MUST support all
	// BuiltinFormats.
	//
	// Implementations may return opentracing.ErrUnsupportedFormat if `format`
	// is or not supported by (or not known by) the implementation.
	//
	// Implementations may return opentracing.ErrInvalidCarrier or any other
	// implementation-specific error if the format is supported but injection
	// fails anyway.
	//
	// See Tracer.Join().
	Inject(sp Span, format interface{}, carrier interface{}) error

	// Join() returns a Span instance with operation name `operationName` given
	// `format` and `carrier`.
	//
	// Join() is responsible for extracting and joining to the trace of a Span
	// instance embedded in a format-specific "carrier" object. Typically the
	// joining will take place on the server side of an RPC boundary, but
	// message queues and other IPC mechanisms are also reasonable places to
	// use Join().
	//
	// OpenTracing defines a common set of `format` values (see BuiltinFormat),
	// and each has an expected carrier type.
	//
	// Other packages may declare their own `format` values, much like the keys
	// used by the `net.Context` package (see
	// https://godoc.org/golang.org/x/net/context#WithValue).
	//
	// Example usage (sans error handling):
	//
	//     carrier := opentracing.HTTPHeaderTextMapCarrier(httpReq.Header)
	//     span, err := tracer.Join(
	//         operationName,
	//         opentracing.TextMap,
	//         carrier)
	//
	// NOTE: All opentracing.Tracer implementations MUST support all
	// BuiltinFormats.
	//
	// Return values:
	//  - A successful join will return a started Span instance and a nil error
	//  - If there was simply no trace to join with in `carrier`, Join()
	//    returns (nil, opentracing.ErrTraceNotFound)
	//  - If `format` is unsupported or unrecognized, Join() returns (nil,
	//    opentracing.ErrUnsupportedFormat)
	//  - If there are more fundamental problems with the `carrier` object,
	//    Join() may return opentracing.ErrInvalidCarrier,
	//    opentracing.ErrTraceCorrupted, or implementation-specific errors.
	//
	// See Tracer.Inject().
	Join(operationName string, format interface{}, carrier interface{}) (Span, error)
}

// StartSpanOptions allows Tracer.StartSpanWithOptions callers to override the
// start timestamp, specify a parent Span, and make sure that Tags are
// available at Span initialization time.
type StartSpanOptions struct {
	// OperationName may be empty (and set later via Span.SetOperationName)
	OperationName string

	// Parent may specify Span instance that caused the new (child) Span to be
	// created.
	//
	// If nil, start a "root" span (i.e., start a new trace).
	Parent Span

	// StartTime overrides the Span's start time, or implicitly becomes
	// time.Now() if StartTime.IsZero().
	StartTime time.Time

	// Tags may have zero or more entries; the restrictions on map values are
	// identical to those for Span.SetTag(). May be nil.
	//
	// If specified, the caller hands off ownership of Tags at
	// StartSpanWithOptions() invocation time.
	Tags map[string]interface{}
}
