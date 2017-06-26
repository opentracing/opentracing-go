// Package harness provides a suite of API compatibility checks. They were originally
// ported from the OpenTracing Python library's "harness" module.
package harness

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// APICheckCapabilities describes options used by APICheckSuite when testing a Tracer.
type APICheckCapabilities struct {
	CheckBaggageValues bool          // whether to check for propagation of baggage values
	CheckExtract       bool          // whether to check if extracting contexts from carriers works
	CheckInject        bool          // whether to check if injecting contexts works
	Probe              APICheckProbe // optional interface providing methods to check recorded data
}

// APICheckProbe exposes methods for testing data recorded by a Tracer.
type APICheckProbe interface {
	// SameTrace helps tests assert that this tracer's spans are from the same trace.
	SameTrace(first, second opentracing.Span) bool
	// SameSpanContext helps tests assert that a span and a context are from the same trace and span.
	SameSpanContext(opentracing.Span, opentracing.SpanContext) bool
}

// APICheckSuite is a testify suite for checking a Tracer against the OpenTracing API.
type APICheckSuite struct {
	suite.Suite
	opts      APICheckCapabilities
	newTracer func() (tracer opentracing.Tracer, closer func())
	tracer    opentracing.Tracer
	closer    func()
}

// NewAPICheckSuite returns a testify suite for checking a Tracer against the OpenTracing API.
// It is provided a function that will be executed to create and destroy a tracer for each test
// in the suite, and API test options described by APICheckCapabilities.
func NewAPICheckSuite(
	newTracer func() (tracer opentracing.Tracer, closer func()),
	opts ...APICheckOption,
) *APICheckSuite {
	s := &APICheckSuite{newTracer: newTracer}
	for _, o := range opts {
		o.Apply(s)
	}
	return s
}

// APICheckOption instances may be passed to NewAPICheckSuite.
type APICheckOption interface {
	Apply(*APICheckSuite)
}

// Apply satisfies the APICheckOption interface.
func (c APICheckCapabilities) Apply(s *APICheckSuite) {
	s.opts = c
}

// CheckEverything enables all API checks.
type CheckEverything struct{}

// Apply satisfies the APICheckOption interface.
func (CheckEverything) Apply(s *APICheckSuite) {
	s.opts.CheckBaggageValues = true
	s.opts.CheckExtract = true
	s.opts.CheckInject = true
}

// UseProbe specifies an APICheckProbe implementation to use.
type UseProbe struct {
	APICheckProbe
}

// Apply satisfies the APICheckOption interface.
func (u UseProbe) Apply(s *APICheckSuite) {
	s.opts.Probe = u.APICheckProbe
}

// BeforeTest creates a tracer for this specific test invocation.
func (s *APICheckSuite) BeforeTest(suiteName, testName string) {
	s.tracer, s.closer = s.newTracer()
	if s.tracer == nil {
		panic(fmt.Sprintf("newTracer returned nil Tracer before running %s, %s", suiteName, testName))
	}
}

// AfterTest closes the tracer, and clears the test-specific tracer.
func (s *APICheckSuite) AfterTest(suiteName, testName string) {
	if s.closer != nil {
		s.closer()
	}
	s.tracer, s.closer = nil, nil
}

// TestStartSpan checks if a Tracer can start a span and calls some span API methods.
func (s *APICheckSuite) TestStartSpan() {
	span := s.tracer.StartSpan(
		"Fry",
		opentracing.Tag{Key: "birthday", Value: "August 14 1974"})
	span.LogFields(
		log.String("hospital", "Brooklyn Pre-Med Hospital"),
		log.String("city", "Old New York"))
	span.Finish()
}

// TestStartSpanWithParent checks if a Tracer can start a span with a specified parent.
func (s *APICheckSuite) TestStartSpanWithParent() {
	parentSpan := s.tracer.StartSpan("parent")
	s.NotNil(parentSpan)

	span := s.tracer.StartSpan(
		"Leela",
		opentracing.ChildOf(parentSpan.Context()))
	span.Finish()
	if s.opts.Probe != nil {
		s.True(s.opts.Probe.SameTrace(parentSpan, span))
	} else {
		s.T().Log("harness.Probe not specified, skipping")
	}

	span = s.tracer.StartSpan(
		"Leela",
		opentracing.FollowsFrom(parentSpan.Context()),
		opentracing.Tag{Key: "birthplace", Value: "sewers"})
	if s.opts.Probe != nil {
		s.True(s.opts.Probe.SameTrace(parentSpan, span))
	} else {
		s.T().Log("harness.Probe not specified, skipping")
	}
	span.Finish()

	parentSpan.Finish()
}

