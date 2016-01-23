package standardtracer

/*
type TraceContextSource interface {
	TraceContextEncoder
	TraceContextDecoder

	// Create a TraceContext which has no parent (and thus begins its own trace).
	// A TraceContextSource must always return the same type in successive calls
	// to NewRootTraceContext().
	NewRootTraceContext() TraceContext

	// NewChildTraceContext creates a child context for `parent`, and returns
	// both that child's own TraceContext as well as any Tags that should be
	// added to the child's Span.
	//
	// The returned TraceContext type must be the same as the type of the
	// TraceContext implementation itself.
	NewChildTraceContext(parent TraceContext) (childCtx TraceContext, childSpanTags Tags)
}
*/
