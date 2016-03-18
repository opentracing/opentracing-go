package opentracing

import "net/http"

// InjectSpanInHeader encodes Span `sp` in `h` as a series of HTTP headers.
// Values are URL-escaped.
func InjectSpanInHeader(sp Span, h http.Header, headerPrefix string) error {
	// First, try to inject using the GoHTTPHeader format (our preference).
	if err := sp.Tracer().Inject(sp, GoHTTPHeader, h); err == nil {
		return nil
	}
	// Use HTTPHeaderTextMapCarrier to wrap around `h`.
	carrier := HTTPHeaderTextMapCarrier{
		HeaderPrefix: headerPrefix,
		Header:       h,
	}
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
func JoinFromHeader(tracer Tracer, operationName string, h http.Header, headerPrefix string) (Span, error) {
	// Use HTTPHeaderTextMapCarrier to wrap around `h`.
	carrier := HTTPHeaderTextMapCarrier{
		HeaderPrefix: headerPrefix,
		Header:       h,
	}
	return tracer.Join(operationName, TextMap, carrier)
}
