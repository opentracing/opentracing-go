package opentracing

import "sync"

// A TraceContextID is the smallest amount of state needed to describe a span's
// identity within a larger [potentially distributed] trace. The TraceContextID
// is not intended to encode the span's operation name, timing, or log data,
// but merely any unique identifiers (etc) needed to contextualize it within a
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

	// Serializes the TraceContextID as a printable ASCII string (e.g.,
	// base64).
	SerializeASCII() string

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

	// Converts the encoded binary data (see
	// `TraceContextID.SerializeBinary()`) into a TraceContextID of the same
	// type as returned by NewRootTraceContextID().
	DeserializeBinaryTraceContextID(encoded []byte) (TraceContextID, error)

	// Converts the encoded ASCII data (see `TraceContextID.SerializeASCII()`)
	// into a TraceContextID of the same type as returned by
	// NewRootTraceContextID().
	DeserializeASCIITraceContextID(encoded string) (TraceContextID, error)
}

// A `TraceContext` builds off of an implementation-provided TraceContextID and
// adds a simple string map of "trace tags". The trace tags are special in that
// they are propagated *in-band*, presumably alongside application data and the
// `TraceContextID` proper. See the documentation for `SetTraceTag()` for more
// details and some important caveats.
//
// Note that the `TraceContext` is managed internally by the opentracer system;
// opentracer implementations only need to concern themselves with the
// `TraceContextID` (which does not know about trace tags).
type TraceContext struct {
	Id TraceContextID

	tagLock   sync.RWMutex
	traceTags map[string]string
}

// `tags` may be nil.
func newTraceContext(id TraceContextID, tags map[string]string) *TraceContext {
	if tags == nil {
		tags = map[string]string{}
	}
	return &TraceContext{
		Id:        id,
		traceTags: tags,
	}
}

// Set a tag on this TraceContext that also propagates to future TraceContext
// children per `NewChild()`.
//
// `SetTraceTag()` enables powerful functionality given a full-stack
// opentracing integration (e.g., arbitrary application data from a mobile app
// can make it, transparently, all the way into the depths of a storage
// system), and with it some powerful costs: use this feature with care.
//
// IMPORTANT NOTE #1: `SetTraceTag()` will only propagate trace tags to
// *future* children of the `TraceContext` (see `NewChild()`) and/or the `Span`
// that references it.
//
// IMPORTANT NOTE #2: Use this thoughtfully and with care. Every key and value
// is copied into every local *and remote* child of this `TraceContext`, and
// that can add up to a lot of network and cpu overhead for large strings.
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

func (t *TraceContext) SerializeASCII() string {
	// XXX: implement correctly if we like this API
	return t.Id.SerializeASCII()
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
) (*TraceContext, error) {
	// XXX: implement correctly if we like this API
	tcid, err := source.DeserializeBinaryTraceContextID(encoded)
	if err != nil {
		return nil, err
	}
	return newTraceContext(tcid, nil), nil
}

func DeserializeASCIITraceContext(
	source TraceContextIDSource,
	encoded string,
) (*TraceContext, error) {
	// XXX: implement correctly if we like this API
	tcid, err := source.DeserializeASCIITraceContextID(encoded)
	if err != nil {
		return nil, err
	}
	return newTraceContext(tcid, nil), nil
}
