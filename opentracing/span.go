package opentracing

import "golang.org/x/net/context"

type Span interface {
	// Creates and starts a child span.
	// XXX: Tags
	StartChild(operationName string, initialTags ...Tags) (Span, context.Context)
	// Creates and starts a child span and adds it to the `context.Context`
	// argument, `parent`, before returning both.
	StartChildWithContext(operationName string, parent context.Context, initialTags ...Tags) (Span, context.Context)

	// Adds a tag to the span. The `value` is immediately coerced into a string
	// using fmt.Sprint().
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	SetTag(key string, value interface{}) Span

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

	// Sets the end timestamp and calls the `ProcessRecorder`s RecordSpan()
	// internally.
	//
	// Finish() should be the last call made to any span instance, and to do
	// otherwise leads to undefined behavior.
	Finish()

	// Suitable for serializing over the wire, etc.
	TraceContext() *TraceContext
}

// A simple, thin interface for Span creation. Though other implementations are
// possible and plausible, most users will be fine with `NewStandardTracer()`.
type OpenTracer interface {
	TraceContextIDSource

	StartEmptyTrace(
		operationName string, initialTags ...Tags,
	) (Span, context.Context)

	// `parent` can either be a `context.Context` or an
	// `opentracing.TraceContext`.
	ContinueTrace(
		operationName string, parent interface{}, initialTags ...Tags,
	) (Span, context.Context)

	// XXX START HERE: adapt comment below to BeginTrace / ContinueTrace above.

	// Starts a new Span for `operationName`.
	//
	// If `parent` is a golang `context.Context`, the returned
	// `context.Context` and `Span` are schematic children of that context and
	// any `Span` found therein.
	//
	// If `parent` is an `opentracing.ContextID`, the returned
	// `context.Context` descends from the `context.Background()` and the
	// returned `Span` descends from the provided `opentracing.ContextID`.
	//
	// If `parent` is omitted, the returned `Span` is a "root" span: i.e., it
	// has no known parent.
	// StartSpan(operationName string, parent ...interface{}) (Span, context.Context)
}

////////////////////////////////////
// begin context.Context boilerplate

type goContextKey int

const activeSpanKey goContextKey = iota

func GoContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, activeSpanKey, span)
}

func SpanFromGoContext(ctx context.Context) Span {
	val := ctx.Value(activeSpanKey)
	if span, ok := val.(Span); ok {
		return span
	}
	return nil
}

// end context.Context boilerplate
////////////////////////////////////
