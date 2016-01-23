package opentracing

// PropagationEncoder is a simple interface to encode a Span *for propagation*
// as a binary byte array or a string-to-string text map.
//
// The encoder should only represent state from the Span that needs to cross
// process boundaries: e.g., a unique id for the larger distributed trace, a
// unique id for the span itself, and any trace annotations (per
// Span.SetTraceAttribute).
//
// The encoded form of a propagate span is divided into two components: the
// core identifying information for the Span and any trace attributes.
type PropagationEncoder interface {
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
}

// PropagationDecoder is a simple interface to decode a binary byte
// array or a string-to-string text map into a TraceContext.
type PropagationDecoder interface {
	// Converts the encoded binary data (see
	// PropagationEncoder.TraceContextToBinary()) into a TraceContext.
	//
	// The first parameter contains the encoder's serialization of the core
	// identifying information in a TraceContext instance.
	//
	// The second parameter contains the encoder's serialization of the trace
	// attributes (per `SetTraceAttribute` and `TraceAttribute`) attached to a
	// TraceContext instance.
	NewSpanFromBinary(
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
	NewSpanFromText(
		operationName string,
		traceContextID map[string]string,
		traceAttrs map[string]string,
	) (Span, error)
}
