package opentracing

import (
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
	// First, try to inject using the GoHTTPHeader format (our preference).
	if err := sp.Tracer().Inject(sp, GoHTTPHeader, h); err == nil {
		return nil
	}

	// Else, fall back on SplitText.
	carrier := NewSplitTextCarrier()
	if err := sp.Tracer().Inject(sp, SplitText, carrier); err != nil {
		return err
	}
	for headerSuffix, val := range carrier.TracerState {
		h.Add(ContextIDHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
	for headerSuffix, val := range carrier.Baggage {
		h.Add(TagsHTTPHeaderPrefix+headerSuffix, url.QueryEscape(val))
	}
	return nil
}

// JoinFromHeader decodes a Span with operation name `operationName` from `h`,
// expecting that header values are URL-escpaed.
//
// If `operationName` is empty, the caller must later call
// `Span.SetOperationName` on the returned `Span`.
func JoinFromHeader(
	operationName string,
	h http.Header,
	tracer Tracer,
) (Span, error) {
	// First, try to Join using the GoHTTPHeader format (our preference).
	span, err := tracer.Join(operationName, GoHTTPHeader, h)
	if err == nil {
		return span, nil
	}

	// Else, fall back on SplitText.
	carrier := NewSplitTextCarrier()
	for key, val := range h {
		if strings.HasPrefix(key, ContextIDHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			if carrier.TracerState == nil {
				carrier.TracerState = map[string]string{}
			}
			carrier.TracerState[strings.TrimPrefix(key, ContextIDHTTPHeaderPrefix)] = unescaped
		} else if strings.HasPrefix(key, TagsHTTPHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			if carrier.Baggage == nil {
				carrier.Baggage = map[string]string{}
			}
			carrier.Baggage[strings.TrimPrefix(key, TagsHTTPHeaderPrefix)] = unescaped
		}
	}
	return tracer.Join(operationName, SplitText, carrier)
}
