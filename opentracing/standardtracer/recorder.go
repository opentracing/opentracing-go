package standardtracer

// ProcessIdentifier is a thin interface that guarantees all implementors
// represent a ProcessName and accepts arbitrary process-level tag assignment
// (e.g., build numbers, platforms, hostnames, etc).
type ProcessIdentifier interface {
	// Every process in the opentracing-instrumented distributed system must
	// have a name.
	ProcessName() string

	// Adds a tag to help identify or classify the recording process (e.g.,
	// the platform, version/build number, host and/or container name, etc).
	//
	// Returns a reference to this `ProcessIdentifier` for chaining, etc.
	SetTag(key string, val interface{}) ProcessIdentifier
}

// A Recorder handles all of the `RawSpan` data generated via an
// associated `Tracer` (see `NewStandardTracer`) instance. It also names
// the containing process and provides access to a straightforward tag map.
type Recorder interface {
	ProcessIdentifier

	// Implementations must determine whether and where to store `span`.
	RecordSpan(span *RawSpan)
}

// NOTE: there should be something like a MultiplexingRecorder which
// itself implements Recorder but trivially redirects all RecordSpan
// calls to multiple "real" Recorder implementations.
//
//     type MultiplexingRecorder ... { ... }
