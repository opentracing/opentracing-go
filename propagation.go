package opentracing

///////////////////////////////////////////////////////////////////////////////
// CORE PROPAGATION INTERFACES:
///////////////////////////////////////////////////////////////////////////////

// PropagationInjector is responsible for injecting Span instances in a manner
// suitable for propagation via a format-specific "carrier" object. Typically
// the injection will take place across an RPC boundary, but message queues and
// other IPC mechanisms are also reasonable places to use a
// PropagationInjector.
//
// See PropagationExtractor and Span.PropagationInjectorForFormat.
type PropagationInjector interface {
	InjectSpan(span Span, carrier interface{})
}

// PropagationExtractor is responsible for extracting Span instances from an
// format-specific "carrier" object. Typically the extraction will take place
// on the server side of an RPC boundary, but message queues and other IPC
// mechanisms are also reasonable places to use a PropagationExtractor.
//
// See PropagationInjector and Tracer.PropagationExtractorForFormat.
type PropagationExtractor interface {
	ExtractSpan(operationName string, carrier interface{}) (Span, error)
}

///////////////////////////////////////////////////////////////////////////////
// BUILTIN PROPAGATION FORMATS:
///////////////////////////////////////////////////////////////////////////////

// BuiltinPropagationFormat is the shared type for builtin OpenTracing in-band
// propagation formats.
type BuiltinPropagationFormat int

const (
	// PROPAGATION_FORMAT_SPLIT_BINARY encodes the Span in a BinaryCarrier
	// instance.
	//
	// The `carrier` for injection and extraction must be a `*BinaryCarrier`
	// instance.
	PROPAGATION_FORMAT_SPLIT_BINARY BuiltinPropagationFormat = iota

	// PROPAGATION_FORMAT_SPLIT_BINARY encodes the Span in a TextCarrier
	// instance.
	//
	// The `carrier` for injection and extraction must be a `*TextCarrier`
	// instance.
	PROPAGATION_FORMAT_SPLIT_TEXT

	// PROPAGATION_FORMAT_GO_HTTP_HEADER encodes the Span into a Go http.Header
	// instance (both the tracer state and any Trace Attributes).
	//
	// The `carrier` for both injection and extraction must be an http.Header
	// instance.
	PROPAGATION_FORMAT_GO_HTTP_HEADER
)

// TextCarrier breaks a propagated Span into two pieces.
//
// The Span is separated in this way for a variety of reasons; the most
// important is to give OpenTracing users a portable way to opt out of Trace
// Attribute propagation entirely if they deem it a stability risk.
type TextCarrier struct {
	// TracerState is Tracer-specific context that must cross process
	// boundaries. For example, in Dapper this would include a trace_id, a
	// span_id, and a bitmask representing the sampling status for the given
	// trace.
	TracerState map[string]string

	// Any Trace Attributes for the encoded Span (per Span.SetTraceAttribute).
	TraceAttributes map[string]string
}

func NewTextCarrier() *TextCarrier {
	return &TextCarrier{
		TracerState:     make(map[string]string),
		TraceAttributes: make(map[string]string),
	}
}

// BinaryCarrier breaks a propagated Span into two pieces.
//
// The Span is separated in this way for a variety of reasons; the most
// important is to give OpenTracing users a portable way to opt out of Trace
// Attribute propagation entirely if they deem it a stability risk.
type BinaryCarrier struct {
	// TracerState is Tracer-specific context that must cross process
	// boundaries. For example, in Dapper this would include a trace_id, a
	// span_id, and a bitmask representing the sampling status for the given
	// trace.
	TracerState []byte

	// Any Trace Attributes for the encoded Span (per Span.SetTraceAttribute).
	TraceAttributes []byte
}

func NewBinaryCarrier() *BinaryCarrier {
	return &BinaryCarrier{
		TracerState:     make([]byte),
		TraceAttributes: make([]byte),
	}
}
