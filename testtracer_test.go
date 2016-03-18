package opentracing

import "strconv"

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
		carrier.(TextMapWriter).Add("opname", span.OperationName)
		carrier.(TextMapWriter).Add("hasparent", strconv.FormatBool(span.HasParent))
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
		err := carrier.(TextMapReader).ReadAllEntries(func(key, val string) error {
			switch key {
			case "hasparent":
				b, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				rval.HasParent = b
			case "opname":
				rval.OperationName = val
			}
			return nil
		})
		return rval, err
	}
	return nil, ErrTraceNotFound
}
