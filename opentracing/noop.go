package opentracing

import "golang.org/x/net/context"

type noopTraceContext struct{}
type noopSpan struct{}
type noopRecorder struct{}
type noopTraceContextSource struct{}
type noopOpenTracer struct {
	noopTraceContextSource
}

var (
	defaultNoopTraceContext       = noopTraceContext{}
	defaultNoopSpan               = noopSpan{}
	defaultNoopRecorder           = noopRecorder{}
	defaultNoopTraceContextSource = noopTraceContextSource{}
	defaultNoopOpenTracer         = noopOpenTracer{}
	emptyTags                     = Tags{}
	emptyBytes                    = []byte{}
	emptyStringMap                = map[string]string{}
)

const (
	emptyString = ""
)

// noopTraceContext:

func (n noopTraceContext) NewChild() (TraceContext, Tags) {
	return defaultNoopTraceContext, emptyTags
}
func (n noopTraceContext) SetTraceTag(key, val string) TraceContext { return n }
func (n noopTraceContext) TraceTag(key string) string               { return emptyString }

// noopSpan:
func (n noopSpan) StartChild(operationName string, keyValueTags ...interface{}) Span {
	return defaultNoopSpan
}
func (n noopSpan) SetTag(key string, value interface{}) Span      { return n }
func (n noopSpan) Info(message string, payload ...interface{})    {}
func (n noopSpan) Warning(message string, payload ...interface{}) {}
func (n noopSpan) Error(message string, payload ...interface{})   {}
func (n noopSpan) Finish()                                        {}
func (n noopSpan) TraceContext() TraceContext                     { return defaultNoopTraceContext }
func (n noopSpan) AddToGoContext(ctx context.Context) (Span, context.Context) {
	return n, GoContextWithSpan(ctx, n)
}

// noopTraceContextSource:
func (n noopTraceContextSource) MarshalBinaryTraceContext(tcid TraceContext) []byte {
	return emptyBytes
}
func (n noopTraceContextSource) MarshalStringMapTraceContext(tcid TraceContext) map[string]string {
	return emptyStringMap
}
func (n noopTraceContextSource) UnmarshalBinaryTraceContext(encoded []byte) (TraceContext, error) {
	return defaultNoopTraceContext, nil
}
func (n noopTraceContextSource) UnmarshalStringMapTraceContext(encoded map[string]string) (TraceContext, error) {
	return defaultNoopTraceContext, nil
}
func (n noopTraceContextSource) NewRootTraceContext() TraceContext {
	return defaultNoopTraceContext
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
