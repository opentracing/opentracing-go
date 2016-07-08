package ext_test

import (
	"reflect"
	"testing"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func assertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Not equal: %#v (expected)\n"+
			"        != %#v (actual)", expected, actual)
	}
}

func TestPeerTags(t *testing.T) {
	if ext.PeerService != "peer.service" {
		t.Fatalf("Invalid PeerService %v", ext.PeerService)
	}
	tracer := mocktracer.New()
	span := tracer.StartSpan("my-trace")
	ext.PeerService.Set(span, "my-service")
	ext.PeerHostname.Set(span, "my-hostname")
	ext.PeerHostIPv4.Set(span, uint32(127<<24|1))
	ext.PeerHostIPv6.Set(span, "::")
	ext.PeerPort.Set(span, uint16(8080))
	ext.SamplingPriority.Set(span, uint16(1))
	ext.SpanKind.Set(span, ext.SpanKindRPCServerEnum)
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	span.Finish()

	rawSpan := tracer.GetFinishedSpans()[0]
	assertEqual(t, map[string]interface{}{
		"peer.service":      "my-service",
		"peer.hostname":     "my-hostname",
		"peer.ipv4":         uint32(127<<24 | 1),
		"peer.ipv6":         "::",
		"peer.port":         uint16(8080),
		"span.kind":         ext.SpanKindRPCClientEnum,
		"sampling.priority": uint16(1),
	}, rawSpan.GetTags())
}

func TestHTTPTags(t *testing.T) {
	tracer := mocktracer.New()
	span := tracer.StartSpan("my-trace", ext.SpanKindRPCServer)
	ext.HTTPUrl.Set(span, "test.biz/uri?protocol=false")
	ext.HTTPMethod.Set(span, "GET")
	ext.HTTPStatusCode.Set(span, 301)
	span.Finish()

	rawSpan := tracer.GetFinishedSpans()[0]
	assertEqual(t, map[string]interface{}{
		"http.url":         "test.biz/uri?protocol=false",
		"http.method":      "GET",
		"http.status_code": uint16(301),
		"span.kind":        ext.SpanKindRPCServerEnum,
	}, rawSpan.GetTags())
}

func TestMiscTags(t *testing.T) {
	tracer := mocktracer.New()
	span := tracer.StartSpan("my-trace")
	ext.Component.Set(span, "my-awesome-library")
	ext.SamplingPriority.Set(span, 1)
	ext.Error.Set(span, true)

	span.Finish()

	rawSpan := tracer.GetFinishedSpans()[0]
	assertEqual(t, map[string]interface{}{
		"component":         "my-awesome-library",
		"sampling.priority": uint16(1),
		"error":             true,
	}, rawSpan.GetTags())
}