// TestSetOperationName attempts to set the operation name on a span after it has been created.
func (s *APICheckSuite) TestSetOperationName() {
	span := s.tracer.StartSpan("").SetOperationName("Farnsworth")
	span.Finish()
}

// TestSpanTagValueTypes sets tags using values of different types.
func (s *APICheckSuite) TestSpanTagValueTypes() {
	span := s.tracer.StartSpan("ManyTypes")
	span.
		SetTag("an_int", 9).
		SetTag("a_bool", true).
		SetTag("a_string", "aoeuidhtns")
}

// TestSpanTagsWithChaining tests chaining of calls to SetTag.
func (s *APICheckSuite) TestSpanTagsWithChaining() {
	span := s.tracer.StartSpan("Farnsworth")
	span.
		SetTag("birthday", "9 April, 2841").
		SetTag("loves", "different lengths of wires")
	span.
		SetTag("unicode_val", "non-ascii: \u200b").
		SetTag("unicode_key_\u200b", "ascii val")
	span.Finish()
}

// TestSpanLogs tests calls to log keys and values with spans.
func (s *APICheckSuite) TestSpanLogs() {
	span := s.tracer.StartSpan("Fry")
	span.LogKV(
		"event", "frozen",
		"year", 1999,
		"place", "Cryogenics Labs")
	span.LogKV(
		"event", "defrosted",
		"year", 2999,
		"place", "Cryogenics Labs")

	ts := time.Now()
	span.FinishWithOptions(opentracing.FinishOptions{
		LogRecords: []opentracing.LogRecord{
			{
				Timestamp: ts,
				Fields: []log.Field{
					log.String("event", "job-assignment"),
					log.String("type", "delivery boy"),
				},
			},
		}})

	// Test deprecated log methods
	span.LogEvent("an arbitrary event")
	span.LogEventWithPayload("y", "z")
	span.Log(opentracing.LogData{Event: "y", Payload: "z"})
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

// TestSpanBaggage tests calls to set and get span baggage, and if the CheckBaggageValues option
// is set, asserts that baggage values were successfully retrieved.
func (s *APICheckSuite) TestSpanBaggage() {
	span := s.tracer.StartSpan("Fry")
	assertEmptyBaggage(s.T(), span.Context())

	spanRef := span.SetBaggageItem("Kiff-loves", "Amy")
	s.Exactly(spanRef, span)

	val := span.BaggageItem("Kiff-loves")
	if s.opts.CheckBaggageValues {
		s.Equal("Amy", val)
	} else {
		s.T().Log("CheckBaggageValues capability not set, skipping")
	}
	span.Finish()
}

// TestContextBaggage tests calls to set and get span baggage, and if the CheckBaggageValues option
// is set, asserts that baggage values were successfully retrieved from the span's SpanContext.
func (s *APICheckSuite) TestContextBaggage() {
	span := s.tracer.StartSpan("Fry")
	assertEmptyBaggage(s.T(), span.Context())

	span.SetBaggageItem("Kiff-loves", "Amy")
	if s.opts.CheckBaggageValues {
		called := false
		span.Context().ForeachBaggageItem(func(k, v string) bool {
			s.False(called)
			called = true
			s.Equal("Kiff-loves", k)
			s.Equal("Amy", v)
			return true
		})
	} else {
		s.T().Log("CheckBaggageValues capability not set, skipping")
	}
	span.Finish()
}

// TestTextPropagation tests if the Tracer can Inject a span into a TextMapCarrier, and later Extract it.
// If CheckExtract is set, it will check if Extract was successful (returned no error). If a Probe is set,
// it will check if the extracted context is in the same trace as the original span.
func (s *APICheckSuite) TestTextPropagation() {
	span := s.tracer.StartSpan("Bender")
	textCarrier := opentracing.TextMapCarrier{}
	err := span.Tracer().Inject(span.Context(), opentracing.TextMap, textCarrier)
	assert.NoError(s.T(), err)

	extractedContext, err := s.tracer.Extract(opentracing.TextMap, textCarrier)
	if s.opts.CheckExtract {
		s.NoError(err)
		assertEmptyBaggage(s.T(), extractedContext)
	} else {
		s.T().Log("CheckExtract capability not set, skipping")
	}
	if s.opts.Probe != nil {
		s.True(s.opts.Probe.SameSpanContext(span, extractedContext))
	} else {
		s.T().Log("harness.Probe not specified, skipping")
	}
	span.Finish()
}

// TestHTTPPropagation tests if the Tracer can Inject a span into HTTP headers, and later Extract it.
// If CheckExtract is set, it will check if Extract was successful (returned no error). If a Probe is set,
// it will check if the extracted context is in the same trace as the original span.
func (s *APICheckSuite) TestHTTPPropagation() {
	span := s.tracer.StartSpan("Bender")
	textCarrier := opentracing.HTTPHeadersCarrier{}
	err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, textCarrier)
	s.NoError(err)

	extractedContext, err := s.tracer.Extract(opentracing.HTTPHeaders, textCarrier)
	if s.opts.CheckExtract {
		s.NoError(err)
		assertEmptyBaggage(s.T(), extractedContext)
	} else {
		s.T().Log("CheckExtract capability not set, skipping")
	}
	if s.opts.Probe != nil {
		s.True(s.opts.Probe.SameSpanContext(span, extractedContext))
	} else {
		s.T().Log("harness.Probe not specified, skipping")
	}
	span.Finish()
}

