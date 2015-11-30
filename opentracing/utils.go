package opentracing

import (
	"fmt"
	"net/http"
	"reflect"
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

func keyValueListToTags(keyValueTags []interface{}) Tags {
	if len(keyValueTags)%2 != 0 {
		panic(fmt.Errorf(
			"there must be an even number of keyValueTags params to split them into pairs: got %v",
			len(keyValueTags)))
	}
	rval := make(Tags, len(keyValueTags)/2)
	var k string
	for i, kOrV := range keyValueTags {
		if i%2 == 0 {
			var ok bool
			k, ok = kOrV.(string)
			if !ok {
				panic(fmt.Errorf(
					"even-indexed keyValueTags (i.e., the keys) must be strings: got %v",
					reflect.TypeOf(kOrV)))
			}
		} else {
			rval[k] = kOrV
		}
	}
	return rval
}
