package testutils

import "testing"

func TestSimpleTraceContextSource(t *testing.T) {
	ctxSrc := &SimpleTraceContextSource{}
	ctxSrc.NewRootTraceContext()
}
