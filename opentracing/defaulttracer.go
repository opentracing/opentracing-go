package opentracing

var (
	defaultTracer Tracer = noopTracer{noopTraceContextSource{}}
)

// InitDefaultTracer sets the [singleton] opentracing.Tracer returned by
// DefaultTracer(). Those who use DefaultTracer (rather than directly manage an
// opentracing.Tracer instance) should call InitDefaultTracer as early as possible in
// main(), prior to calling the `StartTrace` (etc) global funcs below. Prior to
// calling `InitDefaultTracer`, any Spans started via the `StartTrace` (etc)
// globals are noops.
func InitDefaultTracer(tracer Tracer) {
	defaultTracer = tracer
}

// DefaultTracer returns the global singleton `Tracer` implementation.
// Before `InitDefaultTracer()` is called, the `DefaultTracer()` is a noop
// implementation that drops all data handed to it.
func DefaultTracer() Tracer {
	return defaultTracer
}

// StartTrace defers to `Tracer.StartTrace`. See `DefaultTracer()`.
func StartTrace(operationName string) Span {
	return defaultTracer.StartTrace(operationName)
}

// JoinTrace defers to `Tracer.JoinTrace`. See `DefaultTracer()`.
func JoinTrace(operationName string, parent interface{}) Span {
	return defaultTracer.JoinTrace(operationName, parent)
}

// MarshalTraceContextBinary defers to
// `TraceContextMarshaler.MarshalTraceContextBinary`.
//
// See `DefaultTracer()`.
func MarshalTraceContextBinary(ctx TraceContext) ([]byte, []byte) {
	return defaultTracer.MarshalTraceContextBinary(ctx)
}

// MarshalTraceContextStringMap defers to
// `TraceContextMarshaler.MarshalTraceContextStringMap`.
//
// See `DefaultTracer()`.
func MarshalTraceContextStringMap(ctx TraceContext) (map[string]string, map[string]string) {
	return defaultTracer.MarshalTraceContextStringMap(ctx)
}

// UnmarshalTraceContextBinary defers to
// `TraceContextUnmarshaler.UnmarshalTraceContextBinary`.
//
// See `DefaultTracer()`.
func UnmarshalTraceContextBinary(traceContextID []byte, traceTags []byte) (TraceContext, error) {
	return defaultTracer.UnmarshalTraceContextBinary(traceContextID, traceTags)
}

// UnmarshalTraceContextStringMap defers to
// `TraceContextUnmarshaler.UnmarshaTraceContextStringMap`.
//
// See `DefaultTracer()`.
func UnmarshalTraceContextStringMap(traceContextID map[string]string, traceTags map[string]string) (TraceContext, error) {
	return defaultTracer.UnmarshalTraceContextStringMap(traceContextID, traceTags)
}
