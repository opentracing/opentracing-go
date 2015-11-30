// Noop implementations of the core opentracing interfaces.
package opentracing

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
)

const (
	emptyString = ""
)

// noopTraceContextID:

func (n noopTraceContextID) NewChild() (TraceContextID, Tags) {
	return defaultNoopTraceContextID, emptyTags
}
func (n noopTraceContextID) SerializeBinary() []byte {
	return emptyBytes
}
func (n noopTraceContextID) SerializeString() string {
	return emptyString
}

// noopSpan:
func (n noopSpan) StartChild(operationName string, initialTags ...Tags) Span {
	return defaultNoopSpan
}
func (n noopSpan) SetTag(key string, value interface{}) Span      { return n }
func (n noopSpan) Info(message string, payload ...interface{})    {}
func (n noopSpan) Warning(message string, payload ...interface{}) {}
func (n noopSpan) Error(message string, payload ...interface{})   {}
func (n noopSpan) Finish()                                        {}
func (n noopSpan) TraceContext() *TraceContext                    { return defaultNoopTraceContext }

// noopTraceContextIDSource:
func (n noopTraceContextIDSource) DeserializeBinaryTraceContextID(encoded []byte) (TraceContextID, error) {
	return defaultNoopTraceContextID, nil
}
func (n noopTraceContextIDSource) DeserializeStringTraceContextID(encoded string) (TraceContextID, error) {
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
func (n noopOpenTracer) StartNewTrace(operationName string, initialTags ...Tags) Span {
	return defaultNoopSpan
}

func (n noopOpenTracer) JoinTrace(operationName string, parent interface{}, initialTags ...Tags) Span {
	return defaultNoopSpan
}
