package ext_test

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/standardtracer"
	"github.com/opentracing/opentracing-go/testutils"
	"github.com/opentracing/opentracing-go/testutils/assert"
)

func TestPeerTags(t *testing.T) {
	if ext.PeerService != "peer.service" {
		t.Fatalf("Invalid PeerService %v", ext.PeerService)
	}
	recorder := testutils.NewInMemoryRecorder()
	tracer := standardtracer.New(recorder, &opentracing.ImplementationID{})
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
	assert.EqualValues(t, "my-service", rawSpan.Tags[string(ext.PeerService)])
	assert.EqualValues(t, "my-hostname", rawSpan.Tags[string(ext.PeerHostname)])
	assert.EqualValues(t, uint32(127<<24|1), rawSpan.Tags[string(ext.PeerHostIPv4)])
	assert.EqualValues(t, "::", rawSpan.Tags[string(ext.PeerHostIPv6)])
	assert.EqualValues(t, 8080, rawSpan.Tags[string(ext.PeerPort)])
}
