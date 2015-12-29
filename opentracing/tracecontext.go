package opentracing

import (
	"regexp"
	"strings"
)

// Tags are a generic map from an arbitrary string key to an opaque value type.
// The underlying tracing system is responsible for interpreting and
// serializing the values.
type Tags map[string]interface{}

// TraceContext encpasulates the smallest amount of state needed to describe a
// Span's identity within a larger [potentially distributed] trace. The
// TraceContext is not intended to encode the span's operation name, timing,
// or log data, but merely any unique identifiers (etc) needed to contextualize
// it within a larger trace tree.
//
// TraceContexts are sufficient to propagate the, well, *context* of a
// particular trace between processes.
//
// TraceContext also support a simple string map of "trace attributes". These
// trace attributes are special in that they are propagated *in-band*,
// presumably alongside application data. See the documentation for
// SetTraceAttribute() for more details and some important caveats.
type TraceContext interface {
	// SetTraceAttribute sets a tag on this TraceContext that also propagates
	// to future TraceContext children per TraceContext.NewChild.
	//
	// SetTraceAttribute() enables powerful functionality given a full-stack
	// opentracing integration (e.g., arbitrary application data from a mobile
	// app can make it, transparently, all the way into the depths of a storage
	// system), and with it some powerful costs: use this feature with care.
	//
	// IMPORTANT NOTE #1: SetTraceAttribute() will only propagate trace
	// attributes to *future* children of the TraceContext (see NewChild())
	// and/or the Span that references it.
	//
	// IMPORTANT NOTE #2: Use this thoughtfully and with care. Every key and
	// value is copied into every local *and remote* child of this
	// TraceContext, and that can add up to a lot of network and cpu
	// overhead.
	//
	// IMPORTANT NOTE #3: Trace attributes keys have a restricted format:
	// implementations may wish to use them as HTTP header keys (or key
	// suffixes), and of course HTTP headers are case insensitive.
	//
	// As such, `restrictedKey` MUST match the regular expression
	// `(?i:[a-z0-9][-a-z0-9]*)` and is case-insensitive. That is, it must
	// start with a letter or number, and the remaining characters must be
	// letters, numbers, or hyphens. See CanonicalizeTraceAttributeKey(). If
	// `restrictedKey` does not meet these criteria, SetTraceAttribute()
	// results in undefined behavior.
	//
	// Returns a reference to this TraceContext for chaining, etc.
	SetTraceAttribute(restrictedKey, value string) TraceContext

	// Gets the value for a trace tag given its key. Returns the empty string
	// if the value isn't found in this TraceContext.
	//
	// See the `SetTraceAttribute` notes about `restrictedKey`.
	TraceAttribute(restrictedKey string) string
}

// TraceContextMarshaler is a simple interface to marshal a TraceContext to a
// binary byte array or a string-to-string map.
type TraceContextMarshaler interface {
	// Converts the TraceContext into marshaled binary data (see
	// TraceContextUnmarshaler.UnmarshalTraceContextBinary()).
	//
	// The first return value must represent the marshaler's serialization of
	// the core identifying information in `tc`.
	//
	// The second return value must represent the marshaler's serialization of
	// the trace attributes, per `SetTraceAttribute` and `TraceAttribute`.
	MarshalTraceContextBinary(
		tc TraceContext,
	) (
		traceContextID []byte,
		traceAttrs []byte,
	)

	// Converts the TraceContext into a marshaled string:string map (see
	// TraceContextUnmarshaler.UnmarshalTraceContextStringMap()).
	//
	// The first return value must represent the marshaler's serialization of
	// the core identifying information in `tc`.
	//
	// The second return value must represent the marshaler's serialization of
	// the trace attributes, per `SetTraceAttribute` and `TraceAttribute`.
	MarshalTraceContextStringMap(
		tc TraceContext,
	) (
		traceContextID map[string]string,
		traceAttrs map[string]string,
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
	// attributes (per `SetTraceAttribute` and `TraceAttribute`) attached to a
	// TraceContext instance.
	UnmarshalTraceContextBinary(
		traceContextID []byte,
		traceAttrs []byte,
	) (TraceContext, error)

	// Converts the marshaled string:string map (see
	// TraceContextMarshaler.MarshalTraceContextStringMap()) into a TraceContext.
	//
	// The first parameter contains the marshaler's serialization of the core
	// identifying information in a TraceContext instance.
	//
	// The second parameter contains the marshaler's serialization of the trace
	// attributes (per `SetTraceAttribute` and `TraceAttribute`) attached to a
	// TraceContext instance.
	//
	// It's permissible to pass the same map to both parameters (e.g., an HTTP
	// request headers map): the implementation should only unmarshal the
	// subset its interested in.
	UnmarshalTraceContextStringMap(
		traceContextID map[string]string,
		traceAttrs map[string]string,
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

	// NewChildTraceContext creates a child context for `parent`, and returns
	// both that child's own TraceContext as well as any Tags that should be
	// added to the child's Span.
	//
	// The returned TraceContext type must be the same as the type of the
	// TraceContext implementation itself.
	NewChildTraceContext(parent TraceContext) (childCtx TraceContext, childSpanTags Tags)
}

var regexTraceAttribute = regexp.MustCompile("^(?i:[a-z0-9][-a-z0-9]*)$")

// CanonicalizeTraceAttributeKey returns the canonicalized version of trace tag
// key `key`, and true if and only if the key was valid.
func CanonicalizeTraceAttributeKey(key string) (string, bool) {
	if !regexTraceAttribute.MatchString(key) {
		return "", false
	}
	return strings.ToLower(key), true
}

// Merge incorporates the keys and values from `other` into this `Tags`
// instance, then returns same.
func (t Tags) Merge(other Tags) Tags {
	for k, v := range other {
		t[k] = v
	}
	return t
}
