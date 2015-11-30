package opentracing

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// Creates and returns a standard OpenTracer which defers to `rec` and
// `ctxIdSource` as appropriate.
//
// See `SetGlobalTracer()`.
func NewStandardTracer(rec ProcessRecorder, ctxIdSource TraceContextIDSource) OpenTracer {
	return &standardOpenTracer{
		TraceContextIDSource: ctxIdSource,
		recorder:             rec,
	}
}

// A straightforward implementation of the `Span` interface. Created via
// standardOpenTracer (see `NewStandardTracer()`).
type standardSpan struct {
	lock     sync.Mutex
	tracer   *standardOpenTracer
	recorder ProcessRecorder
	raw      RawSpan
}

func (s *standardSpan) StartChild(operationName string, initialTags ...Tags) Span {
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

func (s *standardSpan) TraceContext() *TraceContext {
	// No need for a lock since s.raw.TraceContext is not modified after
	// initialization.
	return s.raw.TraceContext
}

// Implements the `OpenTracer` interface.
type standardOpenTracer struct {
	TraceContextIDSource

	recorder ProcessRecorder
}

func (s *standardOpenTracer) StartNewTrace(
	operationName string,
	initialTags ...Tags,
) Span {
	return s.startSpanGeneric(
		operationName,
		NewRootTraceContext(s),
		nil,
	)
}

func (s *standardOpenTracer) JoinTrace(
	operationName string,
	parent interface{},
	initialTags ...Tags,
) Span {
	if goCtx, ok := parent.(context.Context); ok {
		return s.startSpanWithGoContextParent(operationName, goCtx)
	} else if traceCtx, ok := parent.(*TraceContext); ok {
		return s.startSpanWithTraceContextParent(operationName, traceCtx)
	} else {
		panic(fmt.Errorf("Invalid parent type: %v", reflect.TypeOf(parent)))
	}
}

func (s *standardOpenTracer) startSpanWithGoContextParent(
	operationName string,
	parent context.Context,
) Span {
	if oldSpan := SpanFromGoContext(parent); oldSpan != nil {
		childCtx, tags := oldSpan.TraceContext().NewChild()
		return s.startSpanGeneric(
			operationName,
			childCtx,
			tags,
		)
	} else {
		return s.startSpanGeneric(
			operationName,
			NewRootTraceContext(s),
			nil,
		)
	}
}

func (s *standardOpenTracer) startSpanWithTraceContextParent(
	operationName string,
	parent *TraceContext,
) Span {
	childCtx, tags := parent.NewChild()
	return s.startSpanGeneric(
		operationName,
		childCtx,
		tags,
	)
}

// A helper for standardSpan creation.
func (s *standardOpenTracer) startSpanGeneric(
	operationName string,
	childCtx *TraceContext,
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
