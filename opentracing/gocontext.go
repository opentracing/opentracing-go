package opentracing

import "golang.org/x/net/context"

type goContextKey int

const activeSpanKey goContextKey = iota

// GoContextWithSpan seturns a new `context.Context` that holds a reference to
// the given `Span`.
//
// NOTE: We use the term "GoContext" to minimize confusion with TraceContext.
func GoContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, activeSpanKey, span)
}

// BackgroundGoContextWithSpan is a convenience wrapper around
// `GoContextWithSpan(context.BackgroundContext(), ...)`.
//
// NOTE: We use the term "GoContext" to minimize confusion with TraceContext.
func BackgroundGoContextWithSpan(span Span) context.Context {
	return context.WithValue(context.Background(), activeSpanKey, span)
}

// SpanFromGoContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
//
// NOTE: We use the term "GoContext" to minimize confusion with TraceContext.
func SpanFromGoContext(ctx context.Context) Span {
	val := ctx.Value(activeSpanKey)
	if span, ok := val.(Span); ok {
		return span
	}
	return nil
}
