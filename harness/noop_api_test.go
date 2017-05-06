package harness

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
)

var noopOpts = CheckOpts{CheckBaggageValues: false}
var noopTracer = opentracing.NoopTracer{}

func TestNoopStartSpan(t *testing.T)           { CheckStartSpan(t, noopTracer, noopOpts) }
func TestNoopStartSpanWithParent(t *testing.T) { CheckStartSpanWithParent(t, noopTracer, noopOpts) }
func TestSetOperationName(t *testing.T)        { CheckSetOperationName(t, noopTracer, noopOpts) }
func TestSpanTagValueTypes(t *testing.T)       { CheckSpanTagValueTypes(t, noopTracer, noopOpts) }
func TestSpanTagsWithChaining(t *testing.T)    { CheckSpanTagsWithChaining(t, noopTracer, noopOpts) }
func TestSpanLogs(t *testing.T)                { CheckSpanLogs(t, noopTracer, noopOpts) }
func TestSpanBaggage(t *testing.T)             { CheckSpanBaggage(t, noopTracer, noopOpts) }
func TestContextBaggage(t *testing.T)          { CheckContextBaggage(t, noopTracer, noopOpts) }
func TestTextPropagation(t *testing.T)         { CheckTextPropagation(t, noopTracer, noopOpts) }
func TestHTTPPropagation(t *testing.T)         { CheckHTTPPropagation(t, noopTracer, noopOpts) }
func TestBinaryPropagation(t *testing.T)       { CheckBinaryPropagation(t, noopTracer, noopOpts) }
func TestMandatoryFormats(t *testing.T)        { CheckMandatoryFormats(t, noopTracer, noopOpts) }
func TestUnknownFormat(t *testing.T)           { CheckUnknownFormat(t, noopTracer, noopOpts) }
