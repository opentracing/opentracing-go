package opentracing

import "golang.org/x/net/context"

type goContextKey int

const activeSpanKey goContextKey = iota

func BackgroundContextWithSpan(span Span) context.Context {
	return context.WithValue(context.Background(), activeSpanKey, span)
}

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
