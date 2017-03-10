package opentracing

// Observer can be registered with the Tracer to recieve notifications
// about new Spans.
// It is an interface to create a SpanObserver.
// Client tracer should have an observer initialized and assigned to the tracer.
// E.g :
// func WithObserver(observer opentracing.Observer) TracerOption {
//     return func(opts *TracerOptions) error {
//         opts.observer = observer
//         return nil
//     }
// }
//
// A Package exporting metrics (and wishing to use the Observer interface)
// needs to have an "Observer" struct and a function returning a new Observer.
// E.g :
// type Observer struct {}
// func NewObserver() *Observer {
//     return &Observer{}
// }
//
// Application using the metrics exporter package needs to create an
// observer for that package and send it to the "client" tracer.
// E.g :
// var observer opentracing.Observer = metricsexporter.NewObserver()
// tracer := client.NewTracer(..., client.WithObserver(observer))
//
type Observer interface {
	// Create and return a span oberver. Called when a span starts.
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

// SpanObserver is created by the Observer and receives notifications
// about other Span events.
// SpanObserver is created by the Observer and receives notifications about
// other Span events.
// Client tracers should define these functions for each of the span operations
// which call the registered (metrics exporter) callbacks.
// Metrics exporter packages need to define these functions to do operations
// on each of the span events.
type SpanObserver interface {
	// Callback called from opentracing.Span.SetOperationName()
	OnSetOperationName(operationName string)
	// Callback called from opentracing.Span.SetTag()
	OnSetTag(key string, value interface{})
	// Callback called from opentracing.Span.Finish()
	OnFinish(options FinishOptions)
}
