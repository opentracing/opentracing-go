package opentracing

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	// ContextIDHTTPHeaderPrefix precedes the opentracing-related ContextID HTTP
	// headers.
	ContextIDHTTPHeaderPrefix = "Open-Tracing-Context-Id-"

	// TagsHTTPHeaderPrefix precedes the opentracing-related trace-tags HTTP
	// headers.
	TagsHTTPHeaderPrefix = "Open-Tracing-Trace-Tags-"
)

// InjectSpanInHeader encodes Span `sp` in `h` as a series of HTTP headers.
// Values are URL-escaped.
func InjectSpanInHeader(
	sp Span,
	h http.Header,
) error {
	// First, look for a PROPAGATION_FORMAT_GO_HTTP_HEADER injector (our preference).
	injector := sp.PropagationInjectorForFormat(PROPAGATION_FORMAT_GO_HTTP_HEADER)
	if injector != nil {
		return injector.InjectSpan(sp, h)
	}

	// Else, fall back on PROPAGATION_FORMAT_SPLIT_TEXT.
	if injector = sp.PropagationInjectorForFormat(PROPAGATION_FORMAT_SPLIT_TEXT); injector == nil {
		return errors.New("No suitable injector")
	}
	carrier := NewTextCarrier()
	inject.InjectSpan(sp, carrier)
	for headerSuffix, val := range carrier.TracerState {
		h.Add(ContextIDHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
	for headerSuffix, val := range carrier.TraceAttributes {
		h.Add(TagsHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
}

// JoinTraceFromHeader decodes a Span with operation name `operationName` from
// `h`, expecting that header values are URL-escpaed.
//
// If `operationName` is empty, the caller must later call
// `Span.SetOperationName` on the returned `Span`.
func JoinTraceFromHeader(
	operationName string,
	h http.Header,
	tracer Tracer,
) (Span, error) {
	// First, look for a PROPAGATION_FORMAT_GO_HTTP_HEADER extractor (our
	// preference).
	extractor := tracer.PropagationExtractorForFormat(PROPAGATION_FORMAT_GO_HTTP_HEADER)
	if extractor != nil {
		return extractor.ExtractSpan(operationName, h)
	}

	// Else, fall back on PROPAGATION_FORMAT_SPLIT_TEXT.
	if extractor = tracer.PropagationExtractorForFormat(PROPAGATION_FORMAT_SPLIT_TEXT); extractor == nil {
		return nil, errors.New("No suitable extractor")
	}

	carrier := NewTextCarrier()
	for key, val := range h {
		if strings.HasPrefix(key, ContextIDHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			carrier.TracerState[strings.TrimPrefix(key, ContextIDHTTPHeaderPrefix)] = unescaped
		} else if strings.HasPrefix(key, TagsHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			carrier.TraceAttributes[strings.TrimPrefix(key, TagsHTTPHeaderPrefix)] = unescaped
		}
	}
	return extractor.ExtractSpan(operationName, carrier)
}
