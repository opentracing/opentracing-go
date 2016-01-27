package dapperish

import (
	"fmt"
	"reflect"

	"github.com/opentracing/opentracing-go/standardtracer"
)

// TrivialRecorder implements the standardtracer.Recorder interface.
type TrivialRecorder struct {
	processName string
	tags        map[string]string
}

// NewTrivialRecorder returns a TrivialRecorder for the given `processName`.
func NewTrivialRecorder(processName string) *TrivialRecorder {
	return &TrivialRecorder{
		processName: processName,
		tags:        make(map[string]string),
	}
}

func (t *TrivialRecorder) ProcessName() string { return t.processName }

func (t *TrivialRecorder) SetTag(key string, val interface{}) *TrivialRecorder {
	t.tags[key] = fmt.Sprint(val)
	return t
}

// RecordSpan complies with the standardtracer.Recorder interface.
func (t *TrivialRecorder) RecordSpan(span *standardtracer.RawSpan) {
	fmt.Printf(
		"RecordSpan: %v[%v, %v us] --> %v logs. std context: %v\n",
		span.Operation, span.Start, span.Duration, len(span.Logs),
		span.StandardContext)
	for i, l := range span.Logs {
		fmt.Printf(
			"    log %v @ %v: %v --> %v\n", i, l.Timestamp, l.Event, reflect.TypeOf(l.Payload))
	}
}
