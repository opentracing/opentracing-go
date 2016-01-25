package standardtracer

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

// New creates and returns a standard Tracer which defers to `recorder` and
// `source` as appropriate.
func New(recorder Recorder) opentracing.Tracer {
	return &tracerImpl{
		recorder: recorder,
	}
}

// Implements the `Tracer` interface.
type tracerImpl struct {
	recorder Recorder
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
			Logs:            []*RawLog{},
		},
	}
	return span
}
