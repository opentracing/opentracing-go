package opentracing

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const (
	OpenTracingContextHeaderPrefix = "Opentracing-Context-"
)

func AddTraceContextToHeader(
	ctx TraceContext,
	h http.Header,
	marshaler TraceContextMarshaler,
) {
	for headerSuffix, val := range marshaler.MarshalStringMapTraceContext(ctx) {
		h.Add(OpenTracingContextHeaderPrefix+headerSuffix, val) // XXX escape val
	}
}

func TraceContextFromHeader(
	h http.Header,
	unmarshaler TraceContextUnmarshaler,
) (TraceContext, error) {
	marshaled := make(map[string]string)
	for key, val := range h {
		if strings.HasPrefix(key, OpenTracingContextHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			marshaled[strings.TrimPrefix(key, OpenTracingContextHeaderPrefix)] = val[0]
		}
	}
	return unmarshaler.UnmarshalStringMapTraceContext(marshaled)
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
