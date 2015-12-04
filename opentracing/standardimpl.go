package opentracing

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// Creates and returns a standard OpenTracer which defers to `rec` and
// `ctxIDSource` as appropriate.
func NewStandardTracer(rec Recorder, ctxIDSource TraceContextSource) OpenTracer {
	return &standardOpenTracer{
		TraceContextSource: ctxIDSource,
		recorder:           rec,
	}
}

// Implements the `Span` interface. Created via standardOpenTracer (see
// `NewStandardTracer()`).
type standardSpan struct {
	lock     sync.Mutex
	tracer   *standardOpenTracer
	recorder Recorder
	raw      RawSpan
}

func (s *standardSpan) StartChild(operationName string, keyValueTags ...interface{}) Span {
	childCtx, childTags := s.raw.TraceContext.NewChild()
	return s.tracer.startSpanGeneric(operationName, childCtx, childTags)
}

func (s *standardSpan) SetTag(key string, value interface{}) Span {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.raw.Tags[key] = fmt.Sprint(value)
	return s
}

func (s *standardSpan) Info(message string, payload ...interface{}) {
	s.internalLog(false, message, payload...)
}

func (s *standardSpan) Error(message string, payload ...interface{}) {
	s.internalLog(true, message, payload...)
}

func (s *standardSpan) internalLog(isErr bool, message string, payload ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.raw.Logs = append(s.raw.Logs, &RawLog{
		Timestamp: time.Now(),
		Error:     isErr,
		Message:   message,
		Payload:   payload,
	})
}

func (s *standardSpan) Finish() {
	duration := time.Since(s.raw.Start)
	s.lock.Lock()
	defer s.lock.Unlock()

	s.raw.Duration = duration
	s.recorder.RecordSpan(&s.raw)
}

func (s *standardSpan) TraceContext() TraceContext {
	// No need for a lock since s.raw.TraceContext is not modified after
	// initialization.
	return s.raw.TraceContext
}

func (s *standardSpan) AddToGoContext(ctx context.Context) (Span, context.Context) {
	return s, GoContextWithSpan(ctx, s)
}

// Implements the `OpenTracer` interface.
type standardOpenTracer struct {
	TraceContextSource

	recorder Recorder
}

func (s *standardOpenTracer) StartTrace(
	operationName string,
	keyValueTags ...interface{},
) Span {
	tags := keyValueListToTags(keyValueTags)
	return s.startSpanGeneric(
		operationName,
		s.NewRootTraceContext(),
		tags,
	)
}

func (s *standardOpenTracer) JoinTrace(
	operationName string,
	parent interface{},
	keyValueTags ...interface{},
) Span {
	if goCtx, ok := parent.(context.Context); ok {
		return s.startSpanWithGoContextParent(operationName, goCtx, keyValueTags...)
	} else if traceCtx, ok := parent.(TraceContext); ok {
		return s.startSpanWithTraceContextParent(operationName, traceCtx, keyValueTags...)
	} else {
		panic(fmt.Errorf("Invalid parent type: %v", reflect.TypeOf(parent)))
	}
}

func (s *standardOpenTracer) startSpanWithGoContextParent(
	operationName string,
	parent context.Context,
	keyValueTags ...interface{},
) Span {
	if oldSpan := SpanFromGoContext(parent); oldSpan != nil {
		childCtx, tags := oldSpan.TraceContext().NewChild()
		tags.Merge(keyValueListToTags(keyValueTags))
		return s.startSpanGeneric(
			operationName,
			childCtx,
			tags,
		)
	} else {
		tags := keyValueListToTags(keyValueTags)
		return s.startSpanGeneric(
			operationName,
			s.NewRootTraceContext(),
			tags,
		)
	}
}

func (s *standardOpenTracer) startSpanWithTraceContextParent(
	operationName string,
	parent TraceContext,
	keyValueTags ...interface{},
) Span {
	childCtx, tags := parent.NewChild()
	tags.Merge(keyValueListToTags(keyValueTags))
	return s.startSpanGeneric(
		operationName,
		childCtx,
		tags,
	)
}

// A helper for standardSpan creation.
func (s *standardOpenTracer) startSpanGeneric(
	operationName string,
	childCtx TraceContext,
	tags Tags,
) Span {
	if tags == nil {
		tags = Tags{}
	}
	span := &standardSpan{
		tracer:   s,
		recorder: s.recorder,
		raw: RawSpan{
			TraceContext: childCtx,
			Operation:    operationName,
			Start:        time.Now(),
			Duration:     -1,
			Tags:         tags,
			Logs:         []*RawLog{},
		},
	}
	return span
}
