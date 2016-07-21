package mocktracer

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// New returns a MockTracer opentracing.Tracer implementation that's intended
// to facilitate tests of OpenTracing instrumentation.
func New() *MockTracer {
	t := &MockTracer{
		finishedSpans: []*MockSpan{},
		injectors:     make(map[interface{}]Injector),
		extractors:    make(map[interface{}]Extractor),
	}

	// register default injectors/extractors
	textPropagator := new(TextMapPropagator)
	t.RegisterInjector(opentracing.TextMap, textPropagator)
	t.RegisterExtractor(opentracing.TextMap, textPropagator)

	return t
}

// MockTracer is only intended for testing OpenTracing instrumentation.
//
// It is entirely unsuitable for production use, but appropriate for tests
// that want to verify tracing behavior in other frameworks/applications.
type MockTracer struct {
	sync.RWMutex
	finishedSpans []*MockSpan
	injectors     map[interface{}]Injector
	extractors    map[interface{}]Extractor
}

// Injector is responsible for injecting SpanContext instances in a manner suitable
// for propagation via a format-specific "carrier" object. Typically the
// injection will take place across an RPC boundary, but message queues and
// other IPC mechanisms are also reasonable places to use an Injector.
type Injector interface {
	// Inject takes `SpanContext` and injects it into `carrier`. The actual type
	// of `carrier` depends on the `format` passed to `Tracer.Inject()`.
	//
	// Implementations may return opentracing.ErrInvalidCarrier or any other
	// implementation-specific error if injection fails.
	Inject(ctx *MockSpanContext, carrier interface{}) error
}

// Extractor is responsible for extracting SpanContext instances from a
// format-specific "carrier" object. Typically the extraction will take place
// on the server side of an RPC boundary, but message queues and other IPC
// mechanisms are also reasonable places to use an Extractor.
type Extractor interface {
	// Extract decodes a SpanContext instance from the given `carrier`,
	// or (nil, opentracing.ErrSpanContextNotFound) if no context could
	// be found in the `carrier`.
	Extract(carrier interface{}) (*MockSpanContext, error)
}

// MockSpanContext is an opentracing.SpanContext implementation.
//
// It is entirely unsuitable for production use, but appropriate for tests
// that want to verify tracing behavior in other frameworks/applications.
//
// By default all spans have Sampled=true flag, unless {"sampling.priority": 0}
// tag is set.
type MockSpanContext struct {
	sync.RWMutex

	TraceID int
	SpanID  int
	Sampled bool
	baggage map[string]string
}

// SetBaggageItem belongs to the SpanContext interface
func (s *MockSpanContext) SetBaggageItem(key, val string) opentracing.SpanContext {
	s.Lock()
	defer s.Unlock()
	if s.baggage == nil {
		s.baggage = make(map[string]string)
	}
	s.baggage[key] = val
	return s
}

// BaggageItem belongs to the SpanContext interface
func (s *MockSpanContext) BaggageItem(key string) string {
	s.RLock()
	defer s.RUnlock()
	return s.baggage[key]
}

