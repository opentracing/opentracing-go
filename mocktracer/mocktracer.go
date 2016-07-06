package mocktracer

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
)

// New returns a MockTracer opentracing.Tracer implementation that's intended
// to facilitate tests of OpenTracing instrumentation.
func New() *MockTracer {
	return &MockTracer{
		finishedSpans: []*MockSpan{},
	}
}

// MockTracer is only intended for testing OpenTracing instrumentation.
//
// It is entirely unsuitable for production use, but appropriate for tests
// that want to verify tracing behavior in other frameworks/applications.
type MockTracer struct {
	sync.RWMutex
	finishedSpans []*MockSpan
}

// MockSpanContext is an opentracing.SpanContext implementation.
//
// It is entirely unsuitable for production use, but appropriate for tests
// that want to verify tracing behavior in other frameworks/applications.
type MockSpanContext struct {
	sync.RWMutex

	SpanID  int
	baggage map[string]string
}

// SetBaggageItem belongs to the SpanContext interface
func (s *MockSpanContext) SetBaggageItem(key, val string) opentracing.SpanContext {
	s.Lock()
	defer s.Unlock()
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

// GetBaggage returns a copy of baggage items in the span
func (s *MockSpanContext) GetBaggage() map[string]string {
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

	// All of the below (including spanContext) are protected by spanContext's
	// embedded RWMutex.
	spanContext *MockSpanContext
	tags        map[string]interface{}
	logs        []opentracing.LogData
	tracer      *MockTracer
}

// GetFinishedSpans returns all spans that have been Finish()'ed since the
// MockTracer was constructed or since the last call to its Reset() method.
func (t *MockTracer) GetFinishedSpans() []*MockSpan {
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

// GetTags returns a copy of tags accumulated by the span so far
func (s *MockSpan) GetTags() map[string]interface{} {
	s.spanContext.RLock()
	defer s.spanContext.RUnlock()
	tags := make(map[string]interface{})
	for k, v := range s.tags {
		tags[k] = v
	}
	return tags
}

// GetTag returns a single tag
func (s *MockSpan) GetTag(k string) interface{} {
	s.spanContext.RLock()
	defer s.spanContext.RUnlock()
	return s.tags[k]
}

// GetLogs returns a copy of logs accumulated in the span so far
func (s *MockSpan) GetLogs() []opentracing.LogData {
	s.spanContext.RLock()
	defer s.spanContext.RUnlock()
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

// Inject belongs to the Tracer interface.
func (t *MockTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	spanContext := sm.(*MockSpanContext)
	spanContext.RLock()
	defer spanContext.RUnlock()
	switch format {
	case opentracing.TextMap:
		writer, ok := carrier.(opentracing.TextMapWriter)
		if !ok {
			return opentracing.ErrInvalidCarrier
		}
		// Ids:
		writer.Set(mockTextMapIdsPrefix+"spanid", strconv.Itoa(spanContext.SpanID))
		// Baggage:
		for baggageKey, baggageVal := range spanContext.baggage {
			writer.Set(mockTextMapBaggagePrefix+baggageKey, baggageVal)
		}
		return nil
	}
	return opentracing.ErrUnsupportedFormat
}

// Extract belongs to the Tracer interface.
func (t *MockTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	switch format {
	case opentracing.TextMap:
		rval := newMockSpanContext(0)
		reader, ok := carrier.(opentracing.TextMapReader)
		if !ok {
			return nil, opentracing.ErrInvalidCarrier
		}
		err := reader.ForeachKey(func(key, val string) error {
			lowerKey := strings.ToLower(key)
			switch {
			case lowerKey == mockTextMapIdsPrefix+"spanid":
				// Ids:
				i, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				rval.SpanID = i
			case strings.HasPrefix(lowerKey, mockTextMapBaggagePrefix):
				// Baggage:
				rval.baggage[lowerKey[len(mockTextMapBaggagePrefix):]] = val
			}
			return nil
		})
		if rval.SpanID == 0 {
			return nil, opentracing.ErrSpanContextNotFound
		}
		return rval, err
	}
	return nil, opentracing.ErrUnsupportedFormat
}

var mockIDSource = uint32(42)

func nextMockID() int {
	atomic.AddUint32(&mockIDSource, 1)
	return int(atomic.LoadUint32(&mockIDSource))
}

func newMockSpanContext(spanID int) *MockSpanContext {
	return &MockSpanContext{
		SpanID:  spanID,
		baggage: make(map[string]string),
	}
}

func newMockSpan(t *MockTracer, name string, opts opentracing.StartSpanOptions) *MockSpan {
	tags := opts.Tags
	if tags == nil {
		tags = map[string]interface{}{}
	}
	spanContext := newMockSpanContext(nextMockID())
	parentID := int(0)
	if len(opts.References) > 0 {
		parentID = opts.References[0].Referee.(*MockSpanContext).SpanID
		baggage := opts.References[0].Referee.(*MockSpanContext).GetBaggage()
		spanContext.baggage = baggage
	}
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
		spanContext:   spanContext,

		tracer: t,
	}
}

// Context belongs to the Span interface
func (s *MockSpan) Context() opentracing.SpanContext {
	return s.spanContext
}

// SetTag belongs to the Span interface
func (s *MockSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.spanContext.Lock()
	defer s.spanContext.Unlock()
	s.tags[key] = value
	return s
}

// Finish belongs to the Span interface
func (s *MockSpan) Finish() {
	s.spanContext.Lock()
	s.FinishTime = time.Now()
	s.spanContext.Unlock()
	s.tracer.recordSpan(s)
}

// FinishWithOptions belongs to the Span interface
func (s *MockSpan) FinishWithOptions(opts opentracing.FinishOptions) {
	s.spanContext.Lock()
	s.FinishTime = opts.FinishTime
	s.logs = append(s.logs, opts.BulkLogData...)
	s.spanContext.Unlock()
	s.tracer.recordSpan(s)
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
	s.spanContext.Lock()
	defer s.spanContext.Unlock()
	s.logs = append(s.logs, data)
}

// SetOperationName belongs to the Span interface
func (s *MockSpan) SetOperationName(operationName string) opentracing.Span {
	s.spanContext.Lock()
	defer s.spanContext.Unlock()
	s.OperationName = operationName
	return s
}

// Tracer belongs to the Span interface
func (s *MockSpan) Tracer() opentracing.Tracer {
	return s.tracer
}
