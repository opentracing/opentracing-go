package opentracing

import "time"

// Tracer is a simple, thin interface for Span creation.
//
// A straightforward implementation is available via the
// `opentracing/standardtracer` package's `standardtracer.New()'.
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

	// Return a PropagationExtractor for the given PropagationFormat, or nil if the Tracer
	// does not support such a format.
	PropagationExtractorForFormat(format PropagationFormat) PropagationExtractor
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
