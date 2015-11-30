package main

import (
	"fmt"
	"reflect"

	"github.com/opentracing/api-golang/opentracing"
)

type TrivialRecorder struct {
	processName string
	tags        map[string]string
}

func NewTrivialRecorder(processName string) *TrivialRecorder {
	return &TrivialRecorder{
		processName: processName,
		tags:        make(map[string]string),
	}
}

func (t *TrivialRecorder) ProcessName() string { return t.processName }

func (t *TrivialRecorder) SetTag(key string, val interface{}) {
	t.tags[key] = fmt.Sprint(val)
}

func (t *TrivialRecorder) RecordSpan(span *opentracing.RawSpan) {
	fmt.Printf(
		"RecordSpan: %v[%v, %v us] --> %v logs. trace context: %v\n",
		span.Operation, span.Start, span.Duration, len(span.Logs),
		span.TraceContext.SerializeString())
	for i, l := range span.Logs {
		fmt.Printf(
			"    log %v: %v --> %v\n", i, l.Message, reflect.TypeOf(l.Payload))
	}
}
