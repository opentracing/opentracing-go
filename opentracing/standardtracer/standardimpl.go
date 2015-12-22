package standardtracer

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/opentracing/api-golang/opentracing"

	"golang.org/x/net/context"
)

// New creates and returns a standard OpenTracer which defers to `rec` and
// `source` as appropriate.
func New(rec opentracing.Recorder, source opentracing.TraceContextSource) opentracing.OpenTracer {
	return &standardOpenTracer{
		TraceContextSource: source,
		recorder:           rec,
	}
}

// Implements the `Span` interface. Created via standardOpenTracer (see
// `NewStandardTracer()`).
type standardSpan struct {
	lock     sync.Mutex
	tracer   *standardOpenTracer
	recorder opentracing.Recorder
	raw      opentracing.RawSpan
}

func (s *standardSpan) StartChild(operationName string) opentracing.Span {
	childCtx, childTags := s.raw.TraceContext.NewChild()
	return s.tracer.startSpanGeneric(operationName, childCtx, childTags)
}

func (s *standardSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.raw.Tags[key] = fmt.Sprint(value)
	return s
}

func (s *standardSpan) SetTags(tags opentracing.Tags) opentracing.Span {
	s.lock.Lock()
	defer s.lock.Unlock()

	for k, v := range tags {
		s.raw.Tags[k] = fmt.Sprint(v)
	}
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

	s.raw.Logs = append(s.raw.Logs, &opentracing.RawLog{
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

func (s *standardSpan) TraceContext() opentracing.TraceContext {
	// No need for a lock since s.raw.TraceContext is not modified after
	// initialization.
	return s.raw.TraceContext
}

func (s *standardSpan) AddToGoContext(ctx context.Context) (opentracing.Span, context.Context) {
	return s, opentracing.GoContextWithSpan(ctx, s)
}

// Implements the `OpenTracer` interface.
type standardOpenTracer struct {
	opentracing.TraceContextSource

	recorder opentracing.Recorder
}

func (s *standardOpenTracer) StartTrace(
	operationName string,
) opentracing.Span {
	return s.startSpanGeneric(
		operationName,
		s.NewRootTraceContext(),
		nil,
	)
}

func (s *standardOpenTracer) JoinTrace(
	operationName string,
	parent interface{},
) opentracing.Span {
	if goCtx, ok := parent.(context.Context); ok {
		return s.startSpanWithGoContextParent(operationName, goCtx)
	} else if traceCtx, ok := parent.(opentracing.TraceContext); ok {
		return s.startSpanWithTraceContextParent(operationName, traceCtx)
	} else {
		panic(fmt.Errorf("Invalid parent type: %v", reflect.TypeOf(parent)))
	}
}

func (s *standardOpenTracer) startSpanWithGoContextParent(
	operationName string,
	parent context.Context,
) opentracing.Span {
	if oldSpan := opentracing.SpanFromGoContext(parent); oldSpan != nil {
		childCtx, tags := oldSpan.TraceContext().NewChild()
		return s.startSpanGeneric(
			operationName,
			childCtx,
			tags,
		)
	}

	return s.startSpanGeneric(
		operationName,
		s.NewRootTraceContext(),
		nil,
	)
}

func (s *standardOpenTracer) startSpanWithTraceContextParent(
	operationName string,
	parent opentracing.TraceContext,
) opentracing.Span {
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
	childCtx opentracing.TraceContext,
	tags opentracing.Tags,
) opentracing.Span {
	if tags == nil {
		tags = opentracing.Tags{}
	}
	span := &standardSpan{
		tracer:   s,
		recorder: s.recorder,
		raw: opentracing.RawSpan{
			TraceContext: childCtx,
			Operation:    operationName,
			Start:        time.Now(),
			Duration:     -1,
			Tags:         tags,
			Logs:         []*opentracing.RawLog{},
		},
	}
	return span
}
