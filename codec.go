package opentracing

// TraceContextEncoder is a simple interface to encode a TraceContext as a
// binary byte array or a string-to-string text map.
type TraceContextEncoder interface {
	// Converts the TraceContext into encoded binary data (see
	// TraceContextDecoder.TraceContextFromBinary()).
	//
	// The first return value must represent the encoder's serialization of
	// the core identifying information in `tc`.
	//
	// The second return value must represent the encoder's serialization of
	// the trace attributes, per `SetTraceAttribute` and `TraceAttribute`.
	TraceContextToBinary(
		tc TraceContext,
	) (
		traceContextID []byte,
		traceAttrs []byte,
	)

	// Converts the TraceContext into a encoded string:string text map (see
	// TraceContextDecoder.TraceContextFromText()).
	//
	// The first return value must represent the encoder's serialization of
	// the core identifying information in `tc`.
	//
	// The second return value must represent the encoder's serialization of
	// the trace attributes, per `SetTraceAttribute` and `TraceAttribute`.
	TraceContextToText(
		tc TraceContext,
	) (
		traceContextID map[string]string,
		traceAttrs map[string]string,
	)
}

// TraceContextDecoder is a simple interface to decode a binary byte
// array or a string-to-string text map into a TraceContext.
type TraceContextDecoder interface {
	// Converts the encoded binary data (see
	// TraceContextEncoder.TraceContextToBinary()) into a TraceContext.
	//
	// The first parameter contains the encoder's serialization of the core
	// identifying information in a TraceContext instance.
	//
	// The second parameter contains the encoder's serialization of the trace
	// attributes (per `SetTraceAttribute` and `TraceAttribute`) attached to a
	// TraceContext instance.
	TraceContextFromBinary(
		traceContextID []byte,
		traceAttrs []byte,
	) (TraceContext, error)

	// Converts the encoded string:string text map (see
	// TraceContextEncoder.TraceContextToText()) into a TraceContext.
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
	TraceContextFromText(
		traceContextID map[string]string,
		traceAttrs map[string]string,
	) (TraceContext, error)
}
