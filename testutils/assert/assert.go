package assert

import (
	"reflect"
	"testing"
)

// NOTE: the functions in this module are borrowed from https://github.com/stretchr/testify
// to avoid creating a dependency (until Go comes up with proper dependency manager)

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    assert.EqualValues(t, uint32(123), int32(123), "123 and 123 should be equal")
//
// If assertion is not successful, execution halts via t.Fatalf() call
func EqualValues(t *testing.T, expected, actual interface{}) bool {
	if !ObjectsAreEqualValues(expected, actual) {
		t.Fatalf("Not equal: %#v (expected)\n"+
			"        != %#v (actual)", expected, actual)
	}
	return true
}

// ObjectsAreEqual determines if two objects are considered equal.
func ObjectsAreEqual(expected, actual interface{}) bool {

	if expected == nil || actual == nil {
		return expected == actual
	}

	if reflect.DeepEqual(expected, actual) {
		return true
	}

	return false

}

// ObjectsAreEqualValues determines whether two objects are equal,
// or if their values are equal.
func ObjectsAreEqualValues(expected, actual interface{}) bool {
	if ObjectsAreEqual(expected, actual) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.Type().ConvertibleTo(actualType) {
		// Attempt comparison after type conversion
		if reflect.DeepEqual(actual, expectedValue.Convert(actualType).Interface()) {
			return true
		}
	}

	return false
}
