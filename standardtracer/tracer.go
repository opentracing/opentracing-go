package standardtracer

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

// New creates and returns a standard Tracer which defers to `recorder` after
// RawSpans have been assembled.
func New(recorder SpanRecorder) opentracing.Tracer {
	return &tracerImpl{
		recorder: recorder,
	}
}

// Implements the `Tracer` interface.
type tracerImpl struct {
	recorder SpanRecorder
}

func (s *tracerImpl) StartTrace(
	operationName string,
) opentracing.Span {
	return s.startSpanGeneric(
		operationName,
		NewRootStandardContext(),
	)
}

func (s *tracerImpl) startSpanWithSpanParent(
	operationName string,
	parent opentracing.Span,
) opentracing.Span {
	childCtx := parent.(*spanImpl).raw.StandardContext.NewChild()
	return s.startSpanGeneric(
		operationName,
		childCtx,
	)
}

// A helper for spanImpl creation.
func (s *tracerImpl) startSpanGeneric(
	operationName string,
	childCtx *StandardContext,
) opentracing.Span {
	span := &spanImpl{
		tracer:   s,
		recorder: s.recorder,
		raw: RawSpan{
			StandardContext: childCtx,
			Operation:       operationName,
			Start:           time.Now(),
			Duration:        -1,
			Tags:            opentracing.Tags{},
			Logs:            []*opentracing.LogData{},
		},
	}
	return span
}
