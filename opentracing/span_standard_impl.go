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

func (s *standardSpan) StartChildSpan(operationName string, parent ...context.Context) (Span, context.Context) {
	var parentGoCtx context.Context
	if len(parent) > 0 {
		parentGoCtx = parent[0]
	} else {
		parentGoCtx = context.Background()
	}
	childCtx, childTags := s.raw.TraceContext.NewChild()
	return s.tracer.startSpanGeneric(operationName, parentGoCtx, childCtx, childTags)
}

func (s *standardSpan) SetTag(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.raw.Tags[key] = fmt.Sprint(value)
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
	return &s.raw.TraceContext
}

// Implements the `OpenTracer` interface.
type standardOpenTracer struct {
	TraceContextIDSource

	recorder ProcessRecorder
}

func (s *standardOpenTracer) StartSpan(
	operationName string,
	parent ...interface{},
) (Span, context.Context) {
	if len(parent) == 0 {
		return s.startSpanGeneric(
			operationName,
			context.Background(),
			s.NewRootTraceContextID(),
			nil,
		)
	} else {
		if goCtx, ok := parent[0].(context.Context); ok {
			return s.startSpanWithGoContextParent(operationName, goCtx)
		} else if ctxID, ok := parent[0].(*TraceContext); ok {
			return s.startSpanWithTraceContextParent(operationName, ctxID)
		} else {
			panic(fmt.Errorf("Invalid parent type: %v", reflect.TypeOf(parent[0])))
		}
	}
}

func (s *standardOpenTracer) startSpanWithTraceContextIDParent(
	operationName string,
	parent *TraceContext,
) (Span, context.Context) {
	childCtx, tags := parent.NewChild()
	return s.startSpanGeneric(
		operationName,
		context.Background(),
		childCtx,
		tags,
	)
}

func (s *standardOpenTracer) startSpanWithGoContextParent(
	operationName string,
	parent context.Context,
) (Span, context.Context) {
	if oldSpan := SpanFromGoContext(parent); oldSpan != nil {
		childCtx, tags := oldSpan.TraceContext().NewChild()
		return s.startSpanGeneric(
			operationName,
			parent,
			childCtx,
			tags,
		)
	} else {
		return s.startSpanGeneric(
			operationName,
			parent,
			NewRootTraceContext(s),
			nil,
		)
	}
}

// A helper for standardSpan creation.
func (s *standardOpenTracer) startSpanGeneric(
	operationName string,
	parentGoCtx context.Context,
	childCtx *TraceContext,
	tags Tags,
) (Span, context.Context) {
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
	goCtx := GoContextWithSpan(parentGoCtx, span)
	return span, goCtx
}
