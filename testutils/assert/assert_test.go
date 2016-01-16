package assert_test

import (
	"testing"

	"github.com/opentracing/opentracing-go/testutils/assert"
)

func TestAssert(t *testing.T) {
	assert.EqualValues(t, int16(123), int16(123))
	assert.EqualValues(t, int16(123), int32(123))
}
