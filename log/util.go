package log

import (
	"fmt"
	"math"

	"go.uber.org/zap/zapcore"
)

// InterleavedKVToFields converts keyValues a la Span.LogKV() to a Field slice
// a la Span.LogFields().
func InterleavedKVToFields(keyValues ...interface{}) ([]Field, error) {
	if len(keyValues)%2 != 0 {
		return nil, fmt.Errorf("non-even keyValues len: %d", len(keyValues))
	}
	fields := make([]Field, len(keyValues)/2)
	for i := 0; i*2 < len(keyValues); i++ {
		key, ok := keyValues[i*2].(string)
		if !ok {
			return nil, fmt.Errorf(
				"non-string key (pair #%d): %T",
				i, keyValues[i*2])
		}
		switch typedVal := keyValues[i*2+1].(type) {
		case bool:
			fields[i] = Bool(key, typedVal)
		case string:
			fields[i] = String(key, typedVal)
		case int:
			fields[i] = Int(key, typedVal)
		case int8:
			fields[i] = Int32(key, int32(typedVal))
		case int16:
			fields[i] = Int32(key, int32(typedVal))
		case int32:
			fields[i] = Int32(key, typedVal)
		case int64:
			fields[i] = Int64(key, typedVal)
		case uint:
			fields[i] = Uint64(key, uint64(typedVal))
		case uint64:
			fields[i] = Uint64(key, typedVal)
		case uint8:
			fields[i] = Uint32(key, uint32(typedVal))
		case uint16:
			fields[i] = Uint32(key, uint32(typedVal))
		case uint32:
			fields[i] = Uint32(key, typedVal)
		case float32:
			fields[i] = Float32(key, typedVal)
		case float64:
			fields[i] = Float64(key, typedVal)
		default:
			// When in doubt, coerce to a string
			fields[i] = String(key, fmt.Sprint(typedVal))
		}
	}
	return fields, nil
}

// FieldFromZap returns a new standard Opentracing Field that contains the same information as the
// standard Zap field given as input.
func FieldFromZap(zapField zapcore.Field) Field {
	switch zapField.Type {

	case zapcore.BoolType:
		val := false
		if zapField.Integer >= 1 {
			val = true
		}
		return Bool(zapField.Key, val)
	case zapcore.Float32Type:
		return Float32(zapField.Key, math.Float32frombits(uint32(zapField.Integer)))
	case zapcore.Float64Type:
		return Float64(zapField.Key, math.Float64frombits(uint64(zapField.Integer)))
	case zapcore.Int64Type:
		return Int64(zapField.Key, int64(zapField.Integer))
	case zapcore.Int32Type:
		return Int32(zapField.Key, int32(zapField.Integer))
	case zapcore.StringType:
		return String(zapField.Key, zapField.String)
	case zapcore.Uint64Type:
		return Uint64(zapField.Key, uint64(zapField.Integer))
	case zapcore.Uint32Type:
		return Uint32(zapField.Key, uint32(zapField.Integer))
	case zapcore.ErrorType:
		return Error(zapField.Interface.(error))
	default:
		// By default, use the generic "object" type.
		return Object(zapField.Key, zapField.Interface)
	}
}

// FieldsFromZap returns a slice of Opentracing Fields that contains the same information as a
// standard Zap field.
func FieldsFromZap(zapFields ...zapcore.Field) []Field {
	opentracingFields := make([]Field, len(zapFields))

	for i, zapField := range zapFields {
		opentracingFields[i] = FieldFromZap(zapField)
	}

	return opentracingFields
}
