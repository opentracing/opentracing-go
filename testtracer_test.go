package opentracing

import (
	"strconv"
	"strings"
)

const testHTTPHeaderPrefix = "testprefix-"

// testTracer is a most-noop Tracer implementation that makes it possible for
// unittests to verify whether certain methods were / were not called.
type testTracer struct{}

type testSpan struct {
	OperationName string
	HasParent     bool
}

// testSpan:
func (n testSpan) SetTag(key string, value interface{}) Span             { return n }
func (n testSpan) Finish()                                               {}
func (n testSpan) FinishWithOptions(opts FinishOptions)                  {}
func (n testSpan) SetBaggageItem(key, val string) Span                   { return n }
func (n testSpan) BaggageItem(key string) string                         { return "" }
func (n testSpan) LogEvent(event string)                                 {}
func (n testSpan) LogEventWithPayload(event string, payload interface{}) {}
func (n testSpan) Log(data LogData)                                      {}
func (n testSpan) SetOperationName(operationName string) Span            { return n }
func (n testSpan) Tracer() Tracer                                        { return testTracer{} }

// StartSpan belongs to the Tracer interface.
func (n testTracer) StartSpan(operationName string) Span {
	return testSpan{
		OperationName: operationName,
		HasParent:     false,
	}
}

// StartSpanWithOptions belongs to the Tracer interface.
func (n testTracer) StartSpanWithOptions(opts StartSpanOptions) Span {
	return testSpan{
		OperationName: opts.OperationName,
		HasParent:     opts.Parent != nil,
	}
}

// Inject belongs to the Tracer interface.
func (n testTracer) Inject(sp Span, format interface{}, carrier interface{}) error {
	span := sp.(testSpan)
	switch format {
	case TextMap:
		// Just for testing purposes... generally not a worthwhile thing to
		// propagate.
		carrier.(TextMapWriter).Add(testHTTPHeaderPrefix+"opname", span.OperationName)
		carrier.(TextMapWriter).Add(testHTTPHeaderPrefix+"hasparent", strconv.FormatBool(span.HasParent))
		return nil
	}
	return ErrUnsupportedFormat
}

// Join belongs to the Tracer interface.
func (n testTracer) Join(operationName string, format interface{}, carrier interface{}) (Span, error) {
	switch format {
	case TextMap:
		// Just for testing purposes... generally not a worthwhile thing to
		// propagate.
		rval := testSpan{}
		err := carrier.(TextMapReader).ReadAll(func(key, val string) error {
			switch strings.ToLower(key) {
			case testHTTPHeaderPrefix + "hasparent":
				b, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				rval.HasParent = b
			case testHTTPHeaderPrefix + "opname":
				rval.OperationName = val
			}
			return nil
		})
		return rval, err
	}
	return nil, ErrTraceNotFound
}
