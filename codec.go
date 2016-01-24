package opentracing

// SpanPropagator is an interface both to encode a Span *for propagation* (in-band
// with application data) and to join back to that trace given the encoded
// data. There are two supported encodings: binary byte arrays and
// string-to-string text maps.
//
// The encoded form of a propagated span is divided into two components: the
// core identifying information for the Span and any trace attributes. (These
// are separated for a variety of reasons, though the most important is to
// allow OpenTracing users to opt out of trace attribute propagation entirely
// if they deem it a stability risk)
type SpanPropagator interface {
	// Represents the Span for propagation as encoded binary data (see
	// PropagationDecoder.NewSpanFromBinary()).
	//
	// The first return value must represent the encoder's serialization of
	// core identifying information in `tc`.
	//
	// The second return value must represent the encoder's serialization of
	// trace attributes, per `SetTraceAttribute` and `TraceAttribute`.
	PropagateSpanAsBinary(
		sp Span,
	) (
		traceContextID []byte,
		traceAttrs []byte,
	)

	// XXX XXX
	// Converts the TraceContext into a encoded string:string text map (see
	// PropagationDecoder.TraceContextFromText()).
	//
	// The first return value must represent the encoder's serialization of
	// the core identifying information in `tc`.
	//
	// The second return value must represent the encoder's serialization of
	// the trace attributes, per `SetTraceAttribute` and `TraceAttribute`.
	PropagateSpanAsText(
		sp Span,
	) (
		traceContextID map[string]string,
		traceAttrs map[string]string,
	)

	// Converts the encoded binary data (see
	// PropagationEncoder.TraceContextToBinary()) into a TraceContext.
	//
	// The first parameter contains the encoder's serialization of the core
	// identifying information in a TraceContext instance.
	//
	// The second parameter contains the encoder's serialization of the trace
	// attributes (per `SetTraceAttribute` and `TraceAttribute`) attached to a
	// TraceContext instance.
	JoinTraceFromBinary(
		operationName string,
		traceContextID []byte,
		traceAttrs []byte,
	) (Span, error)

	// Converts the encoded string:string text map (see
	// PropagationEncoder.TraceContextToText()) into a TraceContext.
	//
	// The first parameter contains the encoder's serialization of the core
	// identifying information in a TraceContext instance.
	//
	// The second parameter contains the encoder's serialization of the trace
	// attributes (per `SetTraceAttribute` and `TraceAttribute`) attached to a
	// TraceContext instance.
	//
	// It's permissible to pass the same map to both parameters (e.g., an HTTP
	// request headers map): the implementation should only decode the subset
	// it's interested in.
	JoinTraceFromText(
		operationName string,
		traceContextID map[string]string,
		traceAttrs map[string]string,
	) (Span, error)
}
