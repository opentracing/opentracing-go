// Noop implementations of the core opentracing interfaces.
package opentracing

import "golang.org/x/net/context"

type noopTraceContextID struct{}
type noopSpan struct{}
type noopRecorder struct{}
type noopTraceContextIDSource struct{}
type noopOpenTracer struct {
	noopTraceContextIDSource
}

var (
	defaultNoopTraceContextID       = noopTraceContextID{}
	defaultNoopSpan                 = noopSpan{}
	defaultNoopRecorder             = noopRecorder{}
	defaultNoopTraceContextIDSource = noopTraceContextIDSource{}
	defaultNoopOpenTracer           = noopOpenTracer{}
	emptyTags                       = Tags{}
	emptyBytes                      = []byte{}
)

// noopTraceContextID:

func (n noopTraceContextID) NewChild() (TraceContextID, Tags) {
	return defaultNoopTraceContextID, emptyTags
}
func (n noopTraceContextID) Serialize() []byte {
	return emptyBytes
}

// noopSpan:
func (n noopSpan) StartChildSpan(operationName string, parent ...context.Context) (Span, context.Context) {
	if len(parent) > 0 {
		return defaultNoopSpan, parent[0]
	} else {
		return defaultNoopSpan, context.Background()
	}
}
func (n noopSpan) SetTag(key string, value interface{})           {}
func (n noopSpan) Info(message string, payload ...interface{})    {}
func (n noopSpan) Warning(message string, payload ...interface{}) {}
func (n noopSpan) Error(message string, payload ...interface{})   {}
func (n noopSpan) Finish()                                        {}
func (n noopSpan) TraceContext() TraceContextID                   { return defaultNoopTraceContextID }

// noopTraceContextIDSource:
func (n noopTraceContextIDSource) DeserializeTraceContextID(encoded []byte) (TraceContextID, error) {
	return defaultNoopTraceContextID, nil
}
func (n noopTraceContextIDSource) NewRootTraceContextID() TraceContextID {
	return defaultNoopTraceContextID
}

// noopRecorder:
func (n noopRecorder) SetTag(key string, val interface{}) {}
func (n noopRecorder) RecordSpan(span *RawSpan)           {}
func (n noopRecorder) ProcessName() string                { return "" }

// noopOpenTracer:
func (n noopOpenTracer) StartSpan(operationName string, parent ...interface{}) (Span, context.Context) {
	if len(parent) > 0 {
		if ctx, ok := parent[0].(context.Context); ok {
			return defaultNoopSpan, ctx
		}
	}
	ctx := context.Background()
	return defaultNoopSpan, ctx
}
