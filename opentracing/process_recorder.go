package opentracing

// A ProcessRecorder handles all of the `RawSpan` data generated via an
// associated `OpenTracer` (see `NewStandardTracer`) instance. It also names
// the containing process and provides access to a straightforward tag map.
type ProcessRecorder interface {
	// Every process in the larger distributed system must have a name. It is
	// up to the `ProcessRecorder` implementation to determine how this name is
	// specified.
	ProcessName() string

	// Adds a tag to help identify or classify the recording process (e.g.,
	// the platform, version/build number, host and/or container name, etc).
	//
	// Returns a reference to this `ProcessRecorder` for chaining, etc.
	SetTag(key string, val interface{}) ProcessRecorder

	// Implementations must determine whether and where to store `span`.
	RecordSpan(span *RawSpan)
}

// NOTE: there should be something like a MultiplexingProcessRecorder which
// itself implements ProcessRecorder but trivially redirects all RecordSpan
// calls to multiple "real" ProcessRecorder implementations.
//
//     type MultiplexingProcessRecorder ... { ... }
