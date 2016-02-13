package standardtracer_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/standardtracer"
	"github.com/opentracing/opentracing-go/testutils"
)

func TestSpanPropagator(t *testing.T) {
	var err error
	const op = "test"
	recorder := testutils.NewInMemoryRecorder()
	tracer := standardtracer.New(recorder)

	sp := tracer.StartSpan(op)
	sp.SetTraceAttribute("foo", "bar")

	textCarrier := opentracing.NewSplitTextCarrier()
	err = tracer.Injector(opentracing.SplitText).InjectSpan(sp, textCarrier)
	if err != nil {
		t.Fatal(err)
	}
	binaryCarrier := opentracing.NewSplitBinaryCarrier()
	err = tracer.Injector(opentracing.SplitBinary).InjectSpan(sp, binaryCarrier)
	if err != nil {
		t.Fatal(err)
	}

	sp1, err := tracer.Extractor(opentracing.SplitText).JoinTrace(op, textCarrier)
	if err != nil {
		t.Fatal(err)
	}
	sp2, err := tracer.Extractor(opentracing.SplitBinary).JoinTrace(op, binaryCarrier)
	if err != nil {
		t.Fatal(err)
	}
	sp.Finish()
	for _, sp := range []opentracing.Span{sp1, sp2} {
		sp.Finish()
	}

	spans := recorder.GetSpans()
	if a, e := len(spans), 3; a != e {
		t.Fatalf("expected %d spans, got %d", e, a)
	}

	exp := spans[0]
	exp.Duration = time.Duration(123)
	exp.Start = time.Time{}.Add(1)

	for i, sp := range spans[1:] {
		if a, e := sp.ParentSpanID, exp.SpanID; a != e {
			t.Errorf("%d: ParentSpanID %d does not match expectation %d", i, a, e)
		} else {
			// Prepare for comparison.
			sp.SpanID, sp.ParentSpanID = exp.SpanID, 0
			sp.Duration, sp.Start = exp.Duration, exp.Start
		}
		if a, e := sp.TraceID, exp.TraceID; a != e {
			t.Errorf("%d: TraceID changed from %d to %d", i, e, a)
		}
		if !reflect.DeepEqual(exp, sp) {
			t.Errorf("%d: wanted %+v, got %+v", i, spew.Sdump(exp), spew.Sdump(sp))
		}
	}
}