// TestBinaryPropagation tests if the Tracer can Inject a span into a binary buffer, and later Extract it.
// If CheckExtract is set, it will check if Extract was successful (returned no error). If a Probe is set,
// it will check if the extracted context is in the same trace as the original span.
func (s *APICheckSuite) TestBinaryPropagation() {
	span := s.tracer.StartSpan("Bender")
	buf := new(bytes.Buffer)
	err := span.Tracer().Inject(span.Context(), opentracing.Binary, buf)
	s.NoError(err)

	extractedContext, err := s.tracer.Extract(opentracing.Binary, buf)
	if s.opts.CheckExtract {
		s.NoError(err)
		assertEmptyBaggage(s.T(), extractedContext)
	} else {
		s.T().Log("CheckExtract capability not set, skipping")
	}
	if s.opts.Probe != nil {
		s.True(s.opts.Probe.SameSpanContext(span, extractedContext))
	} else {
		s.T().Log("harness.Probe not specified, skipping")
	}
	span.Finish()
}

// TestMandatoryFormats tests if all mandatory carrier formats are supported. If CheckExtract is set, it
// will check if the call to Extract was successful (returned no error such as ErrUnsupportedFormat).
func (s *APICheckSuite) TestMandatoryFormats() {
	formats := []struct{ Format, Carrier interface{} }{
		{opentracing.TextMap, opentracing.TextMapCarrier{}},
		{opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier{}},
		{opentracing.Binary, new(bytes.Buffer)},
	}
	span := s.tracer.StartSpan("Bender")
	for _, fmtCarrier := range formats {
		err := span.Tracer().Inject(span.Context(), fmtCarrier.Format, fmtCarrier.Carrier)
		s.NoError(err)
		spanCtx, err := s.tracer.Extract(fmtCarrier.Format, fmtCarrier.Carrier)
		if s.opts.CheckExtract {
			s.NoError(err)
			assertEmptyBaggage(s.T(), spanCtx)
		} else {
			s.T().Log("CheckExtract capability not set, skipping")
		}
	}
}

// TestUnknownFormat checks if attempting to Inject or Extract using an unsupported format
// returns ErrUnsupportedFormat, if CheckInject and CheckExtract are set.
func (s *APICheckSuite) TestUnknownFormat() {
	customFormat := "kiss my shiny metal ..."
	span := s.tracer.StartSpan("Bender")

	err := span.Tracer().Inject(span.Context(), customFormat, nil)
	if s.opts.CheckInject {
		s.Equal(opentracing.ErrUnsupportedFormat, err)
	} else {
		s.T().Log("CheckInject capability not set, skipping")
	}
	ctx, err := s.tracer.Extract(customFormat, nil)
	s.Nil(ctx)
	if s.opts.CheckExtract {
		s.Equal(opentracing.ErrUnsupportedFormat, err)
	} else {
		s.T().Log("CheckExtract capability not set, skipping")
	}
}

// ForeignSpanContext satisfies the opentracing.SpanContext interface, but otherwise does nothing.
type ForeignSpanContext struct{}

// ForeachBaggageItem could call handler for each baggage KV, but does nothing.
func (f ForeignSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}

// NotACarrier does not satisfy any of the opentracing carrier interfaces.
type NotACarrier struct{}

