package log

import "math"

type fieldType int

const (
	stringType fieldType = iota
	boolType
	intType
	int32Type
	uint32Type
	int64Type
	uint64Type
	float32Type
	float64Type
	errorType
	objectType
	lazyLoggerType
)

// Field instances are constructed via LogBool, LogString, and so on.
// Tracing implementations may then handle them via the Field.Process
// method.
//
// "heavily influenced by" (i.e., partially stolen from)
// https://github.com/uber-go/zap
type Field struct {
	key          string
	fieldType    fieldType
	numericVal   int64
	stringVal    string
	interfaceVal interface{}
}

func String(key, val string) Field {
	return Field{
		key:       key,
		fieldType: stringType,
		stringVal: val,
	}
}

func Bool(key string, val bool) Field {
	var numericVal int64
	if val {
		numericVal = 1
	}
	return Field{
		key:        key,
		fieldType:  boolType,
		numericVal: numericVal,
	}
}

func Int(key string, val int) Field {
	return Field{
		key:        key,
		fieldType:  intType,
		numericVal: int64(val),
	}
}

func Int32(key string, val int32) Field {
	return Field{
		key:        key,
		fieldType:  int32Type,
		numericVal: int64(val),
	}
}

func Int64(key string, val int64) Field {
	return Field{
		key:        key,
		fieldType:  int64Type,
		numericVal: val,
	}
}

func Uint32(key string, val uint32) Field {
	return Field{
		key:        key,
		fieldType:  uint32Type,
		numericVal: int64(val),
	}
}

func Uint64(key string, val uint64) Field {
	return Field{
		key:        key,
		fieldType:  int64Type,
		numericVal: int64(val),
	}
}

func Float32(key string, val float32) Field {
	return Field{
		key:        key,
		fieldType:  float32Type,
		numericVal: int64(math.Float32bits(val)),
	}
}

func Float64(key string, val float64) Field {
	return Field{
		key:        key,
		fieldType:  float64Type,
		numericVal: int64(math.Float64bits(val)),
	}
}

func Error(err error) Field {
	return Field{
		key:          "error",
		fieldType:    errorType,
		interfaceVal: err,
	}
}

func Object(key string, obj interface{}) Field {
	return Field{
		key:          key,
		fieldType:    objectType,
		interfaceVal: obj,
	}
}

type LazyLogger func(fv FieldVisitor)

func Lazy(key string, ll LazyLogger) Field {
	return Field{
		key:          key,
		fieldType:    lazyLoggerType,
		interfaceVal: ll,
	}
}

// FieldVisitor allows access to the contents of a Field (via a call to
// Field.Visit).
//
// Tracer implementations typically provide an implementation of
// FieldVisitor; OpenTracing callers should not need to concern themselves
// with it.
type FieldVisitor interface {
	AddString(key, value string)
	AddBool(key string, value bool)
	AddInt(key string, value int)
	AddInt32(key string, value int32)
	AddInt64(key string, value int64)
	AddUint32(key string, value uint32)
	AddUint64(key string, value uint64)
	AddFloat32(key string, value float32)
	AddFloat64(key string, value float64)
	AddObject(key string, value interface{})
	AddLazyLogger(key string, value LazyLogger)
}

// Visit passes a Field instance through to the appropriate field-type-specific
// method of a FieldVisitor.
func (lf Field) Visit(visitor FieldVisitor) {
	switch lf.fieldType {
	case stringType:
		visitor.AddString(lf.key, lf.stringVal)
	case boolType:
		visitor.AddBool(lf.key, lf.numericVal != 0)
	case intType:
		visitor.AddInt(lf.key, int(lf.numericVal))
	case int32Type:
		visitor.AddInt32(lf.key, int32(lf.numericVal))
	case int64Type:
		visitor.AddInt64(lf.key, int64(lf.numericVal))
	case uint32Type:
		visitor.AddUint32(lf.key, uint32(lf.numericVal))
	case uint64Type:
		visitor.AddUint64(lf.key, uint64(lf.numericVal))
	case float32Type:
		visitor.AddFloat32(lf.key, math.Float32frombits(uint32(lf.numericVal)))
	case float64Type:
		visitor.AddFloat64(lf.key, math.Float64frombits(uint64(lf.numericVal)))
	case errorType:
		visitor.AddString(lf.key, lf.interfaceVal.(error).Error())
	case objectType:
		visitor.AddObject(lf.key, lf.interfaceVal)
	case lazyLoggerType:
		visitor.AddLazyLogger(lf.key, lf.interfaceVal.(LazyLogger))
	}
}