// ForeachBaggageItem belongs to the SpanContext interface
func (s *MockSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {
	s.RLock()
	defer s.RUnlock()
	for k, v := range s.baggage {
		if !handler(k, v) {
			break
		}
	}
}

// Baggage returns a copy of baggage items in the span
func (s *MockSpanContext) Baggage() map[string]string {
	s.RLock()
	defer s.RUnlock()
	baggage := make(map[string]string)
	for k, v := range s.baggage {
		baggage[k] = v
	}
	return baggage
}

// MockSpan is an opentracing.Span implementation that exports its internal
// state for testing purposes.
type MockSpan struct {
	ParentID int

	OperationName string
	StartTime     time.Time
	FinishTime    time.Time

	// All of the below (including SpanContext) are protected by spanContext's
	// embedded RWMutex.
	SpanContext *MockSpanContext
	tags        map[string]interface{}
	logs        []opentracing.LogData
	tracer      *MockTracer
}

// FinishedSpans returns all spans that have been Finish()'ed since the
// MockTracer was constructed or since the last call to its Reset() method.
func (t *MockTracer) FinishedSpans() []*MockSpan {
	t.RLock()
	defer t.RUnlock()
	spans := make([]*MockSpan, len(t.finishedSpans))
	copy(spans, t.finishedSpans)
	return spans
}

// Reset clears the internally accumulated finished spans. Note that any
// extant MockSpans will still append to finishedSpans when they Finish(),
// even after a call to Reset().
func (t *MockTracer) Reset() {
	t.Lock()
	defer t.Unlock()
	t.finishedSpans = []*MockSpan{}
}

// Tags returns a copy of tags accumulated by the span so far
func (s *MockSpan) Tags() map[string]interface{} {
	s.SpanContext.RLock()
	defer s.SpanContext.RUnlock()
	tags := make(map[string]interface{})
	for k, v := range s.tags {
		tags[k] = v
	}
	return tags
}

// Tag returns a single tag
func (s *MockSpan) Tag(k string) interface{} {
	s.SpanContext.RLock()
	defer s.SpanContext.RUnlock()
	return s.tags[k]
}

// Logs returns a copy of logs accumulated in the span so far
func (s *MockSpan) Logs() []opentracing.LogData {
	s.SpanContext.RLock()
	defer s.SpanContext.RUnlock()
	logs := make([]opentracing.LogData, len(s.logs))
	copy(logs, s.logs)
	return logs
}

// StartSpan belongs to the Tracer interface.
func (t *MockTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	sso := opentracing.StartSpanOptions{}
	for _, o := range opts {
		o.Apply(&sso)
	}
	return newMockSpan(t, operationName, sso)
}

const mockTextMapIdsPrefix = "mockpfx-ids-"
const mockTextMapBaggagePrefix = "mockpfx-baggage-"

// RegisterInjector registers injector for given format
func (t *MockTracer) RegisterInjector(format interface{}, injector Injector) {
	t.injectors[format] = injector
}

// RegisterExtractor registers extractor for given format
func (t *MockTracer) RegisterExtractor(format interface{}, extractor Extractor) {
	t.extractors[format] = extractor
}

// Inject belongs to the Tracer interface.
func (t *MockTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	spanContext, ok := sm.(*MockSpanContext)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}
	injector, ok := t.injectors[format]
	if !ok {
		return opentracing.ErrUnsupportedFormat
	}
	return injector.Inject(spanContext, carrier)
}

// Extract belongs to the Tracer interface.
func (t *MockTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	extractor, ok := t.extractors[format]
	if !ok {
		return nil, opentracing.ErrUnsupportedFormat
	}
	return extractor.Extract(carrier)
}

// TextMapPropagator implements Injector/Extractor for TextMap format.
type TextMapPropagator struct{}

// Inject implements the Injector interface
func (t *TextMapPropagator) Inject(spanContext *MockSpanContext, carrier interface{}) error {
	spanContext.RLock()
	defer spanContext.RUnlock()
	writer, ok := carrier.(opentracing.TextMapWriter)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}
	// Ids:
	writer.Set(mockTextMapIdsPrefix+"traceid", strconv.Itoa(spanContext.TraceID))
	writer.Set(mockTextMapIdsPrefix+"spanid", strconv.Itoa(spanContext.SpanID))
	writer.Set(mockTextMapIdsPrefix+"sampled", fmt.Sprint(spanContext.Sampled))
	// Baggage:
	for baggageKey, baggageVal := range spanContext.baggage {
		writer.Set(mockTextMapBaggagePrefix+baggageKey, baggageVal)
	}
	return nil
}

// Extract implements the Extractor interface
func (t *TextMapPropagator) Extract(carrier interface{}) (*MockSpanContext, error) {
	reader, ok := carrier.(opentracing.TextMapReader)
	if !ok {
		return nil, opentracing.ErrInvalidCarrier
	}
	rval := newMockSpanContext(0, 0, true, nil)
	err := reader.ForeachKey(func(key, val string) error {
		lowerKey := strings.ToLower(key)
		switch {
		case lowerKey == mockTextMapIdsPrefix+"traceid":
			// Ids:
			i, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			rval.TraceID = i
		case lowerKey == mockTextMapIdsPrefix+"spanid":
			// Ids:
			i, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			rval.SpanID = i
		case lowerKey == mockTextMapIdsPrefix+"sampled":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			rval.Sampled = b
		case strings.HasPrefix(lowerKey, mockTextMapBaggagePrefix):
			// Baggage:
			rval.SetBaggageItem(lowerKey[len(mockTextMapBaggagePrefix):], val)
		}
		return nil
	})
	if rval.TraceID == 0 || rval.SpanID == 0 {
		return nil, opentracing.ErrSpanContextNotFound
	}
	if err != nil {
		return nil, err
	}
	return rval, nil
}

