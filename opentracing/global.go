package opentracing

var (
	globalOpenTracer OpenTracer = noopOpenTracer{noopTraceContextIDSource{}}
)

// Should be called as early as possible in main(), prior to calling the
// `StartSpan*` (etc) global funcs below. Prior to calling `InitGlobal`, any
// Spans started via the `StartSpan*` globals are noops.
func InitGlobalTracer(rec ProcessRecorder, ctxIDSource TraceContextIDSource) {
	globalOpenTracer = NewStandardTracer(rec, ctxIDSource)
}

func Global() OpenTracer {
	return globalOpenTracer
}
