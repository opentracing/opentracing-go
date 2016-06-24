package opentracing

import (
	"net/http"
	"strconv"
	"testing"
)

const testHeaderPrefix = "testprefix-"

func TestTextMapCarrierInject(t *testing.T) {
	m := make(map[string]string)
	m["NotOT"] = "blah"
	m["opname"] = "AlsoNotOT"
	tracer := testTracer{}
	span := tracer.StartSpan("someSpan")
	fakeID := span.(testSpan).FakeID

	carrier := TextMapCarrier(m)
	if err := span.Tracer().Inject(span, TextMap, carrier); err != nil {
		t.Fatal(err)
	}

	if len(m) != 3 {
		t.Errorf("Unexpected header length: %v", len(m))
	}
	// The prefix comes from just above; the suffix comes from
	// testTracer.Inject().
	if m["testprefix-fakeid"] != strconv.Itoa(fakeID) {
		t.Errorf("Could not find fakeid at expected key")
	}
}

func TestTextMapCarrierJoin(t *testing.T) {
	m := make(map[string]string)
	m["NotOT"] = "blah"
	m["opname"] = "AlsoNotOT"
	m["testprefix-fakeid"] = "42"
	tracer := testTracer{}

	carrier := TextMapCarrier(m)
	span, err := tracer.Join("ignoredByImpl", TextMap, carrier)
	if err != nil {
		t.Fatal(err)
	}

	if span.(testSpan).FakeID != 42 {
		t.Errorf("Failed to read testprefix-fakeid correctly")
	}
}

func TestHTTPHeaderInject(t *testing.T) {
	h := http.Header{}
	h.Add("NotOT", "blah")
	h.Add("opname", "AlsoNotOT")
	tracer := testTracer{}
	span := tracer.StartSpan("someSpan")
	fakeID := span.(testSpan).FakeID

	// Use HTTPHeaderTextMapCarrier to wrap around `h`.
	carrier := HTTPHeaderTextMapCarrier(h)
	if err := span.Tracer().Inject(span, TextMap, carrier); err != nil {
		t.Fatal(err)
	}

	if len(h) != 3 {
		t.Errorf("Unexpected header length: %v", len(h))
	}
	// The prefix comes from just above; the suffix comes from
	// testTracer.Inject().
	if h.Get("testprefix-fakeid") != strconv.Itoa(fakeID) {
		t.Errorf("Could not find fakeid at expected key")
	}
}

func TestHTTPHeaderJoin(t *testing.T) {
	h := http.Header{}
	h.Add("NotOT", "blah")
	h.Add("opname", "AlsoNotOT")
	h.Add("testprefix-fakeid", "42")
	tracer := testTracer{}

	// Use HTTPHeaderTextMapCarrier to wrap around `h`.
	carrier := HTTPHeaderTextMapCarrier(h)
	span, err := tracer.Join("ignoredByImpl", TextMap, carrier)
	if err != nil {
		t.Fatal(err)
	}

	if span.(testSpan).FakeID != 42 {
		t.Errorf("Failed to read testprefix-fakeid correctly")
	}
}
