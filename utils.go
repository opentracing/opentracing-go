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
	var traceState, attrs map[string]string
	if err := InjectSpan(sp, WIRE_ENCODING_SPLIT_TEXT, &traceState, &attrs); err != nil {
		return err
	}
	for headerSuffix, val := range traceState {
		h.Add(ContextIDHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
	for headerSuffix, val := range attrs {
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
	extractor := tracer.WireExtractorForEncoding(WIRE_ENCODING_SPLIT_TEXT)
	if extractor == nil {
		return nil, errors.New("No WireExtractor for WIRE_ENCODING_SPLIT_TEXT")
	}

	contextIDMap := make(map[string]string)
	tagsMap := make(map[string]string)
	for key, val := range h {
		if strings.HasPrefix(key, ContextIDHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			contextIDMap[strings.TrimPrefix(key, ContextIDHTTPHeaderPrefix)] = unescaped
		} else if strings.HasPrefix(key, TagsHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			tagsMap[strings.TrimPrefix(key, TagsHTTPHeaderPrefix)] = unescaped
		}
	}
	return extractor.ExtractSpan(operationName, contextIDMap, tagsMap)
}
