package opentracing

// Note: The Observer API is at an alpha stage and it is subjected to change.
//
// An observer can be registered with the Tracer to recieve notifications
// about new Spans. Tracers are not required to support the Observer API.
// The actual registration depends on the implementation, which might look
// like the below e.g :
// observer := myobserver.NewObserver()
// tracer := client.NewTracer(..., client.WithObserver(observer))
//
type Observer interface {
	// Create and return a span observer. Called when a span starts.
	// E.g :
	//     func StartSpan(opName string, opts ...opentracing.StartSpanOption) {
	//     var sp opentracing.Span
	//     sso := opentracing.StartSpanOptions{}
	//     var spObs opentracing.SpanObserver = observer.OnStartSpan(span, opName, sso)
	//     ...
	// }
	// OnStartSpan function needs to be defined for a package exporting
	// metrics as well.
	OnStartSpan(sp Span, operationName string, options StartSpanOptions) SpanObserver
}

// SpanObserver is created by the Observer and receives notifications about
// other Span events.
// Client tracers should define these functions for each of the span operations
// which should call the registered (observer) callbacks.
type SpanObserver interface {
	// Callback called from opentracing.Span.SetOperationName()
	OnSetOperationName(operationName string)
	// Callback called from opentracing.Span.SetTag()
	OnSetTag(key string, value interface{})
	// Callback called from opentracing.Span.Finish()
	OnFinish(options FinishOptions)
}
