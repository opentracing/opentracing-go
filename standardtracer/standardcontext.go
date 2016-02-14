package standardtracer

// StandardContext holds the basic Span metadata.
type StandardContext struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID int64

	// A probabilistically unique identifier for a span.
	SpanID int64

	// The SpanID of this StandardContext's parent, or 0 if there is no parent.
	ParentSpanID int64

	// Whether the trace is sampled.
	Sampled bool
}
