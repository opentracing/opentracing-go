package standardtracer

import (
	"fmt"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
)

// Implements the `Span` interface. Created via tracerImpl (see
// `standardtracer.New()`).
type spanImpl struct {
	lock     sync.Mutex
	tracer   *tracerImpl
	recorder SpanRecorder
	raw      RawSpan
}

func (s *spanImpl) StartChild(operationName string) opentracing.Span {
	childCtx := s.raw.StandardContext.NewChild()
	return s.tracer.startSpanGeneric(operationName, childCtx)
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

func (s *spanImpl) Info(message string, payload ...interface{}) {
	s.internalLog(false, message, payload...)
}

func (s *spanImpl) Error(message string, payload ...interface{}) {
	s.internalLog(true, message, payload...)
}

func (s *spanImpl) internalLog(isErr bool, message string, payload ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.raw.Logs = append(s.raw.Logs, &RawLog{
		Timestamp: time.Now(),
		Error:     isErr,
		Message:   message,
		Payload:   payload,
	})
}

func (s *spanImpl) Finish() {
	duration := time.Since(s.raw.Start)
	s.lock.Lock()
	defer s.lock.Unlock()

	s.raw.Duration = duration
	s.recorder.RecordSpan(&s.raw)
}

func (s *spanImpl) SetTraceAttribute(restrictedKey, val string) opentracing.Span {
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.raw.StandardContext.tagLock.Lock()
	defer s.raw.StandardContext.tagLock.Unlock()

	s.raw.StandardContext.traceAttrs[canonicalKey] = val
	return s
}

func (s *spanImpl) TraceAttribute(restrictedKey string) string {
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.raw.StandardContext.tagLock.RLock()
	defer s.raw.StandardContext.tagLock.RUnlock()

	return s.raw.StandardContext.traceAttrs[canonicalKey]
}
