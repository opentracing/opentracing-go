package opentracing

import (
	"strconv"
	"strings"
)

const testHTTPHeaderPrefix = "testprefix-"

// testTracer is a most-noop Tracer implementation that makes it possible for
// unittests to verify whether certain methods were / were not called.
type testTracer struct{}

var fakeIDSource = 1

func nextFakeID() int {
	fakeIDSource++
	return fakeIDSource
}

type testSpan struct {
	OperationName string
	HasParent     bool
	FakeID        int
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
		FakeID:        nextFakeID(),
	}
}

// StartSpanWithOptions belongs to the Tracer interface.
func (n testTracer) StartSpanWithOptions(opts StartSpanOptions) Span {
	return testSpan{
		OperationName: opts.OperationName,
		HasParent:     opts.Parent != nil,
		FakeID:        nextFakeID(),
	}
}

// Inject belongs to the Tracer interface.
func (n testTracer) Inject(sp Span, format interface{}, carrier interface{}) error {
	span := sp.(testSpan)
	switch format {
	case TextMap:
		carrier.(TextMapWriter).Set(testHTTPHeaderPrefix+"fakeid", strconv.Itoa(span.FakeID))
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
		err := carrier.(TextMapReader).ForeachKey(func(key, val string) error {
			switch strings.ToLower(key) {
			case testHTTPHeaderPrefix + "fakeid":
				i, err := strconv.Atoi(val)
				if err != nil {
					return err
				}
				rval.FakeID = i
			}
			return nil
		})
		return rval, err
	}
	return nil, ErrTraceNotFound
}
