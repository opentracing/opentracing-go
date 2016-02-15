package standardtracer

import (
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

// Options allows creating a customized Tracer via NewWithOptions. The object
// must not be updated when there is an active tracer using it.
type Options struct {
	// ShouldSample is a function which is called when creating a new Span and
	// determines whether that Span is sampled. The randomized TraceID is supplied
	// to allow deterministic sampling decisions to be made across different nodes.
	// For example,
	//
	//   func(traceID int64) { return traceID % 64 == 0 }
	//
	// samples every 64th trace on average.
	ShouldSample func(int64) bool
	// TrimUnsampledSpans turns potentially expensive operations on unsampled
	// Spans into no-ops. More precisely, tags, attributes and log events
	// are silently discarded.
	TrimUnsampledSpans bool
	// Recorder receives Spans which have been finished.
	Recorder SpanRecorder
}

// DefaultOptions returns an Options object with a 1 in 64 sampling rate and
// all options disabled. A Recorder needs to be set manually before using the
// returned object with a Tracer.
func DefaultOptions() Options {
	var opts Options
	opts.ShouldSample = func(traceID int64) bool { return traceID%64 == 0 }
	return opts
}

// NewWithOptions creates a customized Tracer.
func NewWithOptions(opts Options) opentracing.Tracer {
	rval := &tracerImpl{
		Options: opts,
		spanPool: sync.Pool{New: func() interface{} {
			return &spanImpl{}
		}},
	}
	rval.textPropagator = &splitTextPropagator{rval}
	rval.binaryPropagator = &splitBinaryPropagator{rval}
	rval.goHTTPPropagator = &goHTTPPropagator{rval.binaryPropagator}
	return rval
}

// New creates and returns a standard Tracer which defers completed Spans to
// `recorder`.
// Spans created by this Tracer support the ext.SamplingPriority tag: Calling
// SetTag(ext.SamplingPriority, nil) causes the Span to be Sampled from that
// point on.
func New(recorder SpanRecorder) opentracing.Tracer {
	opts := DefaultOptions()
	opts.Recorder = recorder
	return NewWithOptions(opts)
}

// Implements the `Tracer` interface.
type tracerImpl struct {
	Options
	spanPool         sync.Pool
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

func (t *tracerImpl) getSpan() *spanImpl {
	sp := t.spanPool.Get().(*spanImpl)
	sp.reset()
	return sp
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
	sp := t.getSpan()
	if opts.Parent == nil {
		sp.raw.TraceID, sp.raw.SpanID = randomID2()
		sp.raw.Sampled = t.ShouldSample(sp.raw.TraceID)
	} else {
		pr := opts.Parent.(*spanImpl)
		sp.raw.TraceID = pr.raw.TraceID
		sp.raw.SpanID = randomID()
		sp.raw.ParentSpanID = pr.raw.SpanID
		sp.raw.Sampled = pr.raw.Sampled

		pr.Lock()
		if l := len(pr.raw.Attributes); l > 0 {
			sp.raw.Attributes = make(map[string]string, len(pr.raw.Attributes))
			for k, v := range pr.raw.Attributes {
				sp.raw.Attributes[k] = v
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
