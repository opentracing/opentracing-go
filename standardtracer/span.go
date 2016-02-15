package standardtracer

import (
	"fmt"
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Implements the `Span` interface. Created via tracerImpl (see
// `standardtracer.New()`).
type spanImpl struct {
	tracer     *tracerImpl
	sync.Mutex // protects the fields below
	raw        RawSpan
}

func (s *spanImpl) reset() {
	s.tracer = nil
	s.raw = RawSpan{}
	s.raw.Attributes = nil // TODO(tschottdorf): is clearing out the map better?
}

func (s *spanImpl) SetOperationName(operationName string) opentracing.Span {
	s.Lock()
	defer s.Unlock()
	s.raw.Operation = operationName
	return s
}

func (s *spanImpl) SetTag(key string, value interface{}) opentracing.Span {
	s.Lock()
	defer s.Unlock()
	if key == string(ext.SamplingPriority) {
		s.raw.Sampled = true
		return s
	}

	if s.raw.Tags == nil {
		s.raw.Tags = opentracing.Tags{}
	}
	s.raw.Tags[key] = value
	return s
}

func (s *spanImpl) LogEvent(event string) {
	s.Log(opentracing.LogData{
		Event: event,
	})
}

func (s *spanImpl) LogEventWithPayload(event string, payload interface{}) {
	s.Log(opentracing.LogData{
		Event:   event,
		Payload: payload,
	})
}

func (s *spanImpl) Log(ld opentracing.LogData) {
	s.Lock()
	defer s.Unlock()

	if ld.Timestamp.IsZero() {
		ld.Timestamp = time.Now()
	}

	s.raw.Logs = append(s.raw.Logs, ld)
}

func (s *spanImpl) Finish() {
	s.FinishWithOptions(opentracing.FinishOptions{})
}

func (s *spanImpl) FinishWithOptions(opts opentracing.FinishOptions) {
	finishTime := opts.FinishTime
	if finishTime.IsZero() {
		finishTime = time.Now()
	}
	duration := finishTime.Sub(s.raw.Start)

	s.Lock()
	defer s.Unlock()
	if opts.BulkLogData != nil {
		s.raw.Logs = append(s.raw.Logs, opts.BulkLogData...)
	}
	s.raw.Duration = duration
	s.tracer.recorder.RecordSpan(s.raw)
	s.tracer.spanPool.Put(s)
}

func (s *spanImpl) SetTraceAttribute(restrictedKey, val string) opentracing.Span {
	// TODO(tschottdorf): this is in a bad place performance-wise. Most callers
	// will have fixed attribute keys, so the are in a much better position to
	// canonicalize here. Alternatively, we could put a LRU cache here, but that
	// is more complexity than is warranted for.
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.Lock()
	defer s.Unlock()
	if s.raw.Attributes == nil {
		s.raw.Attributes = make(map[string]string)
	}
	s.raw.Attributes[canonicalKey] = val
	return s
}

func (s *spanImpl) TraceAttribute(restrictedKey string) string {
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.Lock()
	defer s.Unlock()

	return s.raw.Attributes[canonicalKey]
}

func (s *spanImpl) Tracer() opentracing.Tracer {
	return s.tracer
}
