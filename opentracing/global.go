package opentracing

import (
	"golang.org/x/net/context"
)

var (
	globalOpenTracer OpenTracer = noopOpenTracer{noopContextIDSource{}}
)

// Should be called as early as possible in main(), prior to calling the
// `StartSpan*` (etc) global funcs below. Prior to calling `InitGlobal`, any
// Spans started via the `StartSpan*` globals are noops.
func InitGlobalTracer(rec ComponentRecorder, ctxIDSource ContextIDSource) {
	globalOpenTracer = NewStandardTracer(rec, ctxIDSource)
}

func GlobalTracer() OpenTracer {
	return globalOpenTracer
}

// See `OpenTracer.StartSpan` and `InitGlobal()`.
func StartSpan(operationName string, parent ...interface{}) (Span, context.Context) {
	return globalOpenTracer.StartSpan(operationName, parent...)
}

// See `ContextIDSource.DeserializeContextID()`
func DeserializeContextID(encoded []byte) (ContextID, error) {
	return globalOpenTracer.DeserializeContextID(encoded)
}
