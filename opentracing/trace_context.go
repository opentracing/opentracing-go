package opentracing

// A `TraceContext` is the smallest amount of state needed to describe a span's
// identity within a larger [potentially distributed] trace. The `TraceContext`
// is not intended to encode the span's operation name, timing, or log data,
// but merely any unique identifiers (etc) needed to contextualize it within a
// larger trace tree.
//
// `TraceContext`s are sufficient to propagate the, well, *context* of a
// particular trace between processes.
//
// `TraceContext` also support a simple string map of "trace tags". These trace
// tags are special in that they are propagated *in-band*, presumably alongside
// application data. See the documentation for `SetTraceTag()` for more details
// and some important caveats.
type TraceContext interface {
	// Create a child context for this `TraceContext`, and return both that
	// child's own `TraceContext` as well as any Tags that should be added to
	// the child's Span.
	//
	// The returned `TraceContext` type must be the same as the type of the
	// `TraceContext` implementation itself.
	NewChild() (childCtx TraceContext, childSpanTags Tags)

	// Set a tag on this `TraceContext` that also propagates to future
	// `TraceContext` children per `NewChild()`.
	//
	// `SetTraceTag()` enables powerful functionality given a full-stack
	// opentracing integration (e.g., arbitrary application data from a mobile
	// app can make it, transparently, all the way into the depths of a storage
	// system), and with it some powerful costs: use this feature with care.
	//
	// IMPORTANT NOTE #1: `SetTraceTag()` will only propagate trace tags to
	// *future* children of the `TraceContext` (see `NewChild()`) and/or the
	// `Span` that references it.
	//
	// IMPORTANT NOTE #2: Use this thoughtfully and with care. Every key and
	// value is copied into every local *and remote* child of this
	// `TraceContext`, and that can add up to a lot of network and cpu
	// overhead.
	//
	// Returns a reference to this `TraceContext` for chaining, etc.
	SetTraceTag(key, value string) TraceContext

	// Gets the value for a trace tag given its key. Returns the empty string
	// if the value isn't found in this `TraceContext`.
	TraceTag(key string) string
}

// A simple interface to marshal a `TraceContext` to a binary byte array or a
// string-to-string map.
type TraceContextMarshaler interface {
	// Converts the `TraceContext` into marshaled binary data (see
	// `TraceContextUnmarshaler.UnmarshalTraceContextBinary()`).
	MarshalTraceContextBinary(tcid TraceContext) []byte
	// Converts the `TraceContext` into a marshaled string:string map (see
	// `TraceContextUnmarshaler.UnmarshalTraceContextStringMap()`).
	MarshalTraceContextStringMap(tcid TraceContext) map[string]string
}

// A simple interface to marshal a `TraceContext` to a binary byte array or a
// string-to-string map.
type TraceContextUnmarshaler interface {
	// Converts the marshaled binary data (see
	// `TraceContextMarshaler.MarshalTraceContextBinary()`) into a TraceContext.
	UnmarshalTraceContextBinary(marshaled []byte) (TraceContext, error)
	// Converts the marshaled string:string map (see
	// `TraceContextMarshaler.MarshalTraceContextStringMap()`) into a TraceContext.
	UnmarshalTraceContextStringMap(marshaled map[string]string) (TraceContext, error)
}

// A long-lived interface that knows how to create a root TraceContext and
// serialize/deserialize any other.
type TraceContextSource interface {
	TraceContextMarshaler
	TraceContextUnmarshaler

	// Create a TraceContext which has no parent (and thus begins its own trace).
	// A TraceContextSource must always return the same type in successive calls
	// to NewRootTraceContext().
	NewRootTraceContext() TraceContext
}
