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
// It is entirely unsuitable for production use, but appropriate for tests
// that want to verify tracing behavior in other frameworks/applications.
type MockTracer struct {
	sync.RWMutex
	finishedSpans []*MockSpan
}

// MockSpan is an opentracing.Span implementation that exports its internal
// state for testing purposes.
type MockSpan struct {
	sync.RWMutex

	SpanID   int
	ParentID int

	OperationName string
	StartTime     time.Time
	FinishTime    time.Time

	tags    map[string]interface{}
	baggage map[string]string
	logs    []opentracing.LogData

	tracer *MockTracer
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
	s.RLock()
	defer s.RUnlock()
	tags := make(map[string]interface{})
	for k, v := range s.tags {
		tags[k] = v
	}
	return tags
}

// GetTag returns a single tag
func (s *MockSpan) GetTag(k string) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.tags[k]
}

// GetBaggage returns a copy of baggage items in the span
func (s *MockSpan) GetBaggage() map[string]string {
	s.RLock()
	defer s.RUnlock()
	baggage := make(map[string]string)
	for k, v := range s.baggage {
		baggage[k] = v
	}
	return baggage
}

// GetLogs returns a copy of logs accumulated in the span so far
func (s *MockSpan) GetLogs() []opentracing.LogData {
	s.RLock()
	defer s.RUnlock()
	logs := make([]opentracing.LogData, len(s.logs))
	copy(logs, s.logs)
	return logs
}

// StartSpan belongs to the Tracer interface.
func (t *MockTracer) StartSpan(operationName string) opentracing.Span {
	return newMockSpan(t, opentracing.StartSpanOptions{
		OperationName: operationName,
	})
}

// StartSpanWithOptions belongs to the Tracer interface.
func (t *MockTracer) StartSpanWithOptions(opts opentracing.StartSpanOptions) opentracing.Span {
	return newMockSpan(t, opts)
}

const mockTextMapIdsPrefix = "mockpfx-ids-"
const mockTextMapBaggagePrefix = "mockpfx-baggage-"

// Inject belongs to the Tracer interface.
func (t *MockTracer) Inject(sp opentracing.Span, format interface{}, carrier interface{}) error {
	span := sp.(*MockSpan)
	span.RLock()
	defer span.RUnlock()

	switch format {
	case opentracing.TextMap:
		writer, ok := carrier.(opentracing.TextMapWriter)
		if !ok {
			return opentracing.ErrInvalidCarrier
		}
		// Ids:
		writer.Set(mockTextMapIdsPrefix+"spanid", strconv.Itoa(span.SpanID))
		// Baggage:
		for baggageKey, baggageVal := range span.baggage {
			writer.Set(mockTextMapBaggagePrefix+baggageKey, baggageVal)
		}
		return nil
	}
	return opentracing.ErrUnsupportedFormat
}

// Join belongs to the Tracer interface.
func (t *MockTracer) Join(operationName string, format interface{}, carrier interface{}) (opentracing.Span, error) {
	switch format {
	case opentracing.TextMap:
		rval := newMockSpan(t, opentracing.StartSpanOptions{
			OperationName: operationName,
		})
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
				rval.ParentID = i
			case strings.HasPrefix(lowerKey, mockTextMapBaggagePrefix):
				// Baggage:
				rval.baggage[lowerKey[len(mockTextMapBaggagePrefix):]] = val
			}
			return nil
		})
		return rval, err
	}
	return nil, opentracing.ErrUnsupportedFormat
}

var mockIDSource = uint32(42)

func nextMockID() int {
	atomic.AddUint32(&mockIDSource, 1)
	return int(atomic.LoadUint32(&mockIDSource))
}

func newMockSpan(t *MockTracer, opts opentracing.StartSpanOptions) *MockSpan {
	tags := opts.Tags
	if tags == nil {
		tags = map[string]interface{}{}
	}
	parentID := int(0)
	if opts.Parent != nil {
		parentID = opts.Parent.(*MockSpan).SpanID
	}
	startTime := opts.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}
	return &MockSpan{
		SpanID:   nextMockID(),
		ParentID: parentID,

		OperationName: opts.OperationName,
		StartTime:     startTime,
		tags:          tags,
		baggage:       map[string]string{},
		logs:          []opentracing.LogData{},

		tracer: t,
	}
}

// SetTag belongs to the Span interface
func (s *MockSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.Lock()
	defer s.Unlock()
	s.tags[key] = value
	return s
}

// Finish belongs to the Span interface
func (s *MockSpan) Finish() {
	s.Lock()
	s.FinishTime = time.Now()
	s.Unlock()
	s.tracer.recordSpan(s)
}

// FinishWithOptions belongs to the Span interface
func (s *MockSpan) FinishWithOptions(opts opentracing.FinishOptions) {
	s.Lock()
	s.FinishTime = opts.FinishTime
	s.logs = append(s.logs, opts.BulkLogData...)
	s.Unlock()
	s.tracer.recordSpan(s)
}

func (t *MockTracer) recordSpan(span *MockSpan) {
	t.Lock()
	defer t.Unlock()
	t.finishedSpans = append(t.finishedSpans, span)
}

// SetBaggageItem belongs to the Span interface
func (s *MockSpan) SetBaggageItem(key, val string) opentracing.Span {
	s.Lock()
	defer s.Unlock()
	s.baggage[key] = val
	return s
}

// BaggageItem belongs to the Span interface
func (s *MockSpan) BaggageItem(key string) string {
	s.RLock()
	defer s.RUnlock()
	return s.baggage[key]
}

// ForeachBaggageItem belongs to the Span interface
func (s *MockSpan) ForeachBaggageItem(handler func(k, v string) bool) {
	s.RLock()
	defer s.RUnlock()
	for k, v := range s.baggage {
		if !handler(k, v) {
			break
		}
	}
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
	s.Lock()
	defer s.Unlock()
	s.logs = append(s.logs, data)
}

// SetOperationName belongs to the Span interface
func (s *MockSpan) SetOperationName(operationName string) opentracing.Span {
	s.Lock()
	defer s.Unlock()
	s.OperationName = operationName
	return s
}

// Tracer belongs to the Span interface
func (s *MockSpan) Tracer() opentracing.Tracer {
	return s.tracer
}
