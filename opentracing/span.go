package opentracing

import "golang.org/x/net/context"

// Span represents an active, un-finished span in the opentracing system.
//
// Spans are created by the Tracer interface and Span.StartChild.
type Span interface {
	// Creates and starts a child span.
	StartChild(operationName string) Span

	// Adds a tag to the span.
	//
	// Tag values can be of arbitrary types, however the treatment of complex
	// types is dependent on the underlying tracing system implementation.
	// It is expected that most tracing systems will handle primitive types
	// like strings and numbers. If a tracing system cannot understand how
	// to handle a particular value type, it may ignore the tag, but shall
	// not panic.
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	SetTag(key string, value interface{}) Span

	// SetTags adds multiple tags to this Span instance. Equivalent to calling
	// SetTag separately for each key:value pair.
	SetTags(tags Tags) Span

	// `Message` is a format string and can refer to fields in the payload by path, like so:
	//
	//   "first transaction is worth ${transactions[0].amount} ${transactions[0].currency}"
	//
	// , and the payload might look something like
	//
	//   map[string]interface{}{
	//       transactions: map[string]interface{}[
	//           {amount: 10, currency: "USD"},
	//           {amount: 11, currency: "USD"},
	//       ]}
	Info(message string, payload ...interface{})

	// Like Info(), but for errors.
	Error(message string, payload ...interface{})

	// Sets the end timestamp and calls the `Recorder`s RecordSpan()
	// internally.
	//
	// Finish() should be the last call made to any span instance, and to do
	// otherwise leads to undefined behavior.
	Finish()

	// Suitable for serializing over the wire, etc.
	TraceContext() TraceContext

	// A convenience method. Equivalent to
	//
	//    var goCtx context.Context = ...
	//    var span Span = ...
	//    goCtx := opentracing.GoContextWithSpan(ctx, span)
	//
	//
	// NOTE: We use the term "GoContext" to minimize confusion with
	// TraceContext.
	AddToGoContext(goCtx context.Context) (Span, context.Context)
}
