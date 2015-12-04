package opentracing

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	// OpenTracingContextHTTPHeaderPrefix precedes all opentracing-related HTTP
	// headers.
	OpenTracingContextHTTPHeaderPrefix = "Opentracing-Context-"
)

// AddTraceContextToHeader marshals TraceContext `ctx` to `h` as a series of
// HTTP headers. Values are URL-escaped.
func AddTraceContextToHeader(
	ctx TraceContext,
	h http.Header,
	marshaler TraceContextMarshaler,
) {
	for headerSuffix, val := range marshaler.MarshalTraceContextStringMap(ctx) {
		h.Add(OpenTracingContextHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
}

// TraceContextFromHeader unmarshals a TraceContext from `h`, expecting that
// header values are URL-escpaed.
func TraceContextFromHeader(
	h http.Header,
	unmarshaler TraceContextUnmarshaler,
) (TraceContext, error) {
	marshaled := make(map[string]string)
	for key, val := range h {
		if strings.HasPrefix(key, OpenTracingContextHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			marshaled[strings.TrimPrefix(key, OpenTracingContextHTTPHeaderPrefix)] = unescaped
		}
	}
	return unmarshaler.UnmarshalTraceContextStringMap(marshaled)
}
