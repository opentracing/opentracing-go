package opentracing

import (
	"regexp"
	"strings"
	"time"
)

// Span represents an active, un-finished span in the opentracing system.
//
// Spans are created by the Tracer interface and Span.StartChild.
type Span interface {
	// Sets or changes the operation name.
	SetOperationName(operationName string) Span

	// Creates and starts a child span.
	StartChild(operationName string) Span

	// Adds a tag to the span.
	//
	// Tag values can be of arbitrary types, however the treatment of complex
	// types is dependent on the underlying tracing system implementation.
	// It is expected that most tracing systems will handle primitive types
	// like strings and numbers. If a tracing system cannot understand how
	// to handle a particular value type, it may ignore the tag, but shall
	// not panic.
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	SetTag(key string, value interface{}) Span

	// Sets the end timestamp and calls the `Recorder`s RecordSpan()
	// internally.
	//
	// Finish() should be the last call made to any span instance, and to do
	// otherwise leads to undefined behavior.
	Finish()

	// LogEvent() is equivalent to
	//
	//   Log(LogData{Event: event})
	//
	LogEvent(event string)

	// LogEventWithPayload() is equivalent to
	//
	//   Log(LogData{Event: event, Payload: payload0})
	//
	LogEventWithPayload(event string, payload interface{})

	// Log() records `data` to this Span.
	//
	// See LogData for semantic details.
	Log(data LogData)

	// SetTraceAttribute sets a key:value pair on this Span that also
	// propagates to future Span children.
	//
	// SetTraceAttribute() enables powerful functionality given a full-stack
	// opentracing integration (e.g., arbitrary application data from a mobile
	// app can make it, transparently, all the way into the depths of a storage
	// system), and with it some powerful costs: use this feature with care.
	//
	// IMPORTANT NOTE #1: SetTraceAttribute() will only propagate trace
	// attributes to *future* children of the Span.
	//
	// IMPORTANT NOTE #2: Use this thoughtfully and with care. Every key and
	// value is copied into every local *and remote* child of this Span, and
	// that can add up to a lot of network and cpu overhead.
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
	// Returns a reference to this Span for chaining, etc.
	SetTraceAttribute(restrictedKey, value string) Span

	// Gets the value for a trace tag given its key. Returns the empty string
	// if the value isn't found in this Span.
	//
	// See the `SetTraceAttribute` notes about `restrictedKey`.
	TraceAttribute(restrictedKey string) string
}

// LogData is data associated to a Span. Every LogData instance should specify
// at least one of Event and/or Payload.
type LogData struct {
	// The timestamp of the log record; if set to the default value (the unix
	// epoch), implementations should use time.Now() implicitly.
	Timestamp time.Time

	// Event (if non-empty) should be the stable name of some notable moment in
	// the lifetime of a Span. For instance, a Span representing a browser page
	// load might add an Event for each of the Performance.timing moments
	// here: https://developer.mozilla.org/en-US/docs/Web/API/PerformanceTiming
	//
	// While it is not a formal requirement, Event strings will be most useful
	// if they are *not* unique; rather, tracing systems should be able to use
	// them to understand how two similar Spans relate from an internal timing
	// perspective.
	Event string

	// Payload is a free-form potentially structured object which Tracer
	// implementations may retain and record all, none, or part of.
	//
	// If included, `Payload` should be restricted to data derived from the
	// instrumented application; in particular, it should not be used to pass
	// semantic flags to a Log() implementation.
	//
	// For example, an RPC system could log the wire contents in both
	// directions, or a SQL library could log the query (with or without
	// parameter bindings); tracing implementations may truncate or otherwise
	// record only a snippet of these payloads (or may strip out PII, etc,
	// etc).
	Payload interface{}
}

// Tags are a generic map from an arbitrary string key to an opaque value type.
// The underlying tracing system is responsible for interpreting and
// serializing the values.
type Tags map[string]interface{}

// Merge incorporates the keys and values from `other` into this `Tags`
// instance, then returns same.
func (t Tags) Merge(other Tags) Tags {
	for k, v := range other {
		t[k] = v
	}
	return t
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
