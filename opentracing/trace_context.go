package opentracing

// A TraceContextID is the smallest amount of state needed to describe a span's
// identity within a larger [potentially distributed] trace. The TraceContextID
// is not intended to encode the span's operation name, timing, or log data,
// but merely any unique identifiers (etc) needed to contextualize it within a
// larger trace tree.
//
// TraceContextIDs are sufficient to propagate the, well, *context* of a particular
// trace from process to process.
//
// XXX incorporate:
//
// A `TraceContext` builds off of an implementation-provided TraceContextID and
// adds a simple string map of "trace tags". The trace tags are special in that
// they are propagated *in-band*, presumably alongside application data and the
// `TraceContextID` proper. See the documentation for `SetTraceTag()` for more
// details and some important caveats.
//
// Note that the `TraceContext` is managed internally by the opentracer system;
// opentracer implementations only need to concern themselves with the
// `TraceContextID` (which does not know about trace tags).
type TraceContext interface {
	// Create a child context for this TraceContextID, and return both that child's
	// own TraceContextID as well as any Tags that should be added to the child's
	// Span.
	//
	// The returned TraceContextID type must be the same as the type of the
	// TraceContextID implementation itself.
	NewChild() (childCtx TraceContext, childSpanTags Tags)

	// Set a tag on this TraceContext that also propagates to future
	// TraceContext children per `NewChild()`.
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
	// `TraceContext`, and that can add up to a lot of network and cpu overhead
	// for large strings.
	//
	// Returns a reference to this TraceContext for chaining, etc.
	SetTraceTag(key, value string) TraceContext

	// Gets the value for a trace tag given its key. Returns the empty string
	// if the value isn't found in this TraceContext.
	TraceTag(key string) string
}

// XXX: comment
type TraceContextMarshaler interface {
	MarshalBinaryTraceContext(tcid TraceContext) []byte
	MarshalStringMapTraceContext(tcid TraceContext) map[string]string
}

// XXX: comment
type TraceContextUnmarshaler interface {
	// Converts the marshaled binary data (see
	// `TraceContextMarshaler.MarshalBinary()`) into a TraceContext.
	UnmarshalBinaryTraceContext(marshaled []byte) (TraceContext, error)
	// Converts the marshaled string:string map (see
	// `TraceContextMarshaler.MarshalStringMap()`) into a TraceContext.
	UnmarshalStringMapTraceContext(marshaled map[string]string) (TraceContext, error)
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
