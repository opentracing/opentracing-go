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

// AddString belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddString(key, value string) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddBool belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddBool(key string, value bool) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddInt belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddInt(key string, value int) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddInt32 belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddInt32(key string, value int32) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddInt64 belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddInt64(key string, value int64) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddUint32 belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddUint32(key string, value uint32) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddUint64 belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddUint64(key string, value uint64) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddFloat32 belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddFloat32(key string, value float32) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddFloat64 belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddFloat64(key string, value float64) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddObject belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddObject(key string, value interface{}) {
	m.Key = key
	m.ValueKind = reflect.TypeOf(value).Kind()
	m.ValueString = fmt.Sprint(value)
}

// AddLazyLogger belongs to the log.FieldVisitor interface
func (m *MockKeyValue) AddLazyLogger(key string, value log.LazyLogger) {
	var meta MockKeyValue
	value(&meta)
	m.Key = meta.Key
	m.ValueKind = meta.ValueKind
	m.ValueString = meta.ValueString
}
