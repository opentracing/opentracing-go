package opentracing

import "golang.org/x/net/context"

type noopTraceContext struct{}
type noopSpan struct{}
type noopRecorder struct{}
type noopTraceContextSource struct{}
type noopTracer struct {
	noopTraceContextSource
}

var (
	defaultNoopTraceContext       = noopTraceContext{}
	defaultNoopSpan               = noopSpan{}
	defaultNoopRecorder           = noopRecorder{}
	defaultNoopTraceContextSource = noopTraceContextSource{}
	defaultNoopTracer             = noopTracer{}
	emptyTags                     = Tags{}
	emptyBytes                    = []byte{}
	emptyStringMap                = map[string]string{}
)

const (
	emptyString = ""
)

// noopTraceContext:

func (n noopTraceContext) SetTraceAttribute(key, val string) TraceContext { return n }
func (n noopTraceContext) TraceAttribute(key string) string               { return emptyString }

// noopSpan:
func (n noopSpan) StartChild(operationName string) Span {
	return defaultNoopSpan
}
func (n noopSpan) SetTag(key string, value interface{}) Span      { return n }
func (n noopSpan) SetTags(tags Tags) Span                         { return n }
func (n noopSpan) Info(message string, payload ...interface{})    {}
func (n noopSpan) Warning(message string, payload ...interface{}) {}
func (n noopSpan) Error(message string, payload ...interface{})   {}
func (n noopSpan) Finish()                                        {}
func (n noopSpan) TraceContext() TraceContext                     { return defaultNoopTraceContext }
func (n noopSpan) AddToGoContext(ctx context.Context) (Span, context.Context) {
	return n, GoContextWithSpan(ctx, n)
}

// noopTraceContextSource:
func (n noopTraceContextSource) MarshalTraceContextBinary(tcid TraceContext) ([]byte, []byte) {
	return emptyBytes, emptyBytes
}
func (n noopTraceContextSource) MarshalTraceContextStringMap(tcid TraceContext) (map[string]string, map[string]string) {
	return emptyStringMap, emptyStringMap
}
func (n noopTraceContextSource) UnmarshalTraceContextBinary(
	traceContextID []byte,
	traceAttrs []byte,
) (TraceContext, error) {
	return defaultNoopTraceContext, nil
}
func (n noopTraceContextSource) UnmarshalTraceContextStringMap(
	traceContextID map[string]string,
	traceAttrs map[string]string,
) (TraceContext, error) {
	return defaultNoopTraceContext, nil
}
func (n noopTraceContextSource) NewRootTraceContext() TraceContext {
	return defaultNoopTraceContext
}
func (n noopTraceContextSource) NewChildTraceContext(parent TraceContext) (TraceContext, Tags) {
	return defaultNoopTraceContext, emptyTags
}

// noopRecorder:
func (n noopRecorder) SetTag(key string, val interface{}) ProcessIdentifier { return n }
func (n noopRecorder) RecordSpan(span *RawSpan)                             {}
func (n noopRecorder) ProcessName() string                                  { return "" }

// noopTracer:
func (n noopTracer) StartTrace(operationName string) Span {
	return defaultNoopSpan
}

func (n noopTracer) StartSpanWithContext(operationName string, ctx TraceContext) Span {
	return defaultNoopSpan
}

func (n noopTracer) JoinTrace(operationName string, parent interface{}) Span {
	return defaultNoopSpan
}
