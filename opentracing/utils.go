package opentracing

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func AddContextIDToHttpHeader(ctxID ContextID, h http.Header) {
	h.Add(OpenTracingContextIDHeader, base64.StdEncoding.EncodeToString(ctxID.Serialize()))
}

func GetContextIDFromHttpHeader(h http.Header, ctxIDSource ContextIDSource) (ContextID, error) {
	headerStr := h.Get(OpenTracingContextIDHeader)
	if len(headerStr) == 0 {
		return nil, fmt.Errorf("%q header not found", OpenTracingContextIDHeader)
	}
	ctxBytes, err := base64.StdEncoding.DecodeString(headerStr)
	if err != nil {
		return nil, err
	}
	return ctxIDSource.DeserializeContextID(ctxBytes)
}
