package opentracing

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

const (
	OpenTracingContextHeader = "OpenTracing-Context"
)

func AddContextToHttpHeader(ctx TraceContext, h http.Header) {
	h.Add(OpenTracingContextHeader, ctxID.SerializeString())
}

func GetTraceContextFromHttpHeader(
	h http.Header,
	ctxIDSource TraceContextIDSource,
) (TraceContext, error) {
	headerStr := h.Get(OpenTracingContextHeader)
	if len(headerStr) == 0 {
		return nil, fmt.Errorf("%q header not found", OpenTracingContextHeader)
	}
	ctxBytes, err := base64.StdEncoding.DecodeString(headerStr)
	if err != nil {
		return nil, err
	}
	return DeserializeStringTraceContext(ctxBytes)
}
