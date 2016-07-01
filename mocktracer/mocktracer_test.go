package mocktracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/opentracing/opentracing-go"
)

func TestMockTracer_StartSpanWithOptions(t *testing.T) {
	tracer := New()
	span1 := tracer.StartSpanWithOptions(opentracing.StartSpanOptions{
		OperationName: "a",
		Tags:          map[string]interface{}{"x": "y"}})

	span2 := span1.Tracer().StartSpanWithOptions(opentracing.StartSpanOptions{Parent: span1})
	span2.Finish()
	span1.Finish()
	spans := tracer.GetFinishedSpans()
	assert.Equal(t, 2, len(spans))

	parent := spans[1]
	child := spans[0]
	assert.Equal(t, map[string]interface{}{"x": "y"}, parent.GetTags())
	assert.Equal(t, child.ParentID, parent.SpanID)
}

func TestMockSpan_SetOperationName(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("")
	span.SetOperationName("x")
	assert.Equal(t, "x", span.(*MockSpan).OperationName)
}

func TestMockSpan_Baggage(t *testing.T) {
	tracer := New()
	span := tracer.StartSpan("x")
	span.SetBaggageItem("x", "y")
	assert.Equal(t, "y", span.BaggageItem("x"))
	assert.Equal(t, map[string]string{"x": "y"}, span.(*MockSpan).GetBaggage())

	baggage := make(map[string]string)
	span.ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return true
	})
	assert.Equal(t, map[string]string{"x": "y"}, baggage)

	span.SetBaggageItem("a", "b")
	baggage = make(map[string]string)
	span.ForeachBaggageItem(func(k, v string) bool {
		baggage[k] = v
		return false // exit early
	})
	assert.Equal(t, 2, len(span.(*MockSpan).GetBaggage()))
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
	span.SetBaggageItem("x", "y")

	assert.Equal(t, opentracing.ErrUnsupportedFormat,
		tracer.Inject(span, opentracing.Binary, nil))
	assert.Equal(t, opentracing.ErrInvalidCarrier,
		tracer.Inject(span, opentracing.TextMap, span))

	carrier := make(map[string]string)

	err := tracer.Inject(span, opentracing.TextMap, opentracing.TextMapCarrier(carrier))
	require.NoError(t, err)
	t.Logf("%+v", carrier)
	assert.Equal(t, 2, len(carrier), "expect baggage + id")

	_, err = tracer.Join("y", opentracing.Binary, nil)
	assert.Equal(t, opentracing.ErrUnsupportedFormat, err)
	_, err = tracer.Join("y", opentracing.TextMap, tracer)
	assert.Equal(t, opentracing.ErrInvalidCarrier, err)

	span2, err := tracer.Join("y", opentracing.TextMap, opentracing.TextMapCarrier(carrier))
	require.NoError(t, err)
	assert.Equal(t, "y", span2.BaggageItem("x"))
}
