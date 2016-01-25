package opentracing

var (
	globalTracer Tracer = noopTracer{}
)

// InitGlobalTracer sets the [singleton] opentracing.Tracer returned by
// GlobalTracer(). Those who use GlobalTracer (rather than directly manage an
// opentracing.Tracer instance) should call InitGlobalTracer as early as possible in
// main(), prior to calling the `StartTrace` (etc) global funcs below. Prior to
// calling `InitGlobalTracer`, any Spans started via the `StartTrace` (etc)
// globals are noops.
func InitGlobalTracer(tracer Tracer) {
	globalTracer = tracer
}

// GlobalTracer returns the global singleton `Tracer` implementation.
// Before `InitGlobalTracer()` is called, the `GlobalTracer()` is a noop
// implementation that drops all data handed to it.
func GlobalTracer() Tracer {
	return globalTracer
}

// StartTrace defers to `Tracer.StartTrace`. See `GlobalTracer()`.
func StartTrace(operationName string) Span {
	return globalTracer.StartTrace(operationName)
}
