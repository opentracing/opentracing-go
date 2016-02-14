package standardtracer

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

// New creates and returns a standard Tracer which defers to `recorder` after
// RawSpans have been assembled.
func New(recorder SpanRecorder) opentracing.Tracer {
	rval := &tracerImpl{
		recorder: recorder,
	}
	rval.textPropagator = &splitTextPropagator{rval}
	rval.binaryPropagator = &splitBinaryPropagator{rval}
	rval.goHTTPPropagator = &goHTTPPropagator{rval.binaryPropagator}
	return rval
}

// Implements the `Tracer` interface.
type tracerImpl struct {
	recorder         SpanRecorder
	textPropagator   *splitTextPropagator
	binaryPropagator *splitBinaryPropagator
	goHTTPPropagator *goHTTPPropagator
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

	// Build the new span. This is the only allocation: We'll return this as
	// a opentracing.Span.
	sp := &spanImpl{}
	if opts.Parent == nil {
		sp.raw.TraceID, sp.raw.SpanID = randomID2()
		sp.raw.Sampled = sp.raw.TraceID%64 == 0
	} else {
		pr := opts.Parent.(*spanImpl)
		sp.raw.TraceID = pr.raw.TraceID
		sp.raw.SpanID = randomID()
		sp.raw.ParentSpanID = pr.raw.SpanID
		sp.raw.Sampled = pr.raw.Sampled

		pr.Lock()
		if l := len(pr.traceAttrs); l > 0 {
			sp.traceAttrs = make(map[string]string, len(pr.traceAttrs))
			for k, v := range pr.traceAttrs {
				sp.traceAttrs[k] = v
			}
		}
		pr.Unlock()
	}

	return t.startSpanInternal(
		sp,
		opts.OperationName,
		startTime,
		tags,
	)
}

func (t *tracerImpl) startSpanInternal(
	sp *spanImpl,
	operationName string,
	startTime time.Time,
	tags opentracing.Tags,
) opentracing.Span {
	sp.tracer = t
	sp.raw.Operation = operationName
	sp.raw.Start = startTime
	sp.raw.Duration = -1
	sp.raw.Tags = tags
	return sp
}

func (t *tracerImpl) Extractor(format interface{}) opentracing.Extractor {
	switch format {
	case opentracing.SplitText:
		return t.textPropagator
	case opentracing.SplitBinary:
		return t.binaryPropagator
	case opentracing.GoHTTPHeader:
		return t.goHTTPPropagator
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
		return t.goHTTPPropagator
	}
	return nil
}
