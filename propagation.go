package opentracing

///////////////////////////////////////////////////////////////////////////////
// CORE PROPAGATION INTERFACES:
///////////////////////////////////////////////////////////////////////////////

// Injector is responsible for injecting Span instances in a manner suitable
// for propagation via a format-specific "carrier" object. Typically the
// injection will take place across an RPC boundary, but message queues and
// other IPC mechanisms are also reasonable places to use a Injector.
//
// See Extractor and Span.Injector.
type Injector interface {
	InjectSpan(span Span, carrier interface{}) error
}

// Extractor is responsible for extracting Span instances from an
// format-specific "carrier" object. Typically the extraction will take place
// on the server side of an RPC boundary, but message queues and other IPC
// mechanisms are also reasonable places to use a Extractor.
//
// See Injector and Tracer.Extractor.
type Extractor interface {
	JoinTrace(operationName string, carrier interface{}) (Span, error)
}

///////////////////////////////////////////////////////////////////////////////
// BUILTIN PROPAGATION FORMATS:
///////////////////////////////////////////////////////////////////////////////

// BuiltinFormat is used to demarcate the values within package `opentracing`
// that are intended for use with the Tracer.Injector() and Tracer.Extractor()
// methods.
type BuiltinFormat byte

const (
	// SplitBinary encodes the Span in a SplitBinaryCarrier instance.
	//
	// The `carrier` for injection and extraction must be a
	// `*SplitBinaryCarrier` instance.
	SplitBinary BuiltinFormat = iota

	// SplitText encodes the Span in a SplitTextCarrier instance.
	//
	// The `carrier` for injection and extraction must be a `*SplitTextCarrier`
	// instance.
	SplitText

	// GoHTTPHeader encodes the Span into a Go http.Header instance (both the
	// tracer state and any Trace Attributes).
	//
	// The `carrier` for both injection and extraction must be an http.Header
	// instance.
	GoHTTPHeader
)

// SplitTextCarrier breaks a propagated Span into two pieces.
//
// The Span is separated in this way for a variety of reasons; the most
// important is to give OpenTracing users a portable way to opt out of Trace
// Attribute propagation entirely if they deem it a stability risk.
type SplitTextCarrier struct {
	// TracerState is Tracer-specific context that must cross process
	// boundaries. For example, in Dapper this would include a trace_id, a
	// span_id, and a bitmask representing the sampling status for the given
	// trace.
	TracerState map[string]string

	// Any Trace Attributes for the encoded Span (per Span.SetTraceAttribute).
	TraceAttributes map[string]string
}

func NewSplitTextCarrier() *SplitTextCarrier {
	return &SplitTextCarrier{
		TracerState:     make(map[string]string),
		TraceAttributes: make(map[string]string),
	}
}

// SplitBinaryCarrier breaks a propagated Span into two pieces.
//
// The Span is separated in this way for a variety of reasons; the most
// important is to give OpenTracing users a portable way to opt out of Trace
// Attribute propagation entirely if they deem it a stability risk.
type SplitBinaryCarrier struct {
	// TracerState is Tracer-specific context that must cross process
	// boundaries. For example, in Dapper this would include a trace_id, a
	// span_id, and a bitmask representing the sampling status for the given
	// trace.
	TracerState []byte

	// Any Trace Attributes for the encoded Span (per Span.SetTraceAttribute).
	TraceAttributes []byte
}

func NewSplitBinaryCarrier() *SplitBinaryCarrier {
	return &SplitBinaryCarrier{
		TracerState:     []byte{},
		TraceAttributes: []byte{},
	}
}
