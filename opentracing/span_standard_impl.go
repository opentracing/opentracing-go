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
func NewStandardTracer(rec Recorder, ctxIdSource ContextIDSource) OpenTracer {
	return &standardOpenTracer{
		ContextIDSource: ctxIdSource,
		recorder:        rec,
	}
}

// A straightforward implementation of the `Span` interface. Created via
// standardOpenTracer (see `NewStandardTracer()`).
type standardSpan struct {
	lock     sync.Mutex
	tracer   *standardOpenTracer
	recorder Recorder
	raw      RawSpan
}

func (s *standardSpan) StartChildSpan(operationName string, parent ...context.Context) (Span, context.Context) {
	var parentGoCtx context.Context
	if len(parent) > 0 {
		parentGoCtx = parent[0]
	} else {
		parentGoCtx = context.Background()
	}
	childID, childTags := s.raw.ContextID.NewChild()
	return s.tracer.startSpanGeneric(operationName, parentGoCtx, childID, childTags)
}

func (s *standardSpan) SetTag(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.raw.Tags[key] = fmt.Sprint(value)
}

func (s *standardSpan) Info(message string, payload ...interface{}) {
	s.internalLog(SeverityInfo, message, payload...)
}

func (s *standardSpan) Warning(message string, payload ...interface{}) {
	s.internalLog(SeverityWarning, message, payload...)
}

func (s *standardSpan) Error(message string, payload ...interface{}) {
	s.internalLog(SeverityError, message, payload...)
}

func (s *standardSpan) internalLog(sev Severity, message string, payload ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.raw.Logs = append(s.raw.Logs, &RawLog{
		Timestamp: time.Now(),
		Severity:  sev,
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

func (s *standardSpan) ContextID() ContextID {
	// No need for a lock since s.raw.ContextID is not modified after
	// initialization.
	return s.raw.ContextID
}

// Implements the `OpenTracer` interface.
type standardOpenTracer struct {
	ContextIDSource

	recorder Recorder
}

func (s *standardOpenTracer) StartSpan(
	operationName string,
	parent ...interface{},
) (Span, context.Context) {
	if len(parent) == 0 {
		return s.startSpanGeneric(
			operationName,
			context.Background(),
			s.NewRootContextID(),
			nil,
		)
	} else {
		if goCtx, ok := parent[0].(context.Context); ok {
			return s.startSpanWithGoContextParent(operationName, goCtx)
		} else if ctxID, ok := parent[0].(ContextID); ok {
			return s.startSpanWithContextIDParent(operationName, ctxID)
		} else {
			panic(fmt.Errorf("Invalid parent type: %v", reflect.TypeOf(parent[0])))
		}
	}
}

func (s *standardOpenTracer) startSpanWithContextIDParent(
	operationName string,
	parent ContextID,
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
		childCtxId, tags := oldSpan.ContextID().NewChild()
		return s.startSpanGeneric(
			operationName,
			parent,
			childCtxId,
			tags,
		)
	} else {
		return s.startSpanGeneric(
			operationName,
			parent,
			s.NewRootContextID(),
			nil,
		)
	}
}

// A helper for standardSpan creation.
func (s *standardOpenTracer) startSpanGeneric(
	operationName string,
	parentGoCtx context.Context,
	childCtxID ContextID,
	tags Tags,
) (Span, context.Context) {
	if tags == nil {
		tags = Tags{}
	}
	span := &standardSpan{
		tracer:   s,
		recorder: s.recorder,
		raw: RawSpan{
			ContextID: childCtxID,
			Operation: operationName,
			Start:     time.Now(),
			Duration:  -1,
			Tags:      tags,
			Logs:      []*RawLog{},
		},
	}
	goCtx := GoContextWithSpan(parentGoCtx, span)
	return span, goCtx
}
