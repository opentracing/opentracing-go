package opentracing

var (
	defaultOpenTracer OpenTracer = noopOpenTracer{noopTraceContextSource{}}
)

// Should be called as early as possible in main(), prior to calling the
// `StartSpan*` (etc) global funcs below. Prior to calling `InitDefaultTracer`,
// any Spans started via the `StartSpan*` globals are noops.
func InitDefaultTracer(rec Recorder, ctxIDSource TraceContextSource) {
	defaultOpenTracer = NewStandardTracer(rec, ctxIDSource)
}

// Return the global singleton `OpenTracer` implementation. Before
// `InitDefaultTracer()` is called, the `DefaultTracer()` is a noop
// implementation that drops all data handed to it.
func DefaultTracer() OpenTracer {
	return defaultOpenTracer
}

// Defers to `OpenTracer.StartTrace`. See `DefaultTracer()`.
func StartTrace(operationName string, keyValueTags ...interface{}) Span {
	return defaultOpenTracer.StartTrace(operationName, keyValueTags...)
}

// Defers to `OpenTracer.JoinTrace`. See `DefaultTracer()`.
func JoinTrace(operationName string, parent interface{}, keyValueTags ...interface{}) Span {
	return defaultOpenTracer.JoinTrace(operationName, parent, keyValueTags...)
}

// Defers to `TraceContextMarshaler.MarshalBinaryTraceContext`.
//
// See `DefaultTracer()`.
func MarshalBinaryTraceContext(ctx TraceContext) []byte {
	return defaultOpenTracer.MarshalBinaryTraceContext(ctx)
}

// Defers to `TraceContextMarshaler.MarshalStringMapTraceContext`.
//
// See `DefaultTracer()`.
func MarshalStringMapTraceContext(ctx TraceContext) map[string]string {
	return defaultOpenTracer.MarshalStringMapTraceContext(ctx)
}

// Defers to `TraceContextUnmarshaler.UnmarshalBinaryTraceContext`.
//
// See `DefaultTracer()`.
func UnmarshalBinaryTraceContext(encoded []byte) (TraceContext, error) {
	return defaultOpenTracer.UnmarshalBinaryTraceContext(encoded)
}

// Defers to `TraceContextUnmarshaler.UnmarshalStringMapTraceContext`.
//
// See `DefaultTracer()`.
func UnmarshalStringMapTraceContext(encoded map[string]string) (TraceContext, error) {
	return defaultOpenTracer.UnmarshalStringMapTraceContext(encoded)
}
