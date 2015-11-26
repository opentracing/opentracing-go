// Noop implementations of the core opentracing interfaces.
package opentracing

import "golang.org/x/net/context"

// noopContextID:
type noopContextID struct{}

func (n noopContextID) NewChild() (ContextID, Tags) {
	return noopContextID{}, Tags{}
}
func (n noopContextID) Serialize() []byte {
	return []byte{}
}

// noopSpan:
type noopSpan struct{}

func (n noopSpan) StartChildSpan(operationName string, parent ...context.Context) (Span, context.Context) {
	if len(parent) > 0 {
		return noopSpan{}, parent[0]
	} else {
		return noopSpan{}, context.Background()
	}
}
func (n noopSpan) SetTag(key string, value interface{})           {}
func (n noopSpan) Info(message string, payload ...interface{})    {}
func (n noopSpan) Warning(message string, payload ...interface{}) {}
func (n noopSpan) Error(message string, payload ...interface{})   {}
func (n noopSpan) Finish()                                        {}
func (n noopSpan) ContextID() ContextID                           { return noopContextID{} }

// noopContextIDSource:
type noopContextIDSource struct{}

func (n noopContextIDSource) DeserializeContextID(encoded []byte) (ContextID, error) {
	return noopContextID{}, nil
}
func (n noopContextIDSource) NewRootContextID() ContextID {
	return noopContextID{}
}

// noopRecorder:
type noopRecorder struct{}

func (n noopRecorder) SetTag(key string, val interface{}) {}
func (n noopRecorder) RecordSpan(span *RawSpan)           {}

// noopOpenTracer:
type noopOpenTracer struct {
	noopContextIDSource
}

func (n noopOpenTracer) StartSpan(operationName string, parent ...interface{}) (Span, context.Context) {
	if len(parent) > 0 {
		if ctx, ok := parent[0].(context.Context); ok {
			return noopSpan{}, ctx
		}
	}
	ctx := context.Background()
	return noopSpan{}, ctx
}
