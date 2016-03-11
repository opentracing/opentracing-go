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

func TestStartSpanFromContext(t *testing.T) {
	testTracer := testTracer{}

	// Test the case where there *is* a Span in the Context.
	{
		parentSpan := &testSpan{}
		parentCtx := BackgroundContextWithSpan(parentSpan)
		childSpan, childCtx := startSpanFromContextWithTracer(parentCtx, "child", testTracer)
		if !childSpan.(testSpan).HasParent {
			t.Errorf("Failed to find parent: %v", childSpan)
		}
		if childSpan != SpanFromContext(childCtx) {
			t.Errorf("Unable to find child span in context: %v", childCtx)
		}
	}

	// Test the case where there *is not* a Span in the Context.
	{
		emptyCtx := context.Background()
		childSpan, childCtx := startSpanFromContextWithTracer(emptyCtx, "child", testTracer)
		if childSpan.(testSpan).HasParent {
			t.Errorf("Should not have found parent: %v", childSpan)
		}
		if childSpan != SpanFromContext(childCtx) {
			t.Errorf("Unable to find child span in context: %v", childCtx)
		}
	}
}
