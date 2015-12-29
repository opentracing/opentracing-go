package opentracing

var (
	globalTracer Tracer = noopTracer{noopTraceContextSource{}}
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

// JoinTrace defers to `Tracer.JoinTrace`. See `GlobalTracer()`.
func JoinTrace(operationName string, parent interface{}) Span {
	return globalTracer.JoinTrace(operationName, parent)
}

// MarshalTraceContextBinary defers to
// `TraceContextMarshaler.MarshalTraceContextBinary`.
//
// See `GlobalTracer()`.
func MarshalTraceContextBinary(ctx TraceContext) ([]byte, []byte) {
	return globalTracer.MarshalTraceContextBinary(ctx)
}

// MarshalTraceContextStringMap defers to
// `TraceContextMarshaler.MarshalTraceContextStringMap`.
//
// See `GlobalTracer()`.
func MarshalTraceContextStringMap(ctx TraceContext) (map[string]string, map[string]string) {
	return globalTracer.MarshalTraceContextStringMap(ctx)
}

// UnmarshalTraceContextBinary defers to
// `TraceContextUnmarshaler.UnmarshalTraceContextBinary`.
//
// See `GlobalTracer()`.
func UnmarshalTraceContextBinary(traceContextID []byte, traceTags []byte) (TraceContext, error) {
	return globalTracer.UnmarshalTraceContextBinary(traceContextID, traceTags)
}

// UnmarshalTraceContextStringMap defers to
// `TraceContextUnmarshaler.UnmarshaTraceContextStringMap`.
//
// See `GlobalTracer()`.
func UnmarshalTraceContextStringMap(traceContextID map[string]string, traceTags map[string]string) (TraceContext, error) {
	return globalTracer.UnmarshalTraceContextStringMap(traceContextID, traceTags)
}
