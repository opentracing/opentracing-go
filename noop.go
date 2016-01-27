package opentracing

type noopSpan struct{}
type noopTracer struct{}

var (
	defaultNoopSpan   = noopSpan{}
	defaultNoopTracer = noopTracer{}
	emptyTags         = Tags{}
	emptyBytes        = []byte{}
	emptyStringMap    = map[string]string{}
)

const (
	emptyString = ""
)

// noopSpan:
func (n noopSpan) StartChild(operationName string) Span {
	return defaultNoopSpan
}
func (n noopSpan) SetTag(key string, value interface{}) Span             { return n }
func (n noopSpan) Finish()                                               {}
func (n noopSpan) SetTraceAttribute(key, val string) Span                { return n }
func (n noopSpan) TraceAttribute(key string) string                      { return emptyString }
func (n noopSpan) LogEvent(event string)                                 {}
func (n noopSpan) LogEventWithPayload(event string, payload interface{}) {}
func (n noopSpan) Log(data LogData)                                      {}
func (n noopSpan) SetOperationName(operationName string) Span            { return n }

// noopTracer:
func (n noopTracer) PropagateSpanAsBinary(tcid Span) ([]byte, []byte) {
	return emptyBytes, emptyBytes
}
func (n noopTracer) PropagateSpanAsText(tcid Span) (map[string]string, map[string]string) {
	return emptyStringMap, emptyStringMap
}
func (n noopTracer) JoinTraceFromBinary(
	op string,
	traceContextID []byte,
	traceAttrs []byte,
) (Span, error) {
	return defaultNoopSpan, nil
}
func (n noopTracer) JoinTraceFromText(
	op string,
	traceContextID map[string]string,
	traceAttrs map[string]string,
) (Span, error) {
	return defaultNoopSpan, nil
}

func (n noopTracer) StartTrace(operationName string) Span {
	return defaultNoopSpan
}

func (n noopTracer) JoinTrace(operationName string, parent interface{}) Span {
	return defaultNoopSpan
}
