package ext

import (
	"github.com/opentracing/opentracing-go"
)

// TODO move RPCServerOption into this file? OR, remove it. Code such as:
//
// err := tracer.Extract(...)
// if err != nil { ... }
// tracer.StartSpan("...", RPCServerOption(...))
//
// can be rewritten as:
// tracer.StartSpan("...", SpanKindRPCServer, ChildOfCarrierOption(...))

type extractedOption struct {
	ref             opentracing.SpanReferenceType
	format, carrier interface{}
}

// Apply implements the StartSpanOption interface.
func (e extractedOption) Apply(tracer opentracing.Tracer, o *opentracing.StartSpanOptions) {
	spanCtx, err := tracer.Extract(e.format, e.carrier)
	if spanCtx != nil {
		opentracing.SpanReference{e.ref, spanCtx}.Apply(tracer, o)
	}
	if err != nil {
		opentracing.Tag{string(CarriedContextError), err.Error()}.Apply(tracer, o)
	}
}

// CarriedSpanRefOption extracts a context from "carrier" (of type "format")
// and returns a StartSpanOption that will the associated context as a span
// reference of type "ref" to a newly started span.
func CarriedSpanRefOption(ref opentracing.SpanReferenceType,
	format, carrier interface{}) opentracing.StartSpanOption {
	return extractedOption{ref, format, carrier}
}

// CarriedChildOfOption constructs a opentracing.ChildOfRef SpanReference option
// for use in StartSpan with a carried context.
func CarriedChildOfOption(format, carrier interface{}) opentracing.StartSpanOption {
	return CarriedSpanRefOption(opentracing.ChildOfRef, format, carrier)
}

// FollowsFrom constructs a opentracing.FollowsFrom SpanReference option
// for use in StartSpan with a carried context.
func CarriedFollowsFromOption(format, carrier interface{}) opentracing.StartSpanOption {
	return CarriedSpanRefOption(opentracing.FollowsFromRef, format, carrier)
}