var mockIDSource = uint32(42)

func nextMockID() int {
	atomic.AddUint32(&mockIDSource, 1)
	return int(atomic.LoadUint32(&mockIDSource))
}

func newMockSpanContext(traceID int, spanID int, sampled bool, baggage map[string]string) *MockSpanContext {
	baggageCopy := make(map[string]string)
	for k, v := range baggage {
		baggageCopy[k] = v
	}
	return &MockSpanContext{
		TraceID: traceID,
		SpanID:  spanID,
		Sampled: sampled,
		baggage: baggageCopy,
	}
}

func newMockSpan(t *MockTracer, name string, opts opentracing.StartSpanOptions) *MockSpan {
	tags := opts.Tags
	if tags == nil {
		tags = map[string]interface{}{}
	}
	traceID := nextMockID()
	parentID := int(0)
	var baggage map[string]string
	sampled := true
	if len(opts.References) > 0 {
		traceID = opts.References[0].Referee.(*MockSpanContext).TraceID
		parentID = opts.References[0].Referee.(*MockSpanContext).SpanID
		baggage = opts.References[0].Referee.(*MockSpanContext).Baggage()
		sampled = opts.References[0].Referee.(*MockSpanContext).Sampled
	}
	spanContext := newMockSpanContext(traceID, nextMockID(), sampled, baggage)
	startTime := opts.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}
	return &MockSpan{
		ParentID:      parentID,
		OperationName: name,
		StartTime:     startTime,
		tags:          tags,
		logs:          []opentracing.LogData{},
		SpanContext:   spanContext,

		tracer: t,
	}
}

// Context belongs to the Span interface
func (s *MockSpan) Context() opentracing.SpanContext {
	return s.SpanContext
}

// SetTag belongs to the Span interface
func (s *MockSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.SpanContext.Lock()
	defer s.SpanContext.Unlock()
	if key == string(ext.SamplingPriority) {
		if v, ok := value.(uint16); ok {
			s.SpanContext.Sampled = v > 0
			return s
		}
		if v, ok := value.(int); ok {
			s.SpanContext.Sampled = v > 0
			return s
		}
	}
	s.tags[key] = value
	return s
}

// Finish belongs to the Span interface
func (s *MockSpan) Finish() {
	s.SpanContext.Lock()
	s.FinishTime = time.Now()
	s.SpanContext.Unlock()
	s.tracer.recordSpan(s)
}

// FinishWithOptions belongs to the Span interface
func (s *MockSpan) FinishWithOptions(opts opentracing.FinishOptions) {
	s.SpanContext.Lock()
	s.FinishTime = opts.FinishTime
	s.logs = append(s.logs, opts.BulkLogData...)
	s.SpanContext.Unlock()
	s.tracer.recordSpan(s)
}

// String allows printing span for debugging
func (s *MockSpan) String() string {
	return fmt.Sprintf(
		"traceId=%d, spanId=%d, parentId=%d, sampled=%t, name=%s",
		s.SpanContext.TraceID, s.SpanContext.SpanID, s.ParentID,
		s.SpanContext.Sampled, s.OperationName)
}

func (t *MockTracer) recordSpan(span *MockSpan) {
	t.Lock()
	defer t.Unlock()
	t.finishedSpans = append(t.finishedSpans, span)
}

// LogEvent belongs to the Span interface
func (s *MockSpan) LogEvent(event string) {
	s.Log(opentracing.LogData{
		Event: event,
	})
}

// LogEventWithPayload belongs to the Span interface
func (s *MockSpan) LogEventWithPayload(event string, payload interface{}) {
	s.Log(opentracing.LogData{
		Event:   event,
		Payload: payload,
	})
}

// Log belongs to the Span interface
func (s *MockSpan) Log(data opentracing.LogData) {
	s.SpanContext.Lock()
	defer s.SpanContext.Unlock()
	s.logs = append(s.logs, data)
}

// SetOperationName belongs to the Span interface
func (s *MockSpan) SetOperationName(operationName string) opentracing.Span {
	s.SpanContext.Lock()
	defer s.SpanContext.Unlock()
	s.OperationName = operationName
	return s
}

// Tracer belongs to the Span interface
func (s *MockSpan) Tracer() opentracing.Tracer {
	return s.tracer
}
