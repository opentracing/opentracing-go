package testutils

import (
	"github.com/opentracing/api-golang/opentracing"
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

// MarshalTraceContextBinary implements MarshalTraceContextBinary of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) MarshalTraceContextBinary(
	tc opentracing.TraceContext,
) (
	traceContextID []byte,
	traceAttrs []byte,
) {
	panic("Not implemented")
}

// MarshalTraceContextStringMap implements MarshalTraceContextStringMap of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) MarshalTraceContextStringMap(
	tc opentracing.TraceContext,
) (
	traceContextID map[string]string,
	traceAttrs map[string]string,
) {
	panic("Not implemented")
}

// UnmarshalTraceContextBinary implements UnmarshalTraceContextBinary of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) UnmarshalTraceContextBinary(
	traceContextID []byte,
	traceAttrs []byte,
) (opentracing.TraceContext, error) {
	panic("Not implemented")
}

// UnmarshalTraceContextStringMap implements UnmarshalTraceContextStringMap of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) UnmarshalTraceContextStringMap(
	traceContextID map[string]string,
	traceAttrs map[string]string,
) (opentracing.TraceContext, error) {
	panic("Not implemented")
}
