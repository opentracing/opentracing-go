package ext_test

import (
	"testing"

	"github.com/opentracing/api-golang/opentracing/ext"
	"github.com/opentracing/api-golang/opentracing/standardtracer"
	"github.com/opentracing/api-golang/testutils"
	"github.com/opentracing/api-golang/testutils/assert"
)

func TestRPCTags(t *testing.T) {
	if ext.RPCService.Key != "rpc.service" {
		t.Fatalf("Invalid RPCService.Key %v", ext.RPCService.Key)
	}
	recorder := testutils.NewInMemoryRecorder("test-process")
	tracer := standardtracer.New(recorder, &testutils.SimpleTraceContextSource{})
	span := tracer.StartTrace("my-trace")
	ext.RPCService.Add(span, "my-service")
	ext.RPCHostname.Add(span, "my-hostname")
	ext.RPCHostIPv4.Add(span, uint32(127<<24|1))
	ext.RPCHostIPv6.Add(span, "::")
	ext.RPCPort.Add(span, uint16(8080))
	span.Finish()
	if len(recorder.GetSpans()) != 1 {
		t.Fatal("Span not recorded")
	}
	rawSpan := recorder.GetSpans()[0]
	assert.EqualValues(t, "my-service", rawSpan.Tags[ext.RPCService.Key])
	assert.EqualValues(t, "my-hostname", rawSpan.Tags[ext.RPCHostname.Key])
	assert.EqualValues(t, uint32(127<<24|1), rawSpan.Tags[ext.RPCHostIPv4.Key])
	assert.EqualValues(t, "::", rawSpan.Tags[ext.RPCHostIPv6.Key])
	assert.EqualValues(t, 8080, rawSpan.Tags[ext.RPCPort.Key])
}
