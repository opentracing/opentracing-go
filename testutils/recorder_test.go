package testutils_test

import (
	"testing"
	"time"

	"github.com/opentracing/api-golang/examples/dapperish"
	"github.com/opentracing/api-golang/opentracing/standardtracer"
	"github.com/opentracing/api-golang/testutils"
)

func TestInMemoryRecorderSpans(t *testing.T) {
	recorder := testutils.NewInMemoryRecorder("unit-test")
	var apiRecorder standardtracer.Recorder = recorder
	if apiRecorder.ProcessName() != "unit-test" {
		t.Fatalf("Invalid process name")
	}
	span := &standardtracer.RawSpan{
		TraceContext: &dapperish.TraceContext{},
		Operation:    "test-span",
		Start:        time.Now(),
		Duration:     -1,
	}
	apiRecorder.RecordSpan(span)
	if len(recorder.GetSpans()) != 1 {
		t.Fatal("No spans recorded")
	}
	if recorder.GetSpans()[0] != span {
		t.Fatal("Span not recorded")
	}
}

func TestInMemoryRecorderTags(t *testing.T) {
	recorder := testutils.NewInMemoryRecorder("unit-test")
	recorder.SetTag("tag1", "hello")
	if len(recorder.GetTags()) != 1 {
		t.Fatal("Tag not stored")
	}
	if recorder.GetTags()["tag1"] != "hello" {
		t.Fatal("tag1 != hello")
	}
}