// TestInvalidInject checks if errors are returned when Inject is called with invalid inputs.
func (s *APICheckSuite) TestInvalidInject() {
	if !s.opts.CheckInject {
		s.T().Skip("CheckInject capability not set, skipping")
	}
	span := s.tracer.StartSpan("op")

	// binary inject
	err := span.Tracer().Inject(ForeignSpanContext{}, opentracing.Binary, new(bytes.Buffer))
	s.Equal(opentracing.ErrInvalidSpanContext, err, "Foreign SpanContext should return invalid error")
	err = span.Tracer().Inject(span.Context(), opentracing.Binary, NotACarrier{})
	s.Equal(opentracing.ErrInvalidCarrier, err, "Carrier that's not io.Writer should return error")

	// text inject
	err = span.Tracer().Inject(ForeignSpanContext{}, opentracing.TextMap, opentracing.TextMapCarrier{})
	s.Equal(opentracing.ErrInvalidSpanContext, err, "Foreign SpanContext should return invalid error")
	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, NotACarrier{})
	s.Equal(opentracing.ErrInvalidCarrier, err, "Carrier that's not TextMapWriter should return error")

	// HTTP inject
	err = span.Tracer().Inject(ForeignSpanContext{}, opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier{})
	s.Equal(opentracing.ErrInvalidSpanContext, err, "Foreign SpanContext should return invalid error")
	err = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, NotACarrier{})
	s.Equal(opentracing.ErrInvalidCarrier, err, "Carrier that's not TextMapWriter should return error")
}

// TestInvalidExtract checks if errors are returned when Extract is called with invalid inputs.
func (s *APICheckSuite) TestInvalidExtract() {
	if !s.opts.CheckExtract {
		s.T().Skip("CheckExtract capability not set, skipping")
	}
	span := s.tracer.StartSpan("op")

	// binary extract
	ctx, err := span.Tracer().Extract(opentracing.Binary, NotACarrier{})
	s.Equal(opentracing.ErrInvalidCarrier, err, "Carrier that's not io.Reader should return error")
	s.Nil(ctx)

	// text extract
	ctx, err = span.Tracer().Extract(opentracing.TextMap, NotACarrier{})
	s.Equal(opentracing.ErrInvalidCarrier, err, "Carrier that's not TextMapReader should return error")
	s.Nil(ctx)

	// HTTP extract
	ctx, err = span.Tracer().Extract(opentracing.HTTPHeaders, NotACarrier{})
	s.Equal(opentracing.ErrInvalidCarrier, err, "Carrier that's not TextMapReader should return error")
	s.Nil(ctx)

	span.Finish()
}

// TestMultiBaggage tests calls to set multiple baggage items, and if the CheckBaggageValues option
// is set, asserts that a baggage value was successfully retrieved from the span's SpanContext.
// It also ensures that returning false from the ForeachBaggageItem handler aborts iteration.
func (s *APICheckSuite) TestMultiBaggage() {
	span := s.tracer.StartSpan("op")
	assertEmptyBaggage(s.T(), span.Context())

	span.SetBaggageItem("Bag1", "BaggageVal1")
	span.SetBaggageItem("Bag2", "BaggageVal2")
	if s.opts.CheckBaggageValues {
		s.Equal("BaggageVal1", span.BaggageItem("Bag1"))
		s.Equal("BaggageVal2", span.BaggageItem("Bag2"))
		called := false
		span.Context().ForeachBaggageItem(func(k, v string) bool {
			s.False(called) // should only be called once
			called = true
			return false
		})
		s.True(called)
	} else {
		s.T().Log("CheckBaggageValues capability not set, skipping")
	}
	span.Finish()
}

// TestHTTPBaggage tests whether serializing multiple baggage items that use the same key name
// with different casing will overwrite each other (since HTTP Header names are case-insensitive).
func (s *APICheckSuite) TestHTTPBaggage() {
	if !s.opts.CheckBaggageValues {
		s.T().Skip("CheckBaggageValues capability not set, skipping")
	}
	// set baggage items: should result in two baggage items when serializing using HTTP headers
	span := s.tracer.StartSpan("Bender")
	span.SetBaggageItem("BagItem", "Val1")
	span.SetBaggageItem("bagitem", "Val2")
	span.SetBaggageItem("BAGITEM", "Val3")
	span.SetBaggageItem("KeyName", "Val4")

	carrier := opentracing.HTTPHeadersCarrier{}
	err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	s.NoError(err)

	extractedContext, err := s.tracer.Extract(opentracing.HTTPHeaders, carrier)
	if s.opts.CheckExtract {
		s.NoError(err)
	} else {
		s.T().Log("CheckExtract capability not set, skipping")
	}
	// check baggage items
	count := 0
	items := make(map[string]string)
	extractedContext.ForeachBaggageItem(func(k, v string) bool {
		count++
		items[k] = v
		return true
	})
	s.Len(items, 2, "HTTP Baggage propagation seems to depend on case-sensitive HTTP header field names")
	s.Equal(count, 2, "HTTP Baggage propagation seems to depend on case-sensitive HTTP header field names")
}
