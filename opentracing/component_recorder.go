package opentracing

type ComponentRecorder interface {
	// Every system component must have a name. It is up to the
	// ComponentRecorder implementation to determine how this name is
	// specified.
	ComponentName() string

	// Adds a tag to help identify or classify the recording process (e.g.,
	// the platform, version/build number, host and/or container name, etc).
	SetTag(key string, val interface{})

	RecordSpan(span *RawSpan)
}

// NOTE: there should be something like a MultiplexingComponentRecorder which
// itself implements ComponentRecorder but trivially redirects all RecordSpan
// calls to multiple "real" ComponentRecorder implementations.
//
//     type MultiplexingComponentRecorder ... { ... }
