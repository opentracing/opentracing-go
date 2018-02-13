package log

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestFieldString(t *testing.T) {
	testCases := []struct {
		field    Field
		expected string
	}{
		{
			field:    String("key", "value"),
			expected: "key:value",
		},
		{
			field:    Bool("key", true),
			expected: "key:true",
		},
		{
			field:    Int("key", 5),
			expected: "key:5",
		},
		{
			field:    Error(fmt.Errorf("err msg")),
			expected: "error:err msg",
		},
		{
			field:    Error(nil),
			expected: "error:<nil>",
		},
		{
			field:    Noop(),
			expected: ":<nil>",
		},
	}
	for i, tc := range testCases {
		if str := tc.field.String(); str != tc.expected {
			t.Errorf("%d: expected '%s', got '%s'", i, tc.expected, str)
		}
	}
}

func TestNoopDoesNotMarshal(t *testing.T) {
	mockEncoder := struct {
		Encoder
	}{}
	f := Noop()
	f.Marshal(mockEncoder) // panics if any Encoder method is invoked
}

func TestFieldsFromZap(t *testing.T) {
	testCases := []struct {
		zapField      zapcore.Field
		expectedField Field
	}{
		{
			zap.String("namespace", "123"),
			String("namespace", "123"),
		},
		{
			zap.String("namespace", ""),
			String("namespace", ""),
		},
		{
			zap.String("", "123"),
			String("", "123"),
		},
		{
			zap.Int("namespace", 1),
			Int64("namespace", 1),
		},
		{
			zap.Int32("namespace", 1),
			Int32("namespace", 1),
		},
		{
			zap.Int64("namespace", 1),
			Int64("namespace", 1),
		},
		{
			zap.Uint32("namespace", 1),
			Uint32("namespace", 1),
		},
		{
			zap.Uint64("namespace", 1),
			Uint64("namespace", 1),
		},
		{
			zap.Float32("namespace", 1),
			Float32("namespace", 1),
		},
		{
			zap.Float64("namespace", 1),
			Float64("namespace", 1),
		},
		{
			zap.Bool("namespace", false),
			Bool("namespace", false),
		},
		{
			zap.Bool("namespace", true),
			Bool("namespace", true),
		},
	}

	for i, data := range testCases {
		result := FieldFromZap(data.zapField)
		if result.Key() != data.expectedField.Key() {
			t.Errorf("%d: expected same key. Got %s but expected %s", i, result.Key(), data.expectedField.Key())
		}
		if result.String() != data.expectedField.String() {
			t.Errorf("%d: expected same string. Got %s but expected %s", i, result.String(), data.expectedField.String())
		}
		if result.Value() != data.expectedField.Value() {
			t.Errorf("%d: expected same value. Got %s but expected %s", i, result.Value(), data.expectedField.Value())
		}
	}
}
