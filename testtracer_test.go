package opentracing

// testTracer is a most-noop Tracer implementation that makes it possible for
// unittests to verify whether certain methods were / were not called.
type testTracer struct{}

type testSpan struct {
	OperationName string
	HasParent     bool
}

// testSpan:
func (n testSpan) SetTag(key string, value interface{}) Span             { return n }
func (n testSpan) Finish()                                               {}
func (n testSpan) FinishWithOptions(opts FinishOptions)                  {}
func (n testSpan) SetBaggageItem(key, val string) Span                   { return n }
func (n testSpan) BaggageItem(key string) string                         { return "" }
func (n testSpan) LogEvent(event string)                                 {}
func (n testSpan) LogEventWithPayload(event string, payload interface{}) {}
func (n testSpan) Log(data LogData)                                      {}
func (n testSpan) SetOperationName(operationName string) Span            { return n }
func (n testSpan) Tracer() Tracer                                        { return testTracer{} }

// StartSpan belongs to the Tracer interface.
func (n testTracer) StartSpan(operationName string) Span {
	return testSpan{
		OperationName: operationName,
		HasParent:     false,
	}
}

// StartSpanWithOptions belongs to the Tracer interface.
func (n testTracer) StartSpanWithOptions(opts StartSpanOptions) Span {
	return testSpan{
		OperationName: opts.OperationName,
		HasParent:     opts.Parent != nil,
	}
}

// Inject belongs to the Tracer interface.
func (n testTracer) Inject(sp Span, format interface{}, carrier interface{}) error {
	return nil
}

// Join belongs to the Tracer interface.
func (n testTracer) Join(operationName string, format interface{}, carrier interface{}) (Span, error) {
	return nil, ErrTraceNotFound
}
