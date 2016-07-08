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

func CarriedSpanRefOption(ref opentracing.SpanReferenceType,
	format, carrier interface{}) opentracing.StartSpanOption {
	return extractedOption{ref, format, carrier}
}

func CarriedChildOfOption(format, carrier interface{}) opentracing.StartSpanOption {
	return CarriedSpanRefOption(opentracing.ChildOfRef, format, carrier)
}

func CarriedFollowsFromOption(format, carrier interface{}) opentracing.StartSpanOption {
	return CarriedSpanRefOption(opentracing.FollowsFromRef, format, carrier)
}
