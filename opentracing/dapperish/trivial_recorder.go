package main

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/opentracing/api-golang/opentracing"
)

type TrivialRecorder struct {
	tags map[string]string
}

func NewTrivialRecorder() *TrivialRecorder {
	return &TrivialRecorder{
		tags: make(map[string]string),
	}
}

func (t *TrivialRecorder) SetTag(key string, val interface{}) {
	t.tags[key] = fmt.Sprint(val)
}

func (t *TrivialRecorder) RecordSpan(span *opentracing.RawSpan) {
	str := base64.StdEncoding.EncodeToString(span.ContextID.Serialize())

	fmt.Printf(
		"RecordSpan: %v[%v, %v us] --> %v logs. context base64: %v\n",
		span.Operation, span.Start, span.Duration, len(span.Logs), str)
	for i, l := range span.Logs {
		fmt.Printf(
			"    log %v: %v --> %v\n", i, l.Message, reflect.TypeOf(l.Payload))
	}
}
