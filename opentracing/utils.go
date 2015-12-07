package opentracing

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	// OpenTracingContextIDHTTPHeaderPrefix precedes the opentracing-related
	// ContextID HTTP headers.
	OpenTracingContextIDHTTPHeaderPrefix = "Open-Tracing-Context-Id-"
	// OpenTracingTagsHTTPHeaderPrefix precedes the opentracing-related
	// trace-tags HTTP headers.
	OpenTracingTagsHTTPHeaderPrefix = "Open-Tracing-Trace-Tags-"
)

// AddTraceContextToHeader marshals TraceContext `ctx` to `h` as a series of
// HTTP headers. Values are URL-escaped.
func AddTraceContextToHeader(
	ctx TraceContext,
	h http.Header,
	marshaler TraceContextMarshaler,
) {
	contextIDMap, tagsMap := marshaler.MarshalTraceContextStringMap(ctx)
	for headerSuffix, val := range contextIDMap {
		h.Add(OpenTracingContextIDHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
	for headerSuffix, val := range tagsMap {
		h.Add(OpenTracingTagsHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
}

// TraceContextFromHeader unmarshals a TraceContext from `h`, expecting that
// header values are URL-escpaed.
func TraceContextFromHeader(
	h http.Header,
	unmarshaler TraceContextUnmarshaler,
) (TraceContext, error) {
	contextIDMap := make(map[string]string)
	tagsMap := make(map[string]string)
	for key, val := range h {
		if strings.HasPrefix(key, OpenTracingContextIDHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			contextIDMap[strings.TrimPrefix(key, OpenTracingContextIDHTTPHeaderPrefix)] = unescaped
		} else if strings.HasPrefix(key, OpenTracingTagsHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			tagsMap[strings.TrimPrefix(key, OpenTracingTagsHTTPHeaderPrefix)] = unescaped
		}

	}
	return unmarshaler.UnmarshalTraceContextStringMap(contextIDMap, tagsMap)
}
