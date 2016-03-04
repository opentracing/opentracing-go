package opentracing

import "golang.org/x/net/context"

type contextKey struct{}

var activeSpanKey = contextKey{}

// ContextWithSpan returns a new `context.Context` that holds a reference to
// the given `Span`.
func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, activeSpanKey, span)
}

// BackgroundContextWithSpan is a convenience wrapper around
// `ContextWithSpan(context.BackgroundContext(), ...)`.
func BackgroundContextWithSpan(span Span) context.Context {
	return context.WithValue(context.Background(), activeSpanKey, span)
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
func SpanFromContext(ctx context.Context) Span {
	val := ctx.Value(activeSpanKey)
	if span, ok := val.(Span); ok {
		return span
	}
	return nil
}
