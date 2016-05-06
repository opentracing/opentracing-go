package mocktracer

import (
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
)

// New returns a MockTracer opentracing.Tracer implementation that's intended
// to facilitate tests of OpenTracing instrumentation.
func New() *MockTracer {
	return &MockTracer{
		FinishedSpans: []*MockSpan{},
	}
}

// MockTracer is a for-testing-only opentracing.Tracer implementation. It is
// entirely unsuitable for production use but appropriate for tests that want
// to verify tracing behavior.
type MockTracer struct {
	FinishedSpans []*MockSpan
}

// MockSpan is an opentracing.Span implementation that exports its internal
// state for testing purposes.
type MockSpan struct {
	SpanID   int
	ParentID int

	OperationName string
	StartTime     time.Time
	FinishTime    time.Time
	Tags          map[string]interface{}
	Baggage       map[string]string
	Logs          []opentracing.LogData

	tracer *MockTracer
}

// Reset clears the exported MockTracer.FinishedSpans field. Note that any
// extant MockSpans will still append to FinishedSpans when they Finish(), even
// after a call to Reset().
func (t *MockTracer) Reset() {
	t.FinishedSpans = []*MockSpan{}
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
	switch format {
	case opentracing.TextMap:
		writer := carrier.(opentracing.TextMapWriter)
		// Ids:
		writer.Set(mockTextMapIdsPrefix+"spanid", strconv.Itoa(span.SpanID))
		// Baggage:
		for baggageKey, baggageVal := range span.Baggage {
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
		err := carrier.(opentracing.TextMapReader).ForeachKey(func(key, val string) error {
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
				rval.Baggage[lowerKey[len(mockTextMapBaggagePrefix):]] = val
			}
			return nil
		})
		return rval, err
	}
	return nil, opentracing.ErrTraceNotFound
}

var mockIDSource = 1

func nextMockID() int {
	mockIDSource++
	return mockIDSource
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
		Tags:          tags,
		Baggage:       map[string]string{},
		Logs:          []opentracing.LogData{},

		tracer: t,
	}
}

// SetTag belongs to the Span interface
func (s *MockSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.Tags[key] = value
	return s
}

// Finish belongs to the Span interface
func (s *MockSpan) Finish() {
	s.FinishTime = time.Now()
	s.tracer.FinishedSpans = append(s.tracer.FinishedSpans, s)
}

// FinishWithOptions belongs to the Span interface
func (s *MockSpan) FinishWithOptions(opts opentracing.FinishOptions) {
	s.FinishTime = opts.FinishTime
	s.Logs = append(s.Logs, opts.BulkLogData...)
	s.tracer.FinishedSpans = append(s.tracer.FinishedSpans, s)
}

// SetBaggageItem belongs to the Span interface
func (s *MockSpan) SetBaggageItem(key, val string) opentracing.Span {
	s.Baggage[key] = val
	return s
}

// BaggageItem belongs to the Span interface
func (s *MockSpan) BaggageItem(key string) string {
	return s.Baggage[key]
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
	s.Logs = append(s.Logs, data)
}

// SetOperationName belongs to the Span interface
func (s *MockSpan) SetOperationName(operationName string) opentracing.Span {
	s.OperationName = operationName
	return s
}

// Tracer belongs to the Span interface
func (s *MockSpan) Tracer() opentracing.Tracer {
	return s.tracer
}
