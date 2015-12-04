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

func (s *standardSpan) StartChild(operationName string, keyValueTags ...interface{}) opentracing.Span {
	childCtx, childTags := s.raw.TraceContext.NewChild()
	return s.tracer.startSpanGeneric(operationName, childCtx, childTags)
}

func (s *standardSpan) SetTag(key string, value interface{}) opentracing.Span {
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
	keyValueTags ...interface{},
) opentracing.Span {
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
) opentracing.Span {
	if goCtx, ok := parent.(context.Context); ok {
		return s.startSpanWithGoContextParent(operationName, goCtx, keyValueTags...)
	} else if traceCtx, ok := parent.(opentracing.TraceContext); ok {
		return s.startSpanWithTraceContextParent(operationName, traceCtx, keyValueTags...)
	} else {
		panic(fmt.Errorf("Invalid parent type: %v", reflect.TypeOf(parent)))
	}
}

func (s *standardOpenTracer) startSpanWithGoContextParent(
	operationName string,
	parent context.Context,
	keyValueTags ...interface{},
) opentracing.Span {
	if oldSpan := opentracing.SpanFromGoContext(parent); oldSpan != nil {
		childCtx, tags := oldSpan.TraceContext().NewChild()
		tags.Merge(keyValueListToTags(keyValueTags))
		return s.startSpanGeneric(
			operationName,
			childCtx,
			tags,
		)
	}

	tags := keyValueListToTags(keyValueTags)
	return s.startSpanGeneric(
		operationName,
		s.NewRootTraceContext(),
		tags,
	)
}

func (s *standardOpenTracer) startSpanWithTraceContextParent(
	operationName string,
	parent opentracing.TraceContext,
	keyValueTags ...interface{},
) opentracing.Span {
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

func keyValueListToTags(keyValueTags []interface{}) opentracing.Tags {
	if len(keyValueTags)%2 != 0 {
		panic(fmt.Errorf(
			"there must be an even number of keyValueTags params to split them into pairs: got %v",
			len(keyValueTags)))
	}
	rval := make(opentracing.Tags, len(keyValueTags)/2)
	var k string
	for i, keyOrVal := range keyValueTags {
		if i%2 == 0 {
			var ok bool
			k, ok = keyOrVal.(string)
			if !ok {
				panic(fmt.Errorf(
					"even-indexed keyValueTags (i.e., the keys) must be strings: got %v",
					reflect.TypeOf(keyOrVal)))
			}
		} else {
			rval[k] = keyOrVal
		}
	}
	return rval
}
