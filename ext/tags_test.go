package ext_test

import (
	"reflect"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
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
	tracer := noopTracer{}
	span := tracer.StartSpan("my-trace")
	ext.PeerService.Set(span, "my-service")
	ext.PeerHostname.Set(span, "my-hostname")
	ext.PeerHostIPv4.Set(span, uint32(127<<24|1))
	ext.PeerHostIPv6.Set(span, "::")
	ext.PeerPort.Set(span, uint16(8080))
	ext.SamplingPriority.Set(span, uint16(1))
	ext.SpanKind.Set(span, ext.SpanKindRPCServer)
	ext.SpanKind.Set(span, ext.SpanKindRPCClient)
	span.Finish()

	rawSpan := span.(*noopSpan)
	assertEqual(t, "my-service", rawSpan.Tags[string(ext.PeerService)])
	assertEqual(t, "my-hostname", rawSpan.Tags[string(ext.PeerHostname)])
	assertEqual(t, uint32(127<<24|1), rawSpan.Tags[string(ext.PeerHostIPv4)])
	assertEqual(t, "::", rawSpan.Tags[string(ext.PeerHostIPv6)])
	assertEqual(t, uint16(8080), rawSpan.Tags[string(ext.PeerPort)])
}

// noopTracer and noopSpan with span tags implemented
type noopTracer struct{}

type noopSpan struct {
	Tags opentracing.Tags
}

type noopInjectorExtractor struct{}

func (n noopSpan) SetTag(key string, value interface{}) opentracing.Span {
	n.Tags[key] = value
	return n
}

func (n noopSpan) Finish()                                                {}
func (n noopSpan) FinishWithOptions(opts opentracing.FinishOptions)       {}
func (n noopSpan) SetBaggageItem(key, val string) opentracing.Span        { return n }
func (n noopSpan) BaggageItem(key string) string                          { return "" }
func (n noopSpan) LogEvent(event string)                                  {}
func (n noopSpan) LogEventWithPayload(event string, payload interface{})  {}
func (n noopSpan) Log(data opentracing.LogData)                           {}
func (n noopSpan) SetOperationName(operationName string) opentracing.Span { return n }
func (n noopSpan) Tracer() opentracing.Tracer                             { return nil }

func (n noopTracer) StartSpan(operationName string) opentracing.Span {
	return &noopSpan{Tags: make(opentracing.Tags)}
}

func (n noopTracer) StartSpanWithOptions(opts opentracing.StartSpanOptions) opentracing.Span {
	return noopSpan{Tags: make(opentracing.Tags)}
}

func (n noopTracer) Extractor(format interface{}) opentracing.Extractor {
	return noopInjectorExtractor{}
}

func (n noopTracer) Injector(format interface{}) opentracing.Injector {
	return noopInjectorExtractor{}
}

func (n noopInjectorExtractor) InjectSpan(span opentracing.Span, carrier interface{}) error {
	return nil
}

func (n noopInjectorExtractor) JoinTrace(operationName string, carrier interface{}) (opentracing.Span, error) {
	panic("not implemented")
}
