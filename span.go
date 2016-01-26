package opentracing

import (
	"time"

	"golang.org/x/net/context"
)

// Span represents an active, un-finished span in the opentracing system.
//
// Spans are created by the Tracer interface and Span.StartChild.
type Span interface {
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

	// Suitable for serializing over the wire, etc.
	TraceContext() TraceContext

	// LogEvent() is equivalent to
	//
	//   Log(time.Now(), LogData{Event: event})
	//
	LogEvent(event string)

	// LogEventWithPayload() is equivalent to
	//
	//   Log(time.Now(), LogData{Event: event, Payload: payload0})
	//
	LogEventWithPayload(event string, payload interface{})

	// Log() records `data` to this Span.
	//
	// See LogData for semantic details.
	Log(data LogData)

	// A convenience method. Equivalent to
	//
	//    var goCtx context.Context = ...
	//    var span Span = ...
	//    goCtx := opentracing.GoContextWithSpan(ctx, span)
	//
	//
	// NOTE: We use the term "GoContext" to minimize confusion with
	// TraceContext.
	AddToGoContext(goCtx context.Context) (Span, context.Context)
}

// See Span.Log(). Every LogData instance should specify at least one of Event
// and/or Payload.
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
