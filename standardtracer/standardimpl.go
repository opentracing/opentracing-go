package standardtracer

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"

	"golang.org/x/net/context"
)

// New creates and returns a standard Tracer which defers to `recorder` and
// `source` as appropriate.
func New(recorder Recorder, source opentracing.TraceContextSource) opentracing.Tracer {
	return &standardTracer{
		TraceContextSource: source,
		recorder:           recorder,
	}
}

// Implements the `Span` interface. Created via standardTracer (see
// `NewStandardTracer()`).
type standardSpan struct {
	lock     sync.Mutex
	tracer   *standardTracer
	recorder Recorder
	raw      RawSpan
}

func (s *standardSpan) StartChild(operationName string) opentracing.Span {
	childCtx, childTags := s.tracer.NewChildTraceContext(s.raw.TraceContext)
	return s.tracer.startSpanGeneric(operationName, childCtx, childTags)
}

func (s *standardSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.raw.Tags[key] = value
	return s
}

func (s *standardSpan) LogEvent(event string) {
	s.Log(opentracing.LogData{
		Event: event,
	})
}

func (s *standardSpan) LogEventWithPayload(event string, payload interface{}) {
	s.Log(opentracing.LogData{
		Event:   event,
		Payload: payload,
	})
}

func (s *standardSpan) Log(ld opentracing.LogData) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if ld.Timestamp.IsZero() {
		ld.Timestamp = time.Now()
	}
	s.raw.Logs = append(s.raw.Logs, &ld)
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

// Implements the `Tracer` interface.
type standardTracer struct {
	opentracing.TraceContextSource

	recorder Recorder
}

func (s *standardTracer) StartTrace(
	operationName string,
) opentracing.Span {
	return s.startSpanGeneric(
		operationName,
		s.NewRootTraceContext(),
		nil,
	)
}

func (s *standardTracer) StartSpanWithContext(
	operationName string,
	ctx opentracing.TraceContext,
) opentracing.Span {
	return s.startSpanGeneric(
		operationName,
		ctx,
		nil,
	)
}

func (s *standardTracer) JoinTrace(
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

func (s *standardTracer) startSpanWithGoContextParent(
	operationName string,
	parent context.Context,
) opentracing.Span {
	if oldSpan := opentracing.SpanFromGoContext(parent); oldSpan != nil {
		childCtx, tags := s.NewChildTraceContext(oldSpan.TraceContext())
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

func (s *standardTracer) startSpanWithTraceContextParent(
	operationName string,
	parent opentracing.TraceContext,
) opentracing.Span {
	childCtx, tags := s.NewChildTraceContext(parent)
	return s.startSpanGeneric(
		operationName,
		childCtx,
		tags,
	)
}

// A helper for standardSpan creation.
func (s *standardTracer) startSpanGeneric(
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
		raw: RawSpan{
			TraceContext: childCtx,
			Operation:    operationName,
			Start:        time.Now(),
			Duration:     -1,
			Tags:         tags,
			Logs:         []*opentracing.LogData{},
		},
	}
	return span
}
