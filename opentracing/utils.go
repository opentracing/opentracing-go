package opentracing

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	OpenTracingContextHeaderPrefix = "Opentracing-Context-"
)

// AddTraceContextToHeader marshals TraceContext `ctx` to `h` as a series of
// HTTP headers. Values are URL-escaped.
func AddTraceContextToHeader(
	ctx TraceContext,
	h http.Header,
	marshaler TraceContextMarshaler,
) {
	for headerSuffix, val := range marshaler.MarshalStringMapTraceContext(ctx) {
		h.Add(OpenTracingContextHeaderPrefix+headerSuffix, url.QueryEscape(val))
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
		if strings.HasPrefix(key, OpenTracingContextHeaderPrefix) {
			// We don't know what to do with anything beyond slice item v[0]:
			unescaped, err := url.QueryUnescape(val[0])
			if err != nil {
				return nil, err
			}
			marshaled[strings.TrimPrefix(key, OpenTracingContextHeaderPrefix)] = unescaped
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
