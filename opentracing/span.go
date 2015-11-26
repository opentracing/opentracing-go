package opentracing

import "golang.org/x/net/context"

type Span interface {
	// Creates a child span, optionally modifying a `context.Context` parent.
	//
	// `parent` is optional; if specified, it is used as the descendant for the
	// returned `context.Context` object.
	//
	// Regardless of whether `parent` is specified, `StartChildSpan` returns a
	// `Span` that descends directly from the callee. The returned
	// `context.Context` instance derives from `parent` if.f. it was specified.
	StartChildSpan(operationName string, parent ...context.Context) (Span, context.Context)

	// Adds a tag to the span. The `value` is immediately coerced into a string
	// using fmt.Sprint().
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	SetTag(key string, value interface{})

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

	// Like Info(), but for warnings.
	Warning(message string, payload ...interface{})

	// Like Info(), but for errors.
	Error(message string, payload ...interface{})

	// Sets the end timestamp and calls the `Recorder`s RecordSpan() internally.
	//
	// Calling Finish() multiple times on the same Span instance may lead to
	// undefined behavior.
	Finish()

	// Suitable for serializing over the wire, etc.
	ContextID() ContextID
}

// A simple, thin interface for Span creation. Though other implementations are
// possible and plausible, most users will be fine with `NewStandardTracer()`.
type OpenTracer interface {
	ContextIDSource

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
	StartSpan(operationName string, parent ...interface{}) (Span, context.Context)
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
