package mocktracer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opentracing/opentracing-go"
)

const mockTextMapIdsPrefix = "mockpfx-ids-"
const mockTextMapBaggagePrefix = "mockpfx-baggage-"

// Injector is responsible for injecting SpanContext instances in a manner suitable
// for propagation via a format-specific "carrier" object. Typically the
// injection will take place across an RPC boundary, but message queues and
// other IPC mechanisms are also reasonable places to use an Injector.
type Injector interface {
	// Inject takes `SpanContext` and injects it into `carrier`. The actual type
	// of `carrier` depends on the `format` passed to `Tracer.Inject()`.
	//
	// Implementations may return opentracing.ErrInvalidCarrier or any other
	// implementation-specific error if injection fails.
	Inject(ctx *MockSpanContext, carrier interface{}) error
}

// Extractor is responsible for extracting SpanContext instances from a
// format-specific "carrier" object. Typically the extraction will take place
// on the server side of an RPC boundary, but message queues and other IPC
// mechanisms are also reasonable places to use an Extractor.
type Extractor interface {
	// Extract decodes a SpanContext instance from the given `carrier`,
	// or (nil, opentracing.ErrSpanContextNotFound) if no context could
	// be found in the `carrier`.
	Extract(carrier interface{}) (*MockSpanContext, error)
}

// TextMapPropagator implements Injector/Extractor for TextMap format.
type TextMapPropagator struct{}

// Inject implements the Injector interface
func (t *TextMapPropagator) Inject(spanContext *MockSpanContext, carrier interface{}) error {
	spanContext.RLock()
	defer spanContext.RUnlock()
	writer, ok := carrier.(opentracing.TextMapWriter)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}
	// Ids:
	writer.Set(mockTextMapIdsPrefix+"traceid", strconv.Itoa(spanContext.TraceID))
	writer.Set(mockTextMapIdsPrefix+"spanid", strconv.Itoa(spanContext.SpanID))
	writer.Set(mockTextMapIdsPrefix+"sampled", fmt.Sprint(spanContext.Sampled))
	// Baggage:
	for baggageKey, baggageVal := range spanContext.baggage {
		writer.Set(mockTextMapBaggagePrefix+baggageKey, baggageVal)
	}
	return nil
}

// Extract implements the Extractor interface
func (t *TextMapPropagator) Extract(carrier interface{}) (*MockSpanContext, error) {
	reader, ok := carrier.(opentracing.TextMapReader)
	if !ok {
		return nil, opentracing.ErrInvalidCarrier
	}
	rval := newMockSpanContext(0, 0, true, nil)
	err := reader.ForeachKey(func(key, val string) error {
		lowerKey := strings.ToLower(key)
		switch {
		case lowerKey == mockTextMapIdsPrefix+"traceid":
			// Ids:
			i, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			rval.TraceID = i
		case lowerKey == mockTextMapIdsPrefix+"spanid":
			// Ids:
			i, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			rval.SpanID = i
		case lowerKey == mockTextMapIdsPrefix+"sampled":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			rval.Sampled = b
		case strings.HasPrefix(lowerKey, mockTextMapBaggagePrefix):
			// Baggage:
			rval.SetBaggageItem(lowerKey[len(mockTextMapBaggagePrefix):], val)
		}
		return nil
	})
	if rval.TraceID == 0 || rval.SpanID == 0 {
		return nil, opentracing.ErrSpanContextNotFound
	}
	if err != nil {
		return nil, err
	}
	return rval, nil
}
