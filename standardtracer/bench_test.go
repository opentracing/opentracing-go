package standardtracer

import (
	"fmt"
	"sync/atomic"
	"testing"
)

var tags []string

func init() {
	tags = make([]string, 1000)
	for j := 0; j < len(tags); j++ {
		tags[j] = fmt.Sprintf("%d", randomID())
	}
}

type countingRecorder int32

func (c *countingRecorder) RecordSpan(r RawSpan) {
	atomic.AddInt32((*int32)(c), 1)
}

func benchmarkWithOps(b *testing.B, numEvent, numTag, numAttr int) {
	var r countingRecorder
	t := New(&r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sp := t.StartSpan("test")
		for j := 0; j < numEvent; j++ {
			sp.LogEvent("event")
		}
		for j := 0; j < numTag; j++ {
			sp.SetTag(tags[j], nil)
		}
		for j := 0; j < numAttr; j++ {
			sp.SetTraceAttribute(tags[j], tags[j])
		}
		sp.Finish()
	}
	b.StopTimer()
	if int(r) != b.N {
		b.Fatalf("missing traces: expected %d, got %d", b.N, r)
	}
}

func BenchmarkSpan_Empty(b *testing.B) {
	benchmarkWithOps(b, 0, 0, 0)
}

func BenchmarkSpan_100Events(b *testing.B) {
	benchmarkWithOps(b, 100, 0, 0)
}

func BenchmarkSpan_1000Events(b *testing.B) {
	benchmarkWithOps(b, 100, 0, 0)
}

func BenchmarkSpan_100Tags(b *testing.B) {
	benchmarkWithOps(b, 0, 100, 0)
}

func BenchmarkSpan_1000Tags(b *testing.B) {
	benchmarkWithOps(b, 0, 100, 0)
}
