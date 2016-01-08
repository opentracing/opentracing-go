package testutils

import (
	"github.com/opentracing/api-go/opentracing"
)

// SimpleTraceContextSource is a dummy implementation of TraceContextSource.
type SimpleTraceContextSource struct{}

// NewRootTraceContext implements NewRootTraceContext of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) NewRootTraceContext() opentracing.TraceContext {
	return nil
}

// NewChildTraceContext implements NewChildTraceContext of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) NewChildTraceContext(parent opentracing.TraceContext) (opentracing.TraceContext, opentracing.Tags) {
	return nil, nil
}

// SimpleTraceContextSource implements opentracing.TraceContextEncoder.
func (source *SimpleTraceContextSource) TraceContextToBinary(
	tc opentracing.TraceContext,
) (
	traceContextID []byte,
	traceAttrs []byte,
) {
	panic("Not implemented")
}

// SimpleTraceContextSource implements opentracing.TraceContextEncoder.
func (source *SimpleTraceContextSource) TraceContextToText(
	tc opentracing.TraceContext,
) (
	traceContextID map[string]string,
	traceAttrs map[string]string,
) {
	panic("Not implemented")
}

// SimpleTraceContextSource implements opentracing.TraceContextDecoder.
func (source *SimpleTraceContextSource) TraceContextFromBinary(
	traceContextID []byte,
	traceAttrs []byte,
) (opentracing.TraceContext, error) {
	panic("Not implemented")
}

// SimpleTraceContextSource implements opentracing.TraceContextDecoder.
func (source *SimpleTraceContextSource) TraceContextFromText(
	traceContextID map[string]string,
	traceAttrs map[string]string,
) (opentracing.TraceContext, error) {
	panic("Not implemented")
}
