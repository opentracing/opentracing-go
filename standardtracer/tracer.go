package standardtracer

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

// New creates and returns a standard Tracer which defers to `recorder` after
// RawSpans have been assembled.
func New(recorder SpanRecorder) opentracing.Tracer {
	rval := &tracerImpl{
		recorder: recorder,
	}
	rval.textPropagator = splitTextPropagator{rval}
	rval.binaryPropagator = splitBinaryPropagator{rval}
	return rval
}

// Implements the `Tracer` interface.
type tracerImpl struct {
	recorder         SpanRecorder
	textPropagator   splitTextPropagator
	binaryPropagator splitBinaryPropagator
}

func (t *tracerImpl) StartSpan(
	operationName string,
) opentracing.Span {
	return t.StartSpanWithOptions(
		opentracing.StartSpanOptions{
			OperationName: operationName,
		})
}

func (t *tracerImpl) StartSpanWithOptions(
	opts opentracing.StartSpanOptions,
) opentracing.Span {
	// Start time.
	startTime := opts.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}

	// Tags.
	tags := opts.Tags
	if tags == nil {
		tags = opentracing.Tags{}
	}

	// The context for the new span.
	var newCtx *StandardContext
	if opts.Parent == nil {
		newCtx = NewRootStandardContext()
	} else {
		newCtx = opts.Parent.(*spanImpl).raw.StandardContext.NewChild()
	}

	return t.startSpanInternal(
		newCtx,
		opts.OperationName,
		startTime,
		tags,
	)
}

func (t *tracerImpl) startSpanInternal(
	newCtx *StandardContext,
	operationName string,
	startTime time.Time,
	tags opentracing.Tags,
) opentracing.Span {
	return &spanImpl{
		tracer:   t,
		recorder: t.recorder,
		raw: RawSpan{
			StandardContext: newCtx,
			Operation:       operationName,
			Start:           startTime,
			Duration:        -1,
			Tags:            tags,
			Logs:            []*opentracing.LogData{},
		},
	}
}

func (t *tracerImpl) Extractor(format interface{}) opentracing.Extractor {
	switch format {
	case opentracing.SplitText:
		return t.textPropagator
	case opentracing.SplitBinary:
		return t.binaryPropagator
	case opentracing.GoHTTPHeader:
	}
	return nil
}

func (t *tracerImpl) Injector(format interface{}) opentracing.Injector {
	switch format {
	case opentracing.SplitText:
		return t.textPropagator
	case opentracing.SplitBinary:
		return t.binaryPropagator
	case opentracing.GoHTTPHeader:
	}
	return nil
}
