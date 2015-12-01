package opentracing

var (
	globalOpenTracer OpenTracer = noopOpenTracer{noopTraceContextIDSource{}}
)

// Should be called as early as possible in main(), prior to calling the
// `StartSpan*` (etc) global funcs below. Prior to calling `InitGlobal`, any
// Spans started via the `StartSpan*` globals are noops.
func InitGlobal(rec ProcessRecorder, ctxIDSource TraceContextIDSource) {
	globalOpenTracer = NewStandardTracer(rec, ctxIDSource)
}

// Return the global singleton `OpenTracer` implementation. Before
// `InitGlobal()` is called, the `Global()` is a noop implementation that drops
// all data handed to it.
func Global() OpenTracer {
	return globalOpenTracer
}

// See `OpenTracer.StartTrace` and `InitGlobal()`.
func StartTrace(operationName string, keyValueTags ...interface{}) Span {
	return globalOpenTracer.StartTrace(operationName, keyValueTags...)
}

// See `OpenTracer.JoinTrace` and `InitGlobal()`.
func JoinTrace(operationName string, parent interface{}, keyValueTags ...interface{}) Span {
	return globalOpenTracer.JoinTrace(operationName, parent, keyValueTags...)
}

// Defers to `DeserializeBinaryTraceContext()`. See `InitGlobal()`.
func GlobalDeserializeBinaryTraceContext(encoded []byte) (*TraceContext, error) {
	return DeserializeBinaryTraceContext(globalOpenTracer, encoded)
}

// Defers to `DeserializeASCIITraceContext()`. See `InitGlobal()`.
func GlobalDeserializeASCIITraceContext(encoded string) (*TraceContext, error) {
	return DeserializeASCIITraceContext(globalOpenTracer, encoded)
}
