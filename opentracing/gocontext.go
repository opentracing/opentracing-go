package opentracing

import "golang.org/x/net/context"

type goContextKey int

const activeSpanKey goContextKey = iota

// Return a new `context.Context` including a reference to the given `Span`.
func GoContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, activeSpanKey, span)
}

// A convenience wrapper around `GoContextWithSpan(context.BackgroundContext(), ...)`.
func BackgroundContextWithSpan(span Span) context.Context {
	return context.WithValue(context.Background(), activeSpanKey, span)
}

// Returns the `Span` previously associated with `ctx`, or `nil` if no such
// `Span` could be found.
func SpanFromGoContext(ctx context.Context) Span {
	val := ctx.Value(activeSpanKey)
	if span, ok := val.(Span); ok {
		return span
	}
	return nil
}
