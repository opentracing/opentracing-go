package mocktracer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

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

	// All of the below (including SpanContext) are protected by SpanContext's
	// embedded RWMutex.
	SpanContext *MockSpanContext
	tags        map[string]interface{}
	logs        []opentracing.LogData
	tracer      *MockTracer
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
		sampled = opts.References[0].Referee.(*MockSpanContext).Sampled
		baggage = opts.References[0].Referee.(*MockSpanContext).Baggage()
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
