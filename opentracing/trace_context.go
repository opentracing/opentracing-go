package opentracing

import "sync"

type TraceContext struct {
	Id TraceContextID

	tagLock   sync.RWMutex
	traceTags map[string]string
}

// Set a tag on this TraceContext that also propagates to future TraceContext
// children per `NewChild()`.
//
// Returns a reference to this TraceContext for chaining, etc.
func (t *TraceContext) SetTraceTag(key, value string) *TraceContext {
	t.tagLock.Lock()
	defer t.tagLock.Unlock()

	t.traceTags[key] = value
	return t
}

// Gets the value for a trace tag given its key. Returns the empty string if
// the value isn't found in this TraceContext.
func (t *TraceContext) TraceTag(key string) string {
	t.tagLock.RLock()
	defer t.tagLock.RUnlock()

	return t.traceTags[key]
}

func (t *TraceContext) NewChild() (childCtx *TraceContext, childSpanTags Tags) {
	t.tagLock.RLock()
	newTags := make(map[string]string, len(t.traceTags))
	for k, v := range t.traceTags {
		newTags[k] = v
	}
	t.tagLock.RUnlock()

	ctxId, childSpanTags := t.Id.NewChild()
	return &TraceContext{
		Id:        ctxId,
		traceTags: newTags,
	}, childSpanTags
}

func (t *TraceContext) SerializeString() string {
	// XXX: implement correctly if we like this API
	return t.Id.SerializeString()
}

func (t *TraceContext) SerializeBinary() []byte {
	// XXX: implement correctly if we like this API
	return t.Id.SerializeBinary()
}

func NewRootTraceContext(source TraceContextIDSource) *TraceContext {
	return &TraceContext{
		Id:        source.NewRootTraceContextID(),
		traceTags: make(map[string]string),
	}
}

func DeserializeBinaryTraceContext(
	source TraceContextIDSource,
	encoded []byte,
) (TraceContext, error) {
	// XXX: implement correctly if we like this API
	return source.DeserializeBinaryTraceContext(encoded)
}
func DeserializeStringTraceContext(
	source TraceContextIDSource,
	encoded string,
) (TraceContext, error) {
	// XXX: implement correctly if we like this API
	return source.DeserializeStringTraceContext(encoded)
}

// A TraceContextID is the smallest amount of state needed to describe a span's
// context within a larger [potentially distributed] trace. The TraceContextID is
// not intended to encode the span's operation name, timing, or log data, but
// merely any unique identifiers (etc) needed to contextualize it within a
// larger trace tree.
//
// TraceContextIDs are sufficient to propagate the, well, *context* of a particular
// trace from process to process.
type TraceContextID interface {
	// Create a child context for this TraceContextID, and return both that child's
	// own TraceContextID as well as any Tags that should be added to the child's
	// Span.
	//
	// The returned TraceContextID type must be the same as the type of the
	// TraceContextID implementation itself.
	NewChild() (childCtx TraceContextID, childSpanTags Tags)

	// Serializes the TraceContextID as a valid unicode string.
	SerializeString() string

	// Serializes the TraceContextID as arbitrary binary data.
	SerializeBinary() []byte
}

// A long-lived interface that knows how to create a root TraceContextID and
// serialize/deserialize any other.
type TraceContextIDSource interface {
	// Create a TraceContextID which has no parent (and thus begins its own trace).
	// A TraceContextIDSource must always return the same type in successive calls
	// to NewRootTraceContextID().
	NewRootTraceContextID() TraceContextID

	// Converts the encoded binary data (see `TraceContextID.Serialize()`) into a
	// TraceContextID of the same type as returned by NewRootTraceContextID().
	DeserializeBinaryTraceContextID(encoded []byte) (TraceContextID, error)
	DeserializeStringTraceContextID(encoded string) (TraceContextID, error)
}
