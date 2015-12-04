package opentracing

// OpenTracer is a simple, thin interface for Span creation.
//
// A straightforward implementation is available via the
// `opentracing/standardtracer` package's `standardtracer.New()'.
type OpenTracer interface {
	TraceContextSource

	// Create, start, and return a new Span with the given `operationName`, all
	// without specifying a parent Span that can be used to incorporate the
	// newly-returned Span into an existing trace.
	//
	// `keyValueTags` is a [possibly empty] list of <string, interface{}> tag
	// pairs which are applied to the returned `Span`, just as if they were
	// applied in successive calls to `Span.SetTag()`. If the even-numbered
	// `keyValueTags` parameters are not of `string` type, the implementation
	// is undefined (and may `panic()`).
	//
	// Examples:
	//
	//     val tracer OpenTracer = ...
	//
	//     sp := tracer.StartTrace("GetFeed")
	//
	//     sp := tracer.StartTrace(
	//         "HandleHTTPRequest",
	//         "user_agent", req.UserAgent,
	//         "lucky_number", 42)
	//
	StartTrace(operationName string, keyValueTags ...interface{}) Span

	// Like `StartTrace`, but the return `Span` is made a child of `parent`.
	//
	// The `parent` parameter can either be a `context.Context` or an
	// `opentracing.TraceContext`. In the former case, the implementation
	// attempts to extract an `opentracing.Span` using `SpanFromGoContext()`.
	JoinTrace(operationName string, parent interface{}, keyValueTags ...interface{}) Span
}
