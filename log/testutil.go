package log

import (
	"fmt"
	"reflect"
	"testing"
)

// FieldValidator facilitates testing of Span.Log*() implementations.
//
// FieldValidator SHOULD ONLY BE USED FOR TESTS. It is not intended for
// production code.
//
// Usage:
//
//     fv := log.NewFieldValidator(t, someLogStructure.Fields)
//     fv.
//         ExpectNextFieldEquals("key1", reflect.String, "some string value").
//         ExpectNextFieldEquals("key2", reflect.Uint32, "4294967295")
//
// FieldValidator satisfies the log.Encoder interface and thus is able to
// marshal log.Field instances (which it takes advantage of internally).
type FieldValidator struct {
	t               *testing.T
	fieldIdx        int
	fields          []Field
	nextKey         string
	nextKind        reflect.Kind
	nextValAsString string
}

// NewFieldValidator returns a new validator that will test the contents of
// `fields`.
func NewFieldValidator(t *testing.T, fields []Field) *FieldValidator {
	return &FieldValidator{
		t:      t,
		fields: fields,
	}
}

// ExpectNextFieldEquals facilitates a fluent way of testing the contents
// []Field slices.
func (fv *FieldValidator) ExpectNextFieldEquals(key string, kind reflect.Kind, valAsString string) *FieldValidator {
	if len(fv.fields) < fv.fieldIdx {
		fv.t.Errorf("Expecting more than the %v Fields we have", len(fv.fields))
	}
	fv.nextKey = key
	fv.nextKind = kind
	fv.nextValAsString = valAsString
	fv.fields[fv.fieldIdx].Marshal(fv)
	fv.fieldIdx++
	return fv
}

// EmitString satisfies the Encoder interface
func (fv *FieldValidator) EmitString(key, value string) {
	fv.validateNextField(key, reflect.String, value)
}

// EmitBool satisfies the Encoder interface
func (fv *FieldValidator) EmitBool(key string, value bool) {
	fv.validateNextField(key, reflect.Bool, value)
}

// EmitInt satisfies the Encoder interface
func (fv *FieldValidator) EmitInt(key string, value int) {
	fv.validateNextField(key, reflect.Int, value)
}

// EmitInt32 satisfies the Encoder interface
func (fv *FieldValidator) EmitInt32(key string, value int32) {
	fv.validateNextField(key, reflect.Int32, value)
}

// EmitInt64 satisfies the Encoder interface
func (fv *FieldValidator) EmitInt64(key string, value int64) {
	fv.validateNextField(key, reflect.Int64, value)
}

// EmitUint32 satisfies the Encoder interface
func (fv *FieldValidator) EmitUint32(key string, value uint32) {
	fv.validateNextField(key, reflect.Uint32, value)
}

// EmitUint64 satisfies the Encoder interface
func (fv *FieldValidator) EmitUint64(key string, value uint64) {
	fv.validateNextField(key, reflect.Uint64, value)
}

// EmitFloat32 satisfies the Encoder interface
func (fv *FieldValidator) EmitFloat32(key string, value float32) {
	fv.validateNextField(key, reflect.Float32, value)
}

// EmitFloat64 satisfies the Encoder interface
func (fv *FieldValidator) EmitFloat64(key string, value float64) {
	fv.validateNextField(key, reflect.Float64, value)
}

// EmitObject satisfies the Encoder interface
func (fv *FieldValidator) EmitObject(key string, value interface{}) {
	fv.validateNextField(key, reflect.Interface, value)
}

// EmitLazyLogger satisfies the Encoder interface
func (fv *FieldValidator) EmitLazyLogger(key string, value LazyLogger) {
	fv.t.Error("Test infrastructure does not support EmitLazyLogger yet")
}

func (fv *FieldValidator) validateNextField(key string, actualKind reflect.Kind, value interface{}) {
	if fv.nextKey != key {
		fv.t.Errorf("Bad key: expected %q, found %q", fv.nextKey, key)
	}
	if fv.nextKind != actualKind {
		fv.t.Errorf("Bad reflect.Kind: expected %v, found %v", fv.nextKind, actualKind)
		return
	}
	if fv.nextValAsString != fmt.Sprint(value) {
		fv.t.Errorf("Bad value: expected %q, found %q", fv.nextValAsString, fmt.Sprint(value))
	}
	// All good.
}
