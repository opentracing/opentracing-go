package opentracing

// Observer can be registered with the Tracer to recieve notifications
// about new Spans.
type Observer interface {
	OnStartSpan(sp Span, operationName string, options StartSpanOptions) SpanObserver
}

// SpanObserver is created by the Observer and receives notifications
// about other Span events.
type SpanObserver interface {
	OnSetOperationName(operationName string)
	OnSetTag(key string, value interface{})
	OnFinish(options FinishOptions)
}
