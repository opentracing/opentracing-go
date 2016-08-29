package opentracing

import "math"

const (
	stringType fieldType = iota
	boolType
	intType
	int64Type
	float64Type
	errorType
	objectType
	deferredObjectType
)

// LogField instances are constructed via LogBool, LogString, and so on.
// Tracing implementations may then handle them via the LogField.Process
// method.
type LogField struct {
	key          string
	valType      valType
	numericVal   int64
	stringVal    string
	interfaceVal interface{}
}

func LogString(key, val string) LogField {
	return LogField{
		key:       key,
		valType:   stringType,
		stringVal: val,
	}
}

func LogBool(key string, val bool) LogField {
	var numericVal int64
	if val {
		numericVal = 1
	}
	return LogField{
		key:        key,
		valType:    boolType,
		numericVal: numericVal,
	}
}

func LogInt(key string, val int) LogField {
	var numericVal int64 = int64(val)
	return LogField{
		key:        key,
		valType:    intType,
		numericVal: numericVal,
	}
}

func LogInt64(key string, val int64) LogField {
	return LogField{
		key:        key,
		valType:    int64Type,
		numericVal: val,
	}
}

func LogFloat64(key string, val float64) LogField {
	return LogField{
		key:        key,
		valType:    float64Type,
		numericVal: int64(math.Float64bits(val)),
	}
}

// REVIEWERS: etc etc for other numeric types if we like this direction

func LogError(err error) LogField {
	return LogField{
		key:          "error",
		valType:      errorType,
		interfaceVal: err,
	}
}

func LogObject(key string, obj interface{}) LogField {
	return LogField{
		key:          key,
		valType:      objectType,
		interfaceVal: obj,
	}
}

type DeferredObjectGenerator func() interface{}

func LogDeferredObject(key string, generator DeferredObjectGenerator) LogField {
	return LogField{
		key:          key,
		valType:      deferredObjectType,
		interfaceVal: generator,
	}
}

// LogFieldProcessor allows access to the contents of a LogField (via a call to
// LogField.Process).
//
// Tracer implementations typically provide an implementation of
// LogFieldProcessor; OpenTracing callers should not need to concern themselves
// with it.
type LogFieldProcessor interface {
	AddString(key, value string)
	AddBool(key string, value bool)
	AddInt(key string, value int)
	AddInt64(key string, value int64)
	AddFloat64(key string, value float64)
	AddObject(key string, value interface{})
}

// Process passes a LogField instance through to the appropriate type-specific
// method of a LogFieldProcessor.
func (lf LogField) Process(processor LogFieldProcessor) {
	switch lf.valType {
	case stringType:
		processor.AddString(lf.key, lf.stringVal)
	case boolType:
		processor.AddBool(lf.key, lf.numericVal != 0)
	case intType:
		processor.AddInt(lf.key, int(lf.numericVal))
	case int64Type:
		processor.AddInt64(lf.key, lf.numericVal)
	case float64Type:
		processor.AddFloat64(lf.key, math.Float64frombits(uint64(lf.numericVal)))
	case errorType:
		processor.AddString(lf.key, lf.obj.(error).Error())
	case objectType:
		processor.AddObject(lf.key, lf.interfaceVal)
	case deferredObjectType:
		processor.AddObject(lf.key, lf.interfaceVal.(DeferredObjectGenerator)())
	}
}
