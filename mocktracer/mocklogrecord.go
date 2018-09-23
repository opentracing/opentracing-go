package mocktracer

import (
	"fmt"
	"reflect"
	"time"

	"github.com/opentracing/opentracing-go/log"
)

// MockLogRecord represents data logged to a Span via Span.LogFields or
// Span.LogKV.
type MockLogRecord struct {
	Timestamp time.Time
	Fields    []MockKeyValue
}

// MockKeyValue represents a single key:value pair.
type MockKeyValue struct {
	Key string

	// All MockLogRecord values are coerced to strings via fmt.Sprint(), though
	// we retain their type separately.
	ValueKind   reflect.Kind
	ValueString string
}

// EmitString belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitString(key, value string) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitBool belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitBool(key string, value bool) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitInt belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitInt(key string, value int) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitInt32 belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitInt32(key string, value int32) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitInt64 belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitInt64(key string, value int64) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitUint32 belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitUint32(key string, value uint32) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitUint64 belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitUint64(key string, value uint64) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitFloat32 belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitFloat32(key string, value float32) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitFloat64 belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitFloat64(key string, value float64) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitObject belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitObject(key string, value interface{}) {
	lr.Fields = append(lr.Fields, MockKeyValue{
		Key:         key,
		ValueKind:   reflect.TypeOf(value).Kind(),
		ValueString: fmt.Sprint(value),
	})
}

// EmitLazyLogger belongs to the log.Encoder interface
func (lr *MockLogRecord) EmitLazyLogger(value log.LazyLogger) {
	value(lr)
}
