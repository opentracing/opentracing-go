package opentracing

var (
	defaultOpenTracer OpenTracer = noopOpenTracer{noopTraceContextSource{}}
)

// InitDefaultTracer sets the [singleton] OpenTracer returned by
// DefaultTracer(). Those who use DefaultTracer (rather than directly manage an
// OpenTracer instance) should call InitDefaultTracer as early as possible in
// main(), prior to calling the `StartTrace` (etc) global funcs below. Prior to
// calling `InitDefaultTracer`, any Spans started via the `StartTrace` (etc)
// globals are noops.
func InitDefaultTracer(tracer OpenTracer) {
	defaultOpenTracer = tracer
}

// DefaultTracer returns the global singleton `OpenTracer` implementation.
// Before `InitDefaultTracer()` is called, the `DefaultTracer()` is a noop
// implementation that drops all data handed to it.
func DefaultTracer() OpenTracer {
	return defaultOpenTracer
}

// StartTrace defers to `OpenTracer.StartTrace`. See `DefaultTracer()`.
func StartTrace(operationName string, keyValueTags ...interface{}) Span {
	return defaultOpenTracer.StartTrace(operationName, keyValueTags...)
}

// JoinTrace defers to `OpenTracer.JoinTrace`. See `DefaultTracer()`.
func JoinTrace(operationName string, parent interface{}, keyValueTags ...interface{}) Span {
	return defaultOpenTracer.JoinTrace(operationName, parent, keyValueTags...)
}

// MarshalTraceContextBinary defers to
// `TraceContextMarshaler.MarshalTraceContextBinary`.
//
// See `DefaultTracer()`.
func MarshalTraceContextBinary(ctx TraceContext) []byte {
	return defaultOpenTracer.MarshalTraceContextBinary(ctx)
}

// MarshalTraceContextStringMap defers to
// `TraceContextMarshaler.MarshalTraceContextStringMap`.
//
// See `DefaultTracer()`.
func MarshalTraceContextStringMap(ctx TraceContext) map[string]string {
	return defaultOpenTracer.MarshalTraceContextStringMap(ctx)
}

// UnmarshalTraceContextBinary defers to
// `TraceContextUnmarshaler.UnmarshalTraceContextBinary`.
//
// See `DefaultTracer()`.
func UnmarshalTraceContextBinary(encoded []byte) (TraceContext, error) {
	return defaultOpenTracer.UnmarshalTraceContextBinary(encoded)
}

// UnmarshalTraceContextStringMap defers to
// `TraceContextUnmarshaler.UnmarshaTraceContextStringMap`.
//
// See `DefaultTracer()`.
func UnmarshalTraceContextStringMap(encoded map[string]string) (TraceContext, error) {
	return defaultOpenTracer.UnmarshalTraceContextStringMap(encoded)
}
