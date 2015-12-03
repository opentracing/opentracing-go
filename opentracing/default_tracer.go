package opentracing

var (
	globalOpenTracer OpenTracer = noopOpenTracer{noopTraceContextSource{}}
)

// Should be called as early as possible in main(), prior to calling the
// `StartSpan*` (etc) global funcs below. Prior to calling `InitDefaultTracer`,
// any Spans started via the `StartSpan*` globals are noops.
func InitDefaultTracer(rec ProcessRecorder, ctxIDSource TraceContextSource) {
	globalOpenTracer = NewStandardTracer(rec, ctxIDSource)
}

// Return the global singleton `OpenTracer` implementation. Before
// `InitDefaultTracer()` is called, the `DefaultTracer()` is a noop
// implementation that drops all data handed to it.
func DefaultTracer() OpenTracer {
	return globalOpenTracer
}

// See `OpenTracer.StartTrace` and `InitDefaultTracer()`.
func StartTrace(operationName string, keyValueTags ...interface{}) Span {
	return globalOpenTracer.StartTrace(operationName, keyValueTags...)
}

// See `OpenTracer.JoinTrace` and `InitDefaultTracer()`.
func JoinTrace(operationName string, parent interface{}, keyValueTags ...interface{}) Span {
	return globalOpenTracer.JoinTrace(operationName, parent, keyValueTags...)
}

// Defers to `MarshalBinaryTraceContext()`. See `InitDefaultTracer()`.
func DefaultMarshalBinaryTraceContext(ctx TraceContext) []byte {
	return globalOpenTracer.MarshalBinaryTraceContext(ctx)
}

// Defers to `MarshalStringMapTraceContext()`. See `InitDefaultTracer()`.
func DefaultMarshalStringMapTraceContext(ctx TraceContext) map[string]string {
	return globalOpenTracer.MarshalStringMapTraceContext(ctx)
}

// Defers to `UnmarshalBinaryTraceContext()`. See `InitDefaultTracer()`.
func DefaultUnmarshalBinaryTraceContext(encoded []byte) (TraceContext, error) {
	return globalOpenTracer.UnmarshalBinaryTraceContext(encoded)
}

// Defers to `UnmarshalStringMapTraceContext()`. See `InitDefaultTracer()`.
func DefaultUnmarshalStringMapTraceContext(encoded map[string]string) (TraceContext, error) {
	return globalOpenTracer.UnmarshalStringMapTraceContext(encoded)
}
