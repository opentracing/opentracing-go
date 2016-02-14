package testutils

import (
	"sync"

	"github.com/opentracing/opentracing-go/standardtracer"
)

// InMemoryRecorder is a simple thread-safe implementation of
// standardtracer.SpanRecorder that stores all reported spans in memory, accessible
// via reporter.GetSpans()
type InMemoryRecorder struct {
	spans []standardtracer.RawSpan
	lock  sync.Mutex
}

// NewInMemoryRecorder instantiates a new InMemoryRecorder for testing purposes.
func NewInMemoryRecorder() *InMemoryRecorder {
	return &InMemoryRecorder{
		spans: make([]standardtracer.RawSpan, 0),
	}
}

// RecordSpan implements RecordSpan() of standardtracer.SpanRecorder.
//
// The recorded spans can be retrieved via recorder.Spans slice.
func (recorder *InMemoryRecorder) RecordSpan(span standardtracer.RawSpan) {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	recorder.spans = append(recorder.spans, span)
}

// GetSpans returns a snapshot of spans recorded so far.
func (recorder *InMemoryRecorder) GetSpans() []standardtracer.RawSpan {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	spans := make([]standardtracer.RawSpan, len(recorder.spans))
	copy(spans, recorder.spans)
	return spans
}
