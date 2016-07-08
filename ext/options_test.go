package ext_test

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestCarriedOption(t *testing.T) {
	tracer := mocktracer.New()
	parent := tracer.StartSpan("my-trace")

	carrier := opentracing.HTTPHeaderTextMapCarrier{}
	err := tracer.Inject(parent.Context(), opentracing.TextMap, carrier)
	if err != nil {
		t.Fatal(err)
	}

	tracer.StartSpan("my-child", ext.CarriedChildOfOption(opentracing.TextMap, carrier)).Finish()

	rawSpan := tracer.GetFinishedSpans()[0]
	assertEqual(t, rawSpan.ParentID, parent.Context().(*mocktracer.MockSpanContext).SpanID)
	assertEqual(t, rawSpan.ParentRelationship, opentracing.ChildOfRef)
}
