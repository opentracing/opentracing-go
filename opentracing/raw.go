package opentracing

import "time"

type Tags map[string]interface{}

type RawSpan struct {
	ContextID

	// The name of the "operation" this span is an instance of. (Called a "span
	// name" in some implementations)
	Operation string

	// We store <start, duration> rather than <start, end> so that only
	// one of the timestamps has global clock uncertainty issues.
	Start    time.Time
	Duration time.Duration

	// Essentially an extension mechanism. Can be used for many purposes,
	// not to be enumerated here.
	Tags Tags

	// The span's "microlog".
	Logs []*RawLog
}

type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
)

type RawLog struct {
	Timestamp time.Time

	// Info/Warning/Error.
	Severity

	// `Message` is a format string and can refer to fields in the payload by path, like so:
	//
	//   "first transaction is worth ${transactions[0].amount} ${transactions[0].currency}"
	//
	// , and the payload might look something like
	//
	//   {transactions: [{amount: 10, currency: "USD"}, {amount: 11, currency: "USD"}]}
	Message string

	// `Payload` can be a POD type, a string, or nested maps and slices; i.e.,
	// it is a base type or an anonymous struct.
	Payload interface{}
}

type Recorder interface {
	// Adds a tag to help identify or classify the recording process (e.g.,
	// the platform, version/build number, host and/or container name, etc).
	SetTag(key string, val interface{})

	RecordSpan(span *RawSpan)
}

// NOTE: there should be something like a MultiplexingRecorder which itself
// implements Recorder but trivially redirects all RecordSpan calls to
// multiple "real" Recorder implementations.
//
//     type MultiplexingRecorder ... { ... }
