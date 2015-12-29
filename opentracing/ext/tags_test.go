package ext_test

import (
	"testing"

	"github.com/opentracing/api-golang/opentracing/ext"
	"github.com/opentracing/api-golang/opentracing/standardtracer"
	"github.com/opentracing/api-golang/testutils"
	"github.com/opentracing/api-golang/testutils/assert"
)

func TestPeerTags(t *testing.T) {
	if ext.PeerService.Key != "peer.service" {
		t.Fatalf("Invalid PeerService.Key %v", ext.PeerService.Key)
	}
	recorder := testutils.NewInMemoryRecorder("test-process")
	tracer := standardtracer.New(recorder, &testutils.SimpleTraceContextSource{})
	span := tracer.StartTrace("my-trace")
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
	assert.EqualValues(t, "my-service", rawSpan.Tags[ext.PeerService.Key])
	assert.EqualValues(t, "my-hostname", rawSpan.Tags[ext.PeerHostname.Key])
	assert.EqualValues(t, uint32(127<<24|1), rawSpan.Tags[ext.PeerHostIPv4.Key])
	assert.EqualValues(t, "::", rawSpan.Tags[ext.PeerHostIPv6.Key])
	assert.EqualValues(t, 8080, rawSpan.Tags[ext.PeerPort.Key])
}
