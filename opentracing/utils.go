package opentracing

import (
	"fmt"
	"net/http"
)

const (
	OpenTracingContextHeader = "OpenTracing-Context"
)

func AddTraceContextToHttpHeader(ctx *TraceContext, h http.Header) {
	h.Add(OpenTracingContextHeader, ctx.SerializeString())
}

func GetTraceContextFromHttpHeader(
	h http.Header,
	ctxIDSource TraceContextIDSource,
) (*TraceContext, error) {
	headerStr := h.Get(OpenTracingContextHeader)
	if len(headerStr) == 0 {
		return nil, fmt.Errorf("%q header not found", OpenTracingContextHeader)
	}
	return DeserializeStringTraceContext(ctxIDSource, headerStr)
}
