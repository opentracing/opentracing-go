package dapperish

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/standardtracer"
)

// NewTracer returns a new dapperish Tracer instance.
func NewTracer(processName string) opentracing.Tracer {
	return standardtracer.New(
		NewTrivialRecorder(processName),
		&opentracing.ImplementationID{
			Name:    "dapperish",
			Version: "0.1.0",
		})
}
