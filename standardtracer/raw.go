package standardtracer

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

// RawSpan encapsulates all state associated with a (finished) Span.
type RawSpan struct {
	opentracing.TraceContext

	// The name of the "operation" this span is an instance of. (Called a "span
	// name" in some implementations)
	Operation string

	// We store <start, duration> rather than <start, end> so that only
	// one of the timestamps has global clock uncertainty issues.
	Start    time.Time
	Duration time.Duration

	// Essentially an extension mechanism. Can be used for many purposes,
	// not to be enumerated here.
	Tags opentracing.Tags

	// The span's "microlog".
	Logs []*RawLog
}

// RawLog encapsolutes all state associated with a log element in a Span.
type RawLog struct {
	Timestamp time.Time

	// Self-explanatory :)
	Error bool

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
