package opentracing

const (
	OpenTracingContextIDHeader = "OpenTracing-Context-Id"
)

// A ContextID is the smallest amount of state needed to describe a span's
// context within a larger [potentially distributed] trace. The ContextID is
// not intended to encode the span's operation name, timing, or log data, but
// merely any unique identifiers (etc) needed to contextualize it within a
// larger trace tree.
//
// ContextIDs are sufficient to propagate the, well, *context* of a particular
// trace from process to process.
type ContextID interface {
	// Create a child context for this ContextID, and return both that child's
	// own ContextID as well as any Tags that should be added to the child's
	// Span.
	//
	// The returned ContextID type must be the same as the type of the
	// ContextID implementation itself.
	NewChild() (childCtx ContextID, initialChildSpanTags Tags)

	// Serializes the ContextID as an arbitrary byte string.
	Serialize() []byte
}

// A long-lived interface that knows how to create a root ContextID and
// serialize/deserialize any other.
type ContextIDSource interface {
	// Create a ContextID which has no parent (and thus begins its own trace).
	// A ContextIDSource must always return the same type in successive calls
	// to NewRootContextID().
	NewRootContextID() ContextID

	// Converts the encoded binary data (see `SerializeContextID`) into a
	// ContextID of the same type as returned by NewRootContextID().
	DeserializeContextID(encoded []byte) (ContextID, error)
}
