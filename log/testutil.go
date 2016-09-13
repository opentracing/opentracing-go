package log

import (
	"fmt"
	"reflect"
	"testing"
)

// LogFieldValidator facilitates testing of Span.Log*() implementations.
//
// LogFieldValidator SHOULD ONLY BE USED FOR TESTS. It is not intended for
// production code.
//
// Usage:
//
//     lfv := log.NewLogFieldValidator(t, someLogStructure.Fields)
//     lfv.
//         ExpectNextFieldEquals("key1", reflect.String, "some string value").
//         ExpectNextFieldEquals("key2", reflect.Uint8, "255")
//
// LogFieldValidator satisfies the log.Encoder interface and thus is able to
// marshal log.Field instances (which it takes advantage of internally).
type LogFieldValidator struct {
	t               *testing.T
	fieldIdx        int
	fields          []Field
	nextKey         string
	nextKind        reflect.Kind
	nextValAsString string
}

// NewLogFieldValidator returns a new validator that will test the contents of
// `fields`.
func NewLogFieldValidator(t *testing.T, fields []Field) *LogFieldValidator {
	return &LogFieldValidator{
		t:      t,
		fields: fields,
	}
}

// ExpectNextFieldEquals facilitates a fluent way of testing the contents
// []Field slices.
func (lfv *LogFieldValidator) ExpectNextFieldEquals(key string, kind reflect.Kind, valAsString string) *LogFieldValidator {
	if len(lfv.fields) < lfv.fieldIdx {
		lfv.t.Errorf("Expecting more than the %v Fields we have", len(lfv.fields))
	}
	lfv.nextKey = key
	lfv.nextKind = kind
	lfv.nextValAsString = valAsString
	lfv.fields[lfv.fieldIdx].Marshal(lfv)
	lfv.fieldIdx++
	return lfv
}

// EmitString satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitString(key, value string) {
	lfv.validateNextField(key, reflect.String, value)
}

// EmitBool satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitBool(key string, value bool) {
	lfv.validateNextField(key, reflect.Bool, value)
}

// EmitInt satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitInt(key string, value int) {
	lfv.validateNextField(key, reflect.Int, value)
}

// EmitInt32 satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitInt32(key string, value int32) {
	lfv.validateNextField(key, reflect.Int32, value)
}

// EmitInt64 satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitInt64(key string, value int64) {
	lfv.validateNextField(key, reflect.Int64, value)
}

// EmitUint32 satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitUint32(key string, value uint32) {
	lfv.validateNextField(key, reflect.Uint32, value)
}

// EmitUint64 satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitUint64(key string, value uint64) {
	lfv.validateNextField(key, reflect.Uint64, value)
}

// EmitFloat32 satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitFloat32(key string, value float32) {
	lfv.validateNextField(key, reflect.Float32, value)
}

// EmitFloat64 satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitFloat64(key string, value float64) {
	lfv.validateNextField(key, reflect.Float64, value)
}

// EmitObject satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitObject(key string, value interface{}) {
	lfv.validateNextField(key, reflect.Interface, value)
}

// EmitLazyLogger satisfies the Encoder interface
func (lfv *LogFieldValidator) EmitLazyLogger(key string, value LazyLogger) {
	lfv.t.Error("Test infrastructure does not support EmitLazyLogger yet")
}

func (lfv *LogFieldValidator) validateNextField(key string, actualKind reflect.Kind, value interface{}) {
	if lfv.nextKey != key {
		lfv.t.Errorf("Bad key: expected %q, found %q", lfv.nextKey, key)
	}
	if lfv.nextKind != actualKind {
		lfv.t.Errorf("Bad reflect.Kind: expected %v, found %v", lfv.nextKind, actualKind)
		return
	}
	if lfv.nextValAsString != fmt.Sprint(value) {
		lfv.t.Errorf("Bad value: expected %q, found %q", lfv.nextValAsString, fmt.Sprint(value))
	}
	// All good.
}
