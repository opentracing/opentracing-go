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

	// Log() records `data` to this Span at `timestamp`.
	//
	// See LogData for semantic details.
	Log(timestamp time.Time, data LogData)

	// Info() is equivalent to
	//
	//   Log(time.Now(), LogData{Message: fmt.Sprint(args...)})
	//
	Info(args ...interface{})

	// Infof() is equivalent to
	//
	//   Log(time.Now(), LogData{Message: fmt.Sprintf(format, args...)})
	//
	Infof(format string, args ...interface{})

	// Error() is equivalent to
	//
	//   Log(time.Now(), LogData{Message: fmt.Sprint(args...), Severity: ERROR})
	//
	Error(args ...interface{})

	// Errorf() is equivalent to
	//
	//   Log(time.Now(), LogData{Message: fmt.Sprintf(format, args...), Severity: ERROR})
	//
	Errorf(format string, args ...interface{})

	// Event() is equivalent to
	//
	//   Log(time.Now(), LogData{Event: event})
	//
	// if payload is unspecified; otherwise, only one payload argument is
	// accepted, and Event() is equivalent to
	//
	//   Log(time.Now(), LogData{Event: event, Payload: payload[0]})
	//
	Event(event string, payload ...interface{})

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

// See Span.Log(). Every LogData instance should specify at least one of
// Message, Event, or Payload.
type LogData struct {
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

	// Message (if non-empty) is a free-form debugging string, much in keeping
	// with general practices around console-oriented logging of human-readable
	// messages.
	Message string

	// Severity will most often be non-V0 when `Message` is non-empty, but that
	// is not a formal requirement.
	Severity Severity

	// Payload is a free-form potentially structured object which Tracer
	// implementations may retain and record all, none, or part of.
	//
	// If included, `Payload` should be restricted to data derived from the
	// instrumented application; in particular, it should not be used to pass
	// semantic flags to a Log() implementation.
	Payload interface{}
}

type Severity int

const (
	// Loosely modeled after https://github.com/google/glog/blob/v0.3.4/src/glog/log_severity.h#L47
	INFO    = 0
	WARNING = 1
	ERROR   = 2
)
