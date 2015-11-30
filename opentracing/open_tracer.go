package opentracing

// A simple, thin interface for Span creation. Though other implementations are
// possible and plausible, most users will be fine with `NewStandardTracer()`.
type OpenTracer interface {
	TraceContextIDSource

	StartNewTrace(operationName string, initialTags ...Tags) Span

	// `parent` can either be a `context.Context` or an
	// `opentracing.TraceContext`.
	JoinTrace(operationName string, parent interface{}, initialTags ...Tags) Span

	// XXX: adapt comment below to StartNewTrace / JoinTrace above.

	// Starts a new Span for `operationName`.
	//
	// If `parent` is a golang `context.Context`, the returned
	// `context.Context` and `Span` are schematic children of that context and
	// any `Span` found therein.
	//
	// If `parent` is an `opentracing.ContextID`, the returned
	// `context.Context` descends from the `context.Background()` and the
	// returned `Span` descends from the provided `opentracing.ContextID`.
	//
	// If `parent` is omitted, the returned `Span` is a "root" span: i.e., it
	// has no known parent.
	// StartSpan(operationName string, parent ...interface{}) (Span, context.Context)
}
