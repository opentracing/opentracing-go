package opentracing

// A NoopTracer is a trivial implementation of Tracer for which all operations
// are no-ops.
type NoopTracer struct{}

type noopSpan struct{}
type noopInjectorExtractor struct{}

var (
	defaultNoopSpan              = noopSpan{}
	defaultNoopInjectorExtractor = noopInjectorExtractor{}
	defaultNoopTracer            = NoopTracer{}
	emptyTags                    = Tags{}
	emptyBytes                   = []byte{}
	emptyStringMap               = map[string]string{}
)

const (
	emptyString = ""
)

// noopSpan:
func (n noopSpan) SetTag(key string, value interface{}) Span             { return n }
func (n noopSpan) Finish()                                               {}
func (n noopSpan) FinishWithOptions(opts FinishOptions)                  {}
func (n noopSpan) SetTraceAttribute(key, val string) Span                { return n }
func (n noopSpan) TraceAttribute(key string) string                      { return emptyString }
func (n noopSpan) LogEvent(event string)                                 {}
func (n noopSpan) LogEventWithPayload(event string, payload interface{}) {}
func (n noopSpan) Log(data LogData)                                      {}
func (n noopSpan) SetOperationName(operationName string) Span            { return n }
func (n noopSpan) Tracer() Tracer                                        { return defaultNoopTracer }

// PropagateSpanAsBinary belongs to the Tracer interface.
func (n NoopTracer) PropagateSpanAsBinary(tcid Span) ([]byte, []byte) {
	return emptyBytes, emptyBytes
}

// PropagateSpanAsText belongs to the Tracer interface.
func (n NoopTracer) PropagateSpanAsText(tcid Span) (map[string]string, map[string]string) {
	return emptyStringMap, emptyStringMap
}

// JoinTraceFromBinary belongs to the Tracer interface.
func (n NoopTracer) JoinTraceFromBinary(
	op string,
	traceContextID []byte,
	traceAttrs []byte,
) (Span, error) {
	return defaultNoopSpan, nil
}

// JoinTraceFromText belongs to the Tracer interface.
func (n NoopTracer) JoinTraceFromText(
	op string,
	traceContextID map[string]string,
	traceAttrs map[string]string,
) (Span, error) {
	return defaultNoopSpan, nil
}

// StartSpan belongs to the Tracer interface.
func (n NoopTracer) StartSpan(operationName string) Span {
	return defaultNoopSpan
}

// StartSpanWithOptions belongs to the Tracer interface.
func (n NoopTracer) StartSpanWithOptions(opts StartSpanOptions) Span {
	return defaultNoopSpan
}

// Extractor belongs to the Tracer interface.
func (n NoopTracer) Extractor(format interface{}) Extractor {
	return defaultNoopInjectorExtractor
}

// Injector belongs to the Tracer interface.
func (n NoopTracer) Injector(format interface{}) Injector {
	return defaultNoopInjectorExtractor
}

// JoinTrace belongs to the Tracer interface.
func (n NoopTracer) JoinTrace(operationName string, parent interface{}) Span {
	return defaultNoopSpan
}

func (n noopInjectorExtractor) InjectSpan(span Span, carrier interface{}) error {
	return nil
}

func (n noopInjectorExtractor) JoinTrace(operationName string, carrier interface{}) (Span, error) {
	return nil, nil
}
