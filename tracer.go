package opentracing

import "time"

// Tracer is a simple, thin interface for Span creation.
//
// A straightforward implementation is available via the
// `opentracing/basictracer-go` package's `standardtracer.New()'.
type Tracer interface {
	// Create, start, and return a new Span with the given `operationName`, all
	// without specifying a parent Span that can be used to incorporate the
	// newly-returned Span into an existing trace. (I.e., the returned Span is
	// the "root" of its trace).
	//
	// Examples:
	//
	//     var tracer opentracing.Tracer = ...
	//
	//     sp := tracer.StartSpan("GetFeed")
	//
	//     sp := tracer.StartSpanWithOptions(opentracing.SpanOptions{
	//         OperationName: "LoggedHTTPRequest",
	//         Tags: opentracing.Tags{"user_agent", loggedReq.UserAgent},
	//         StartTime: loggedReq.Timestamp,
	//     })
	//
	StartSpan(operationName string, opts ...StartSpanOption) Span

	// Inject() takes the `sc` SpanContext instance and represents it for
	// propagation within `carrier`. The actual type of `carrier` depends on
	// the value of `format`.
	//
	// OpenTracing defines a common set of `format` values (see BuiltinFormat),
	// and each has an expected carrier type.
	//
	// Other packages may declare their own `format` values, much like the keys
	// used by the `net.Context` package (see
	// https://godoc.org/golang.org/x/net/context#WithValue).
	//
	// Example usage (sans error handling):
	//
	//     carrier := opentracing.HTTPHeaderTextMapCarrier(httpReq.Header)
	//     tracer.Inject(
	//         span.SpanContext(),
	//         opentracing.TextMap,
	//         carrier)
	//
	// NOTE: All opentracing.Tracer implementations MUST support all
	// BuiltinFormats.
	//
	// Implementations may return opentracing.ErrUnsupportedFormat if `format`
	// is or not supported by (or not known by) the implementation.
	//
	// Implementations may return opentracing.ErrInvalidCarrier or any other
	// implementation-specific error if the format is supported but injection
	// fails anyway.
	//
	// See Tracer.Join().
	Inject(sc SpanContext, format interface{}, carrier interface{}) error

	// Extract() returns a SpanContext instance given `format` and `carrier`.
	//
	// OpenTracing defines a common set of `format` values (see BuiltinFormat),
	// and each has an expected carrier type.
	//
	// Other packages may declare their own `format` values, much like the keys
	// used by the `net.Context` package (see
	// https://godoc.org/golang.org/x/net/context#WithValue).
	//
	// Example usage (sans error handling):
	//
	//     carrier := opentracing.HTTPHeaderTextMapCarrier(httpReq.Header)
	//     spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
	//
	// NOTE: All opentracing.Tracer implementations MUST support all
	// BuiltinFormats.
	//
	// Return values:
	//  - A successful Extract returns a SpanContext instance and a nil error
	//  - If there was simply no SpanContext to extract in `carrier`, Extract()
	//    returns (nil, opentracing.ErrTraceNotFound)
	//  - If `format` is unsupported or unrecognized, Extract() returns (nil,
	//    opentracing.ErrUnsupportedFormat)
	//  - If there are more fundamental problems with the `carrier` object,
	//    Extract() may return opentracing.ErrInvalidCarrier,
	//    opentracing.ErrTraceCorrupted, or implementation-specific errors.
	//
	// See Tracer.Inject().
	Extract(format interface{}, carrier interface{}) (SpanContext, error)
}

// StartSpanOptions allows Tracer.StartSpanWithOptions callers to override the
// start timestamp, specify a parent Span, and make sure that Tags are
// available at Span initialization time.
type StartSpanOptions struct {
	// OperationName may be empty (and set later via Span.SetOperationName)
	//
	// XXX: get rid of this... StartSpan() requires an opname anyway.
	OperationName string

	// Zero or more causal references to other Spans/SpanContexts. If empty,
	// start a "root" Span (i.e., start a new trace).
	CausalReferences []CausalReference

	// StartTime overrides the Span's start time, or implicitly becomes
	// time.Now() if StartTime.IsZero().
	StartTime time.Time

	// Tags may have zero or more entries; the restrictions on map values are
	// identical to those for Span.SetTag(). May be nil.
	//
	// If specified, the caller hands off ownership of Tags at
	// StartSpan() invocation time.
	Tags map[string]interface{}
}

type CausalReferenceType int

const (
	// RefStartsBefore refers to a span which MUST start before the Span
	// that's starting.
	RefStartsBefore CausalReferenceType = iota

	// RefBlockedOnFinish refers to a span which CAN NOT finish successfully
	// until the Span that's starting has finished.
	RefBlockedOnFinish

	// RefFinishesBefore refers to a span which MUST finish before the Span
	// that's starting.
	RefFinishesBefore

	// RefBlockedParent is the union of RefStartedBefore and RefBlockedOnFinish.
	RefBlockedParent

	// RefRPCClient is the special case of RefBlockedParent for the RPC client
	// peer of an RPC server span.
	RefRPCClient

	// TODO: etc etc, per
	// https://github.com/opentracing/opentracing.github.io/issues/28
)

type CausalReference struct {
	RefType CausalReferenceType
	SpanContext
}

type StartSpanOption func(*StartSpanOptions)

func OperationName(opName string) StartSpanOption {
	return func(opts *StartSpanOptions) {
		opts.OperationName = opName
	}
}

func Reference(t CausalReferenceType, sc SpanContext) StartSpanOption {
	return func(opts *StartSpanOptions) {
		opts.CausalReferences = append(opts.CausalReferences, CausalReference{
			RefType:     t,
			SpanContext: sc,
		})
	}
}

func StartTime(t time.Time) StartSpanOption {
	return func(opts *StartSpanOptions) {
		opts.StartTime = t
	}
}

func StartTags(t map[string]interface{}) StartSpanOption {
	return func(opts *StartSpanOptions) {
		opts.Tags = t
	}
}
