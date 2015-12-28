package assert_test

import (
	"github.com/opentracing/api-golang/testutils/assert"
	"testing"
)

func TestAssert(t *testing.T) {
	assert.EqualValues(t, int16(123), int16(123))
	assert.EqualValues(t, int16(123), int32(123))
}
