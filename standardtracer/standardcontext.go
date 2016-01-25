package standardtracer

import "sync"

type StandardContext struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID int64

	// A probabilistically unique identifier for a span.
	SpanID int64

	// The SpanID of this StandardContext's parent, or 0 if there is no parent.
	ParentSpanID int64

	// Whether the trace is sampled.
	Sampled bool

	// `tagLock` protects the `traceAttrs` map, which in turn supports
	// `SetTraceAttribute` and `TraceAttribute`.
	tagLock    sync.RWMutex
	traceAttrs map[string]string
}

func NewRootStandardContext() *StandardContext {
	return &StandardContext{
		TraceID:    randomID(),
		SpanID:     randomID(),
		Sampled:    randomID()%64 == 0,
		traceAttrs: make(map[string]string),
	}
}

func (c *StandardContext) NewChild() *StandardContext {
	c.tagLock.RLock()
	newTags := make(map[string]string, len(c.traceAttrs))
	for k, v := range c.traceAttrs {
		newTags[k] = v
	}
	c.tagLock.RUnlock()

	return &StandardContext{
		TraceID:      c.TraceID,
		SpanID:       randomID(),
		ParentSpanID: c.SpanID,
		Sampled:      c.Sampled,
		traceAttrs:   newTags,
	}
}
