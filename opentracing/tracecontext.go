package opentracing

// TraceContext encpasulates the smallest amount of state needed to describe a
// Span's identity within a larger [potentially distributed] trace. The
// TraceContext is not intended to encode the span's operation name, timing,
// or log data, but merely any unique identifiers (etc) needed to contextualize
// it within a larger trace tree.
//
// TraceContexts are sufficient to propagate the, well, *context* of a
// particular trace between processes.
//
// TraceContext also support a simple string map of "trace tags". These trace
// tags are special in that they are propagated *in-band*, presumably alongside
// application data. See the documentation for SetTraceTag() for more details
// and some important caveats.
type TraceContext interface {
	// NewChild creates a child context for this TraceContext, and returns both
	// that child's own TraceContext as well as any Tags that should be added
	// to the child's Span.
	//
	// The returned TraceContext type must be the same as the type of the
	// TraceContext implementation itself.
	NewChild() (childCtx TraceContext, childSpanTags Tags)

	// SetTraceTag sets a tag on this TraceContext that also propagates to
	// future TraceContext children per TraceContext.NewChild.
	//
	// SetTraceTag() enables powerful functionality given a full-stack
	// opentracing integration (e.g., arbitrary application data from a mobile
	// app can make it, transparently, all the way into the depths of a storage
	// system), and with it some powerful costs: use this feature with care.
	//
	// IMPORTANT NOTE #1: SetTraceTag() will only propagate trace tags to
	// *future* children of the TraceContext (see NewChild()) and/or the
	// Span that references it.
	//
	// IMPORTANT NOTE #2: Use this thoughtfully and with care. Every key and
	// value is copied into every local *and remote* child of this
	// TraceContext, and that can add up to a lot of network and cpu
	// overhead.
	//
	// IMPORTANT NOTE #3: Trace tags are case-insensitive: implementations may
	// wish to use them as HTTP header keys (or key suffixes), and of course
	// HTTP headers are case insensitive.
	//
	// Returns a reference to this TraceContext for chaining, etc.
	SetTraceTag(caseInsensitiveKey, value string) TraceContext

	// Gets the value for a trace tag given its key. Returns the empty string
	// if the value isn't found in this TraceContext.
	TraceTag(caseInsensitiveKey string) string
}

// TraceContextMarshaler is a simple interface to marshal a TraceContext to a
// binary byte array or a string-to-string map.
type TraceContextMarshaler interface {
	// Converts the TraceContext into marshaled binary data (see
	// TraceContextUnmarshaler.UnmarshalTraceContextBinary()).
	//
	// The first return value must represent the marshaler's serialization of
	// the core identifying information in `tcid`.
	//
	// The second return value must represent the marshaler's serialization of
	// the trace tags, per `SetTraceTag` and `TraceTag`.
	MarshalTraceContextBinary(
		tcid TraceContext,
	) (
		traceContextID []byte,
		traceTags []byte,
	)

	// Converts the TraceContext into a marshaled string:string map (see
	// TraceContextUnmarshaler.UnmarshalTraceContextStringMap()).
	//
	// The first return value must represent the marshaler's serialization of
	// the core identifying information in `tcid`.
	//
	// The second return value must represent the marshaler's serialization of
	// the trace tags, per `SetTraceTag` and `TraceTag`.
	MarshalTraceContextStringMap(
		tcid TraceContext,
	) (
		traceContextID map[string]string,
		traceTags map[string]string,
	)
}

// TraceContextUnmarshaler is a simple interface to unmarshal a binary byte
// array or a string-to-string map into a TraceContext.
type TraceContextUnmarshaler interface {
	// Converts the marshaled binary data (see
	// TraceContextMarshaler.MarshalTraceContextBinary()) into a TraceContext.
	//
	// The first parameter contains the marshaler's serialization of the core
	// identifying information in a TraceContext instance.
	//
	// The second parameter contains the marshaler's serialization of the trace
	// tags (per `SetTraceTag` and `TraceTag`) attached to a TraceContext
	// instance.
	UnmarshalTraceContextBinary(
		traceContextID []byte,
		traceTags []byte,
	) (TraceContext, error)

	// Converts the marshaled string:string map (see
	// TraceContextMarshaler.MarshalTraceContextStringMap()) into a TraceContext.
	//
	// The first parameter contains the marshaler's serialization of the core
	// identifying information in a TraceContext instance.
	//
	// The second parameter contains the marshaler's serialization of the trace
	// tags (per `SetTraceTag` and `TraceTag`) attached to a TraceContext
	// instance.
	//
	// It's permissable to pass the same map to both parameters (e.g., an HTTP
	// request headers map): the implementation should only unmarshal the
	// subset its interested in.
	UnmarshalTraceContextStringMap(
		traceContextID map[string]string,
		traceTags map[string]string,
	) (TraceContext, error)
}

// TraceContextSource is a long-lived interface that knows how to create a root
// TraceContext and marshal/unmarshal any other.
type TraceContextSource interface {
	TraceContextMarshaler
	TraceContextUnmarshaler

	// Create a TraceContext which has no parent (and thus begins its own trace).
	// A TraceContextSource must always return the same type in successive calls
	// to NewRootTraceContext().
	NewRootTraceContext() TraceContext
}
