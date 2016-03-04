package opentracing

import (
	"testing"
	"golang.org/x/net/context"
)

func TestContextWithSpan(t *testing.T) {
	span := &noopSpan{}
	ctx := BackgroundContextWithSpan(span)
	span2 := SpanFromContext(ctx)
	if span != span2 {
		t.Errorf("Not the same span returned from context, expected=%+v, actual=%+v", span, span2)
	}

	ctx = context.Background()
	span2 = SpanFromContext(ctx)
	if span2 != nil {
		t.Errorf("Expected nil span, found %+v", span2)
	}

	ctx = ContextWithSpan(ctx, span)
	span2 = SpanFromContext(ctx)
	if span != span2 {
		t.Errorf("Not the same span returned from context, expected=%+v, actual=%+v", span, span2)
	}
}
