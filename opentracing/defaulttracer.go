package opentracing

var (
	defaultOpenTracer OpenTracer = noopOpenTracer{noopTraceContextSource{}}
)

// Should be called as early as possible in main(), prior to calling the
// `StartSpan*` (etc) global funcs below. Prior to calling `InitDefaultTracer`,
// any Spans started via the `StartSpan*` globals are noops.
//
// See `NewStandardTracer(...)` to create an `OpenTracer` instance.
func InitDefaultTracer(tracer OpenTracer) {
	defaultOpenTracer = tracer
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

// Defers to `TraceContextMarshaler.MarshalTraceContextBinary`.
//
// See `DefaultTracer()`.
func MarshalTraceContextBinary(ctx TraceContext) []byte {
	return defaultOpenTracer.MarshalTraceContextBinary(ctx)
}

// Defers to `TraceContextMarshaler.MarshalStringMapTraceContext`.
//
// See `DefaultTracer()`.
func MarshalTraceContextStringMap(ctx TraceContext) map[string]string {
	return defaultOpenTracer.MarshalTraceContextStringMap(ctx)
}

// Defers to `TraceContextUnmarshaler.UnmarshalTraceContextBinary`.
//
// See `DefaultTracer()`.
func UnmarshalTraceContextBinary(encoded []byte) (TraceContext, error) {
	return defaultOpenTracer.UnmarshalTraceContextBinary(encoded)
}

// Defers to `TraceContextUnmarshaler.UnmarshalStringMapTraceContext`.
//
// See `DefaultTracer()`.
func UnmarshalTraceContextStringMap(encoded map[string]string) (TraceContext, error) {
	return defaultOpenTracer.UnmarshalTraceContextStringMap(encoded)
}
