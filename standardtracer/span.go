package standardtracer

import (
	"fmt"
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

// Implements the `Span` interface. Created via tracerImpl (see
// `standardtracer.New()`).
type spanImpl struct {
	lock     sync.Mutex
	tracer   *tracerImpl
	recorder SpanRecorder
	raw      RawSpan
}

func (s *spanImpl) SetOperationName(operationName string) opentracing.Span {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.raw.Operation = operationName
	return s
}

func (s *spanImpl) SetTag(key string, value interface{}) opentracing.Span {
	s.lock.Lock()
	defer s.lock.Unlock()
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
	s.lock.Lock()
	defer s.lock.Unlock()

	if ld.Timestamp.IsZero() {
		ld.Timestamp = time.Now()
	}

	s.raw.Logs = append(s.raw.Logs, &ld)
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

	s.lock.Lock()
	defer s.lock.Unlock()
	if opts.BulkLogData != nil {
		s.raw.Logs = append(s.raw.Logs, opts.BulkLogData...)
	}
	s.raw.Duration = duration
	s.recorder.RecordSpan(&s.raw)
}

func (s *spanImpl) SetTraceAttribute(restrictedKey, val string) opentracing.Span {
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.raw.StandardContext.attrMu.Lock()
	defer s.raw.StandardContext.attrMu.Unlock()

	s.raw.StandardContext.traceAttrs[canonicalKey] = val
	return s
}

func (s *spanImpl) TraceAttribute(restrictedKey string) string {
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.raw.StandardContext.attrMu.RLock()
	defer s.raw.StandardContext.attrMu.RUnlock()

	return s.raw.StandardContext.traceAttrs[canonicalKey]
}

func (s *spanImpl) Tracer() opentracing.Tracer {
	return s.tracer
}
