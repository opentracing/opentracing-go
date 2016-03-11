package opentracing

// A NoopTracer is a trivial implementation of Tracer for which all operations
// are no-ops.
type NoopTracer struct{}

type noopSpan struct{}

var (
	defaultNoopSpan   = noopSpan{}
	defaultNoopTracer = NoopTracer{}
)

const (
	emptyString = ""
)

// noopSpan:
func (n noopSpan) SetTag(key string, value interface{}) Span             { return n }
func (n noopSpan) Finish()                                               {}
func (n noopSpan) FinishWithOptions(opts FinishOptions)                  {}
func (n noopSpan) SetBaggageItem(key, val string) Span                   { return n }
func (n noopSpan) BaggageItem(key string) string                         { return emptyString }
func (n noopSpan) LogEvent(event string)                                 {}
func (n noopSpan) LogEventWithPayload(event string, payload interface{}) {}
func (n noopSpan) Log(data LogData)                                      {}
func (n noopSpan) SetOperationName(operationName string) Span            { return n }
func (n noopSpan) Tracer() Tracer                                        { return defaultNoopTracer }

// StartSpan belongs to the Tracer interface.
func (n NoopTracer) StartSpan(operationName string) Span {
	return defaultNoopSpan
}

// StartSpanWithOptions belongs to the Tracer interface.
func (n NoopTracer) StartSpanWithOptions(opts StartSpanOptions) Span {
	return defaultNoopSpan
}

// Inject belongs to the Tracer interface.
func (n NoopTracer) Inject(sp Span, format interface{}, carrier interface{}) error {
	return nil
}

// Join belongs to the Tracer interface.
func (n NoopTracer) Join(operationName string, format interface{}, carrier interface{}) (Span, error) {
	return nil, ErrTraceNotFound
}
