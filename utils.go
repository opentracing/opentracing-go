package opentracing

import (
	"net/http"
	"net/url"
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

	// Else, fall back on TextMap.
	carrier := TextMapCarrier{}
	if err := sp.Tracer().Inject(sp, TextMap, carrier); err != nil {
		return err
	}
	for key, val := range carrier {
		h.Add(key, url.QueryEscape(val))
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

	// Else, fall back on TextMap.
	carrier := TextMapCarrier{}
	for key, val := range h {
		// We don't know what to do with anything beyond slice item v[0]:
		unescaped, err := url.QueryUnescape(val[0])
		if err != nil {
			continue
		}
		carrier[key] = unescaped
	}
	return tracer.Join(operationName, TextMap, carrier)
}
