package opentracing

type ProcessRecorder interface {
	// Every process in the system must have a name. It is up to the
	// ProcessRecorder implementation to determine how this name is specified.
	ProcessName() string

	// Adds a tag to help identify or classify the recording process (e.g.,
	// the platform, version/build number, host and/or container name, etc).
	SetTag(key string, val interface{}) ProcessRecorder

	RecordSpan(span *RawSpan)
}

// NOTE: there should be something like a MultiplexingProcessRecorder which
// itself implements ProcessRecorder but trivially redirects all RecordSpan
// calls to multiple "real" ProcessRecorder implementations.
//
//     type MultiplexingProcessRecorder ... { ... }
