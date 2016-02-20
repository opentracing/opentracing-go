package ext_test

import (
	"reflect"
	"testing"

	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/basictracer-go/testutils"
	"github.com/opentracing/opentracing-go/ext"
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
	recorder := testutils.NewInMemoryRecorder()
	tracer := basictracer.New(recorder)
	span := tracer.StartSpan("my-trace")
	ext.PeerService.Add(span, "my-service")
	ext.PeerHostname.Add(span, "my-hostname")
	ext.PeerHostIPv4.Add(span, uint32(127<<24|1))
	ext.PeerHostIPv6.Add(span, "::")
	ext.PeerPort.Add(span, uint16(8080))
	span.Finish()
	if len(recorder.GetSpans()) != 1 {
		t.Fatal("Span not recorded")
	}
	rawSpan := recorder.GetSpans()[0]
	assertEqual(t, "my-service", rawSpan.Tags[string(ext.PeerService)])
	assertEqual(t, "my-hostname", rawSpan.Tags[string(ext.PeerHostname)])
	assertEqual(t, uint32(127<<24|1), rawSpan.Tags[string(ext.PeerHostIPv4)])
	assertEqual(t, "::", rawSpan.Tags[string(ext.PeerHostIPv6)])
	assertEqual(t, uint16(8080), rawSpan.Tags[string(ext.PeerPort)])
}
