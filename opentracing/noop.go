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
	defaultNoopTraceContext         = newTraceContext(defaultNoopTraceContextID, nil)
	emptyTags                       = Tags{}
	emptyBytes                      = []byte{}
	emptyStringMap                  = map[string]string{}
)

// noopTraceContextID:

func (n noopTraceContextID) NewChild() (TraceContextID, Tags) {
	return defaultNoopTraceContextID, emptyTags
}

// noopSpan:
func (n noopSpan) StartChild(operationName string, keyValueTags ...interface{}) Span {
	return defaultNoopSpan
}
func (n noopSpan) SetTag(key string, value interface{}) Span      { return n }
func (n noopSpan) Info(message string, payload ...interface{})    {}
func (n noopSpan) Warning(message string, payload ...interface{}) {}
func (n noopSpan) Error(message string, payload ...interface{})   {}
func (n noopSpan) Finish()                                        {}
func (n noopSpan) TraceContext() *TraceContext                    { return defaultNoopTraceContext }
func (n noopSpan) AddToGoContext(ctx context.Context) (Span, context.Context) {
	return n, GoContextWithSpan(ctx, n)
}

// noopTraceContextIDSource:
func (n noopTraceContextIDSource) MarshalBinaryTraceContextID(tcid TraceContextID) []byte {
	return emptyBytes
}
func (n noopTraceContextIDSource) MarshalStringMapTraceContextID(tcid TraceContextID) map[string]string {
	return emptyStringMap
}
func (n noopTraceContextIDSource) UnmarshalBinaryTraceContextID(encoded []byte) (TraceContextID, error) {
	return defaultNoopTraceContextID, nil
}
func (n noopTraceContextIDSource) UnmarshalStringMapTraceContextID(encoded map[string]string) (TraceContextID, error) {
	return defaultNoopTraceContextID, nil
}
func (n noopTraceContextIDSource) NewRootTraceContextID() TraceContextID {
	return defaultNoopTraceContextID
}

// noopRecorder:
func (n noopRecorder) SetTag(key string, val interface{}) ProcessRecorder { return n }
func (n noopRecorder) RecordSpan(span *RawSpan)                           {}
func (n noopRecorder) ProcessName() string                                { return "" }

// noopOpenTracer:
func (n noopOpenTracer) StartTrace(operationName string, keyValueTags ...interface{}) Span {
	return defaultNoopSpan
}

func (n noopOpenTracer) JoinTrace(operationName string, parent interface{}, keyValueTags ...interface{}) Span {
	return defaultNoopSpan
}
