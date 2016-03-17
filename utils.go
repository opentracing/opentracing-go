package opentracing

import (
	"net/http"
	"net/url"
)

// InjectSpanInHeader encodes Span `sp` in `h` as a series of HTTP headers.
// Values are URL-escaped.
func InjectSpanInHeader(sp Span, h http.Header) error {
	// First, try to inject using the GoHTTPHeader format (our preference).
	if err := sp.Tracer().Inject(sp, GoHTTPHeader, h); err == nil {
		return nil
	}

	// Else, fall back on TextMap.
	carrier := HTTPHeaderTextMapCarrier{http.Header{}}
	if err := sp.Tracer().Inject(sp, TextMap, carrier); err != nil {
		return err
	}
	return nil
}

// JoinFromHeader decodes a Span with operation name `operationName` from `h`,
// expecting that header values are URL-escpaed.
//
// If `operationName` is empty, the caller must later call
// `Span.SetOperationName` on the returned `Span`.
func JoinFromHeader(operationName string, h http.Header, tracer Tracer) (Span, error) {
	// First, try to Join using the GoHTTPHeader format (our preference).
	span, err := tracer.Join(operationName, GoHTTPHeader, h)
	if err == nil {
		return span, nil
	}

	// Else, fall back on TextMap.
	carrier := HTTPHeaderTextMapCarrier{h}
	return tracer.Join(operationName, TextMap, carrier)
}

type HTTPHeaderTextMapCarrier struct {
	http.Header
}

func (c HTTPHeaderTextMapCarrier) Add(key, val string) {
	c.Header.Add(key, url.QueryEscape(val))
}
func (c HTTPHeaderTextMapCarrier) GetAll(handler func(key, val string) error) error {
	for k, vals := range c.Header {
		for _, v := range vals {
			rawV, err := url.QueryUnescape(v)
			if err != nil {
				continue
			}
			if err = handler(k, rawV); err != nil {
				return err
			}
		}
	}
	return nil
}
