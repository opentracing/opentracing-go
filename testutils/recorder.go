package testutils

import (
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/standardtracer"
)

// InMemoryRecorder is a simple thread-safe implementation of
// standardtracer.Recorder that stores all reported spans in memory, accessible
// via reporter.GetSpans()
type InMemoryRecorder struct {
	processName string
	spans       []*standardtracer.RawSpan
	tags        opentracing.Tags
	lock        sync.Mutex
}

// NewInMemoryRecorder instantiates a new InMemoryRecorder with the given `processName`
func NewInMemoryRecorder(processName string) *InMemoryRecorder {
	return &InMemoryRecorder{
		processName: processName,
		spans:       make([]*standardtracer.RawSpan, 0),
		tags:        make(opentracing.Tags),
	}
}

// ProcessName implements ProcessName() of standardtracer.Recorder
func (recorder *InMemoryRecorder) ProcessName() string {
	return recorder.processName
}

// SetTag implements SetTag() of standardtracer.Recorder. Tags can be
// retrieved via recorder.GetTags()
func (recorder *InMemoryRecorder) SetTag(key string, val interface{}) standardtracer.ProcessIdentifier {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	recorder.tags[key] = val
	return recorder
}

// RecordSpan implements RecordSpan() of standardtracer.Recorder.
//
// The recorded spans can be retrieved via recorder.Spans slice.
func (recorder *InMemoryRecorder) RecordSpan(span *standardtracer.RawSpan) {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	recorder.spans = append(recorder.spans, span)
}

// GetSpans returns a snapshot of spans recorded so far.
func (recorder *InMemoryRecorder) GetSpans() []*standardtracer.RawSpan {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	spans := make([]*standardtracer.RawSpan, len(recorder.spans))
	copy(spans, recorder.spans)
	return spans
}

// GetTags returns a snapshot of tags.
func (recorder *InMemoryRecorder) GetTags() opentracing.Tags {
	recorder.lock.Lock()
	defer recorder.lock.Unlock()
	tags := make(opentracing.Tags)
	for k, v := range recorder.tags {
		tags[k] = v
	}
	return tags
}
