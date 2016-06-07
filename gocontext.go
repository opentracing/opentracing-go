package opentracing

import "golang.org/x/net/context"

type contextKey struct{}

var activeSpanKey = contextKey{}

// ContextWithSpan returns a new `context.Context` that holds a reference to
// the given `Span`'s `SpanContext`.
func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, activeSpanKey, span)
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
func SpanFromContext(ctx context.Context) Span {
	val := ctx.Value(activeSpanKey)
	if sp, ok := val.(Span); ok {
		return sp
	}
	return nil
}

// StartSpanFromContext starts and returns a Span with `operationName`, using
// any Span found within `ctx` as a parent. If no such parent could be found,
// StartSpanFromContext creates a root (parentless) Span.
//
// The second return value is a context.Context object built around the
// returned Span.
//
// Example usage:
//
//    SomeFunction(ctx context.Context, ...) {
//        sp, ctx := opentracing.StartSpanFromContext(ctx, "SomeFunction")
//        defer sp.Finish()
//        ...
//    }
func StartSpanFromContext(ctx context.Context, operationName string) (Span, context.Context) {
	return startSpanFromContextWithTracer(ctx, operationName, GlobalTracer())
}

// startSpanFromContextWithTracer is factored out for testing purposes.
func startSpanFromContextWithTracer(ctx context.Context, operationName string, tracer Tracer) (Span, context.Context) {
	var parentSpanCtx SpanContext
	if span := SpanFromContext(ctx); span != nil {
		parentSpanCtx = span.SpanContext()
	}
	span := tracer.StartSpanWithOptions(StartSpanOptions{
		OperationName: operationName,
		Parent:        parentSpanCtx,
	})
	return span, ContextWithSpan(ctx, span)
}
