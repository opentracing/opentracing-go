package harness

import (
	"bytes"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/assert"
)

type CheckOpts struct {
	CheckBaggageValues bool
}

func CheckStartSpan(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan(
		"Fry",
		opentracing.Tag{Key: "birthday", Value: "August 14 1974"})
	span.LogFields(
		log.String("hospital", "Brooklyn Pre-Med Hospital"),
		log.String("city", "Old New York"))
	span.Finish()
}

func CheckStartSpanWithParent(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	parentSpan := tracer.StartSpan("parent")
	assert.NotNil(t, parentSpan)

	span := tracer.StartSpan(
		"Leela",
		opentracing.ChildOf(parentSpan.Context()))
	span.Finish()

	span = tracer.StartSpan(
		"Leela",
		opentracing.FollowsFrom(parentSpan.Context()),
		opentracing.Tag{Key: "birthplace", Value: "sewers"})
	span.Finish()

	parentSpan.Finish()
}

func CheckSetOperationName(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("").SetOperationName("Farnsworth")
	span.Finish()
}

func CheckSpanTagValueTypes(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("ManyTypes")
	span.
		SetTag("an_int", 9).
		SetTag("a_bool", true).
		SetTag("a_string", "aoeuidhtns")
}

func CheckSpanTagsWithChaining(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("Farnsworth")
	span.
		SetTag("birthday", "9 April, 2841").
		SetTag("loves", "different lengths of wires")
	span.
		SetTag("unicode_val", "non-ascii: \u200b").
		SetTag("unicode_key_\u200b", "ascii val")
	span.Finish()
}

func CheckSpanLogs(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("Fry")
	span.LogKV(
		"frozen.year", 1999,
		"frozen.place", "Cryogenics Labs")
	span.LogKV(
		"defrosted.year", 2999,
		"defrosted.place", "Cryogenics Labs")

	// XXX add LogFields
	// XXX add LogRecords FinishOptions with timestamp
	span.Finish()
}

func assertEmptyBaggage(t *testing.T, spanContext opentracing.SpanContext) {
	if !assert.NotNil(t, spanContext, "assertEmptyBaggage got empty context") {
		return
	}
	spanContext.ForeachBaggageItem(func(k, v string) bool {
		assert.Fail(t, "new span shouldn't have baggage")
		return false
	})
}

func CheckSpanBaggage(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("Fry")
	assertEmptyBaggage(t, span.Context())

	spanRef := span.SetBaggageItem("Kiff-loves", "Amy")
	assert.Exactly(t, spanRef, span)

	val := span.BaggageItem("Kiff-loves")
	if opts.CheckBaggageValues {
		assert.Equal(t, "Amy", val)
	}
	span.Finish()
}

func CheckContextBaggage(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("Fry")
	assertEmptyBaggage(t, span.Context())

	span.SetBaggageItem("Kiff-loves", "Amy")
	if opts.CheckBaggageValues {
		called := false
		span.Context().ForeachBaggageItem(func(k, v string) bool {
			assert.False(t, called)
			called = true
			assert.Equal(t, "Kiff-loves", k)
			assert.Equal(t, "Amy", v)
			return true
		})
	}
	span.Finish()
}

func CheckTextPropagation(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("Bender")
	textCarrier := opentracing.TextMapCarrier{}
	err := span.Tracer().Inject(span.Context(), opentracing.TextMap, textCarrier)
	assert.NoError(t, err)

	extractedContext, err := tracer.Extract(opentracing.TextMap, textCarrier)
	assert.NoError(t, err)
	assertEmptyBaggage(t, extractedContext)
	// XXX add option to check if propagation "works"
	span.Finish()
}

func CheckHTTPPropagation(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("Bender")
	textCarrier := opentracing.HTTPHeadersCarrier{}
	// XXX add same test cases around valid HTTP header characters, casing
	err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, textCarrier)
	assert.NoError(t, err)

	extractedContext, err := tracer.Extract(opentracing.HTTPHeaders, textCarrier)
	assert.NoError(t, err)
	assertEmptyBaggage(t, extractedContext)
	// XXX add option to check if propagation "works"
	span.Finish()
}

func CheckBinaryPropagation(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	span := tracer.StartSpan("Bender")
	buf := new(bytes.Buffer)
	err := span.Tracer().Inject(span.Context(), opentracing.Binary, buf)
	assert.NoError(t, err)

	extractedContext, err := tracer.Extract(opentracing.Binary, buf)
	assert.NoError(t, err)
	assertEmptyBaggage(t, extractedContext)
	// XXX add option to check if propagation "works"
	span.Finish()
}

func CheckMandatoryFormats(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	formats := []struct{ Format, Carrier interface{} }{
		{opentracing.TextMap, opentracing.TextMapCarrier{}},
		{opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier{}},
		{opentracing.Binary, new(bytes.Buffer)},
	}
	span := tracer.StartSpan("Bender")
	for _, fmtCarrier := range formats {
		err := span.Tracer().Inject(span.Context(), fmtCarrier.Format, fmtCarrier.Carrier)
		assert.NoError(t, err)
		spanCtx, err := tracer.Extract(fmtCarrier.Format, fmtCarrier.Carrier)
		assert.NoError(t, err)
		assertEmptyBaggage(t, spanCtx)
	}
}

func CheckUnknownFormat(t *testing.T, tracer opentracing.Tracer, opts CheckOpts) {
	customFormat := "kiss my shiny metal ..."
	span := tracer.StartSpan("Bender")

	err := span.Tracer().Inject(span.Context(), customFormat, nil)
	assert.Equal(t, opentracing.ErrUnsupportedFormat, err)

	ctx, err := tracer.Extract(customFormat, nil)
	assert.Nil(t, ctx)
	assert.Equal(t, opentracing.ErrUnsupportedFormat, err)
}
