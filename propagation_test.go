package opentracing

import (
	"net/http"
	"testing"
)

const testHeaderPrefix = "testprefix-"

func TestHTTPHeaderInject(t *testing.T) {
	h := http.Header{}
	h.Add("NotOT", "blah")
	h.Add("opname", "AlsoNotOT")
	tracer := testTracer{}
	span := tracer.StartSpan("someSpan")

	// Use HTTPHeaderTextMapCarrier to wrap around `h`.
	carrier := HTTPHeaderTextMapCarrier{
		HeaderPrefix: testHeaderPrefix,
		Header:       h,
	}
	if err := span.Tracer().Inject(span, TextMap, carrier); err != nil {
		t.Fatal(err)
	}

	if len(h) != 4 {
		t.Errorf("Unexpected header length: %v", len(h))
	}
	// The prefix comes from just above; the suffix comes from
	// testTracer.Inject().
	if h.Get("testprefix-opname") != "someSpan" {
		t.Errorf("Could not find opname at expected key")
	}
	if h.Get("testprefix-hasparent") != "false" {
		t.Errorf("Could not find hasparent at expected key")
	}
}

func TestHTTPHeaderJoin(t *testing.T) {
	h := http.Header{}
	h.Add("NotOT", "blah")
	h.Add("opname", "AlsoNotOT")
	h.Add("testprefix-opname", "someSpan")
	h.Add("testprefix-hasparent", "true")
	tracer := testTracer{}

	// Use HTTPHeaderTextMapCarrier to wrap around `h`.
	carrier := HTTPHeaderTextMapCarrier{
		HeaderPrefix: testHeaderPrefix,
		Header:       h,
	}
	span, err := tracer.Join("ignoredByImpl", TextMap, carrier)
	if err != nil {
		t.Fatal(err)
	}

	if !span.(testSpan).HasParent {
		t.Errorf("Failed to read testprefix-hasparent correctly")
	}
	if span.(testSpan).OperationName != "someSpan" {
		t.Errorf("Failed to read testprefix-opname correctly")
	}
}
