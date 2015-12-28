package testutils

import (
	"github.com/opentracing/api-golang/opentracing"
)

// SimpleTraceContextSource is a dummy implementation of TraceContextSource.
type SimpleTraceContextSource struct{}

// NewRootTraceContext imlements NewRootTraceContext of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) NewRootTraceContext() opentracing.TraceContext {
	return nil
}

// MarshalTraceContextBinary imlements MarshalTraceContextBinary of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) MarshalTraceContextBinary(
	tc opentracing.TraceContext,
) (
	traceContextID []byte,
	traceTags []byte,
) {
	panic("Not implemented")
}

// MarshalTraceContextStringMap imlements MarshalTraceContextStringMap of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) MarshalTraceContextStringMap(
	tc opentracing.TraceContext,
) (
	traceContextID map[string]string,
	traceTags map[string]string,
) {
	panic("Not implemented")
}

// UnmarshalTraceContextBinary imlements UnmarshalTraceContextBinary of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) UnmarshalTraceContextBinary(
	traceContextID []byte,
	traceTags []byte,
) (opentracing.TraceContext, error) {
	panic("Not implemented")
}

// UnmarshalTraceContextStringMap imlements UnmarshalTraceContextStringMap of opentracing.TraceContextSource.
func (source *SimpleTraceContextSource) UnmarshalTraceContextStringMap(
	traceContextID map[string]string,
	traceTags map[string]string,
) (opentracing.TraceContext, error) {
	panic("Not implemented")
}
