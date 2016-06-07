package mocktracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opentracing/opentracing-go"
)

func TestMockTracer_StartSpan(t *testing.T) {
	tracer := New()
	span1 := tracer.StartSpan(
		"a",
		opentracing.Tags(map[string]interface{}{"x": "y"}))

	span2 := span1.Tracer().StartSpan(
		"", opentracing.ChildOf(span1.Context()))
	span2.Finish()
	span1.Finish()
	spans := tracer.GetFinishedSpans()
	assert.Equal(t, 2, len(spans))

	parent := spans[1]
	child := spans[0]
	assert.Equal(t, map[string]interface{}{"x": "y"}, parent.GetTags())
	assert.Equal(t, child.ParentID, parent.Context().(*MockSpanContext).SpanID)
}

func TestMockSpan_SetOperationName(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("")
	span.SetOperationName("x")
	assert.Equal(t, "x", span.(*MockSpan).OperationName)
}

func TestMockSpanContext_Baggage(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("x")
	span.Context().SetBaggageItem("x", "y")
	assert.Equal(t, "y", span.Context().BaggageItem("x"))
	assert.Equal(t, map[string]string{"x": "y"}, span.Context().(*MockSpanContext).GetBaggage())

	baggage := make(map[string]string)
	span.Context().ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return true
	})
	assert.Equal(t, map[string]string{"x": "y"}, baggage)

	span.Context().SetBaggageItem("a", "b")
	baggage = make(map[string]string)
	span.Context().ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return false // exit early
	})
	assert.Equal(t, 2, len(span.Context().(*MockSpanContext).GetBaggage()))
	assert.Equal(t, 1, len(baggage))
}

func TestMockSpan_GetTag(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("x")
	span.SetTag("x", "y")
	assert.Equal(t, "y", span.(*MockSpan).GetTag("x"))
}

func TestMockSpan_GetTags(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("x")
	span.SetTag("x", "y")
	assert.Equal(t, map[string]interface{}{"x": "y"}, span.(*MockSpan).GetTags())
}

func TestMockTracer_GetFinishedSpans_and_Reset(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("x")
	span.SetTag("x", "y")
	span.Finish()
	spans := tracer.GetFinishedSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, map[string]interface{}{"x": "y"}, spans[0].GetTags())

	tracer.Reset()
	spans = tracer.GetFinishedSpans()
	assert.Equal(t, 0, len(spans))
}

func TestMockSpan_Logs(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("x")
	span.LogEvent("x")
	span.LogEventWithPayload("y", "z")
	span.Log(opentracing.LogData{Event: "a"})
	span.FinishWithOptions(opentracing.FinishOptions{
		BulkLogData: []opentracing.LogData{opentracing.LogData{Event: "f"}}})
	spans := tracer.GetFinishedSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, []opentracing.LogData{
		opentracing.LogData{Event: "x"},
		opentracing.LogData{Event: "y", Payload: "z"},
		opentracing.LogData{Event: "a"},
		opentracing.LogData{Event: "f"},
	}, spans[0].GetLogs())
}

func TestMockTracer_Propagation(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("x")
	span.Context().SetBaggageItem("x", "y")

	assert.Equal(t, opentracing.ErrUnsupportedFormat,
		tracer.Inject(span.Context(), opentracing.Binary, nil))
	assert.Equal(t, opentracing.ErrInvalidCarrier,
		tracer.Inject(span.Context(), opentracing.TextMap, span))

	carrier := make(map[string]string)

	err := tracer.Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(carrier))
	require.NoError(t, err)
	t.Logf("%+v", carrier)
	assert.Equal(t, 2, len(carrier), "expect baggage + id")

	_, err = tracer.Extract(opentracing.Binary, nil)
	assert.Equal(t, opentracing.ErrUnsupportedFormat, err)
	_, err = tracer.Extract(opentracing.TextMap, tracer)
	assert.Equal(t, opentracing.ErrInvalidCarrier, err)

	extractedContext, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(carrier))
	require.NoError(t, err)
	assert.Equal(t, "y", extractedContext.BaggageItem("x"))
}
