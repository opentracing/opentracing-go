package opentracing

import (
	"net/http"
	"testing"
)

func TestInjectSpanInHeader(t *testing.T) {
	h := http.Header{}
	h.Add("NotOT", "blah")
	h.Add("opname", "AlsoNotOT")
	tracer := testTracer{}
	span := tracer.StartSpan("someSpan")
	InjectSpanInHeader(span, h, "testprefix-")
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

func TestJoinFromHeader(t *testing.T) {
	h := http.Header{}
	h.Add("NotOT", "blah")
	h.Add("opname", "AlsoNotOT")
	h.Add("testprefix-opname", "someSpan")
	h.Add("testprefix-hasparent", "true")
	tracer := testTracer{}
	span, err := JoinFromHeader(tracer, "joinedOpname", h, "testprefix-")
	if err != nil {
		t.Fatal(err)
	}
	if !span.(testSpan).HasParent {
		t.Errorf("Failed to read testprefix-hasparent correctly")
	}
}
