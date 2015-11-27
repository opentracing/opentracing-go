// Noop implementations of the core opentracing interfaces.
package opentracing

import "golang.org/x/net/context"

type noopContextID struct{}
type noopSpan struct{}
type noopRecorder struct{}
type noopContextIDSource struct{}
type noopOpenTracer struct {
	noopContextIDSource
}

var (
	defaultNoopContextID       = noopContextID{}
	defaultNoopSpan            = noopSpan{}
	defaultNoopRecorder        = noopRecorder{}
	defaultNoopContextIDSource = noopContextIDSource{}
	defaultNoopOpenTracer      = noopOpenTracer{}
	emptyTags                  = Tags{}
	emptyBytes                 = []byte{}
)

// noopContextID:

func (n noopContextID) NewChild() (ContextID, Tags) {
	return defaultNoopContextID, emptyTags
}
func (n noopContextID) Serialize() []byte {
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
func (n noopSpan) ContextID() ContextID                           { return defaultNoopContextID }

// noopContextIDSource:
func (n noopContextIDSource) DeserializeContextID(encoded []byte) (ContextID, error) {
	return defaultNoopContextID, nil
}
func (n noopContextIDSource) NewRootContextID() ContextID {
	return defaultNoopContextID
}

// noopRecorder:
func (n noopRecorder) SetTag(key string, val interface{}) {}
func (n noopRecorder) RecordSpan(span *RawSpan)           {}
func (n noopRecorder) ComponentName() string              { return "" }

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
