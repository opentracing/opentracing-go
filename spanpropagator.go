package opentracing

// SpanPropagator is responsible (a) for encoding Span instances in a manner
// suitable for propagation, and (b) for taking that encoded data and using it
// to generate Span instances that are placed appropriately in the overarching
// Trace. Typically the propagation will take place across an RPC boundary, but
// message queues and other IPC mechanisms are also reasonable places to use
// a SpanPropagator.
//
// The encoded form of a propagated span is divided into two components:
//
//   1) The core identifying information for the Span (for example, in Dapper
//      this would include a trace_id, a span_id, and a bitmask representing
//      the sampling status for the given trace)
//   2) Any trace attributes (per Span.SetTraceAttribute)
//
// The encoded data is separated in this way for a variety of reasons; the
// most important is to give OpenTracing users a portable way to opt out of
// Trace Attribute propagation entirely if they deem it a stability risk.
//
// The PropagateSpanAs*() and JoinTraceFrom*() methods come in two flavors:
// binary and text. The text format is better-suited to pretty-printing and
// debugging, and the binary format is better-suited to compact,
// high-performance encoding, decoding, and transmission.
//
type SpanPropagator interface {
	// Represents the Span for propagation as opaque binary data (see
	// JoinTraceFromBinary()).
	//
	// The first return value (contextSnapshot) must represent the
	// SpanPropagator's encoding of the core identifying information in `sp`.
	//
	// The second return value (traceAttrs) must represent the SpanPropagator's
	// encoding of trace attributes, per `Span.SetTraceAttribute`.
	PropagateSpanAsBinary(
		sp Span,
	) (
		contextSnapshot []byte,
		traceAttrs []byte,
	)

	// Represents the Span for propagation as string:string text maps (see
	// JoinTraceFromText()).
	//
	// The first return value (contextSnapshot) must represent the
	// SpanPropagator's encoding of the core identifying information in `sp`.
	//
	// The second return value (traceAttrs) must represent the SpanPropagator's
	// encoding of trace attributes, per `Span.SetTraceAttribute`.
	PropagateSpanAsText(
		sp Span,
	) (
		contextSnapshot map[string]string,
		traceAttrs map[string]string,
	)

	// JoinTraceFromBinary starts a new Span with the given `operationName`
	// that's joined to the Span that was binary-encoded as `contextSnapshot` and
	// `traceAttrs` (see SpanPropagator.PropagateSpanAsBinary()).
	//
	// If `operationName` is empty, the caller must later call
	// `Span.SetOperationName` on the returned `Span`.
	JoinTraceFromBinary(
		operationName string,
		contextSnapshot []byte,
		traceAttrs []byte,
	) (Span, error)

	// JoinTraceFromBinary starts a new Span with the given `operationName`
	// that's joined to the Span that was text-encoded as `contextSnapshot` and
	// `traceAttrs` (see SpanPropagator.PropagateSpanAsBinary()).
	//
	// If `operationName` is empty, the caller must later call
	// `Span.SetOperationName` on the returned `Span`.
	//
	// It's permissible to pass the same map to both parameters (e.g., an HTTP
	// request headers map): the implementation should only decode the subset
	// of keys it's interested in.
	JoinTraceFromText(
		operationName string,
		contextSnapshot map[string]string,
		traceAttrs map[string]string,
	) (Span, error)
}
