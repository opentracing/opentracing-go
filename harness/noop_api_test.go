package harness

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/suite"
)

func TestAPI(t *testing.T) {
	apiSuite := NewAPICheckSuite(func() (tracer opentracing.Tracer, closer func()) {
		return opentracing.NoopTracer{}, nil
	}, APICheckCapabilities{
		// NoopTracer doesn't do much
		CheckBaggageValues: false,
		CheckInject:        false,
		CheckExtract:       false,
	})
	suite.Run(t, apiSuite)
}
