package opentracing

// WireEncoding identifies a WireInjector/WireExtractor pair.
type WireEncoding int

const (
	// WIRE_ENCODING_SPLIT_BINARY divides the encoded Span into two components:
	//
	//   1) The "tracer state" for the Span (for example, in Dapper this would
	//      include a trace_id, a span_id, and a bitmask representing the
	//      sampling status for the given trace)
	//   2) Any Trace Attributes (per Span.SetTraceAttribute)
	//
	// The encoded data is separated in this way for a variety of reasons; the
	// most important is to give OpenTracing users a portable way to opt out of
	// Trace Attribute propagation entirely if they deem it a stability risk.
	//
	// The `carrier` for injection must be two pointers to `[]byte`s for the
	// two components described above, each of which is represented as
	// arbitrary binary data. The `carrier` for extraction must be two
	// (non-pointer) `[]byte`s respectively.
	WIRE_ENCODING_SPLIT_BINARY WireEncoding = iota

	// WIRE_ENCODING_SPLIT_TEXT divides the encoded Span into two components:
	//
	//   1) The "tracer state" for the Span (for example, in Dapper this would
	//      include a trace_id, a span_id, and a bitmask representing the
	//      sampling status for the given trace)
	//   2) Any Trace Attributes (per Span.SetTraceAttribute)
	//
	// The encoded data is separated in this way for a variety of reasons; the
	// most important is to give OpenTracing users a portable way to opt out of
	// Trace Attribute propagation entirely if they deem it a stability risk.
	//
	// The `carrier` for injection must be two pointers to `map[string]string`s
	// for the two components described above. The `carrier` for extraction
	// must be two (non-pointer) `map[string]string`s respectively.
	WIRE_ENCODING_SPLIT_TEXT
)

// WireInjector is responsible for injecting Span instances in a manner suitable
// for propagation in an WireEncoding-specific "carrier" object or objects.
// Typically the injection will take place across an RPC boundary, but message
// queues and other IPC mechanisms are also reasonable places to use a
// WireInjector.
//
// The specific encoding for an injected Span depends on the WireEncoding.
// OpenTracing defines a common set of WireEncodings, and each has an expected
// carrier type and format. See the WireEncoding enum comments for details.
//
// See WireExtractor and Span.WireInjectorForEncoding.
type WireInjector interface {
	InjectSpan(span Span, carrier ...interface{})
}

// WireExtractor is responsible for extracting Span instances from an
// encoding-specific "carrier" object. Typically the extraction will take place
// on the server side of an RPC boundary, but message queues and other IPC
// mechanisms are also reasonable places to use a WireExtractor.
//
// The specific encoding for a Span extraction depends on the WireEncoding.
// OpenTracing defines a common set of WireEncodings, and each has an expected
// carrier type and format. See the WireEncoding enum comments for details.
//
// See WireInjector and Tracer.WireExtractorForEncoding.
type WireExtractor interface {
	ExtractSpan(operationName string, carrier ...interface{}) (Span, error)
}
