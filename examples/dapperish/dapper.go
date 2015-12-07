package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"

	"github.com/opentracing/api-golang/opentracing"
	"github.com/opentracing/api-golang/opentracing/standardtracer"
)

// NewTracer returns a new dapperish OpenTracer instance.
func NewTracer(processName string) opentracing.OpenTracer {
	return standardtracer.New(
		NewTrivialRecorder(processName),
		NewTraceContextSource())
}

// DapperishTraceContext is an implementation of opentracing.TraceContext.
type DapperishTraceContext struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID int64

	// A probabilistically unique identifier for a span.
	SpanID int64

	// Whether the trace is sampled.
	Sampled bool

	// `tagLock` protects the `traceTags` map, which in turn supports
	// `SetTraceTag` and `TraceTag`.
	tagLock   sync.RWMutex
	traceTags map[string]string
}

const (
	// Note that these strings are designed to be unchanged by the conversion
	// into standard HTTP headers (which messes with capitalization).
	fieldNameTraceID = "Traceid"
	fieldNameSpanID  = "Spanid"
	fieldNameSampled = "Sampled"
)

// NewChild complies with the opentracing.TraceContext interface.
func (d *DapperishTraceContext) NewChild() (opentracing.TraceContext, opentracing.Tags) {
	d.tagLock.RLock()
	newTags := make(map[string]string, len(d.traceTags))
	for k, v := range d.traceTags {
		newTags[k] = v
	}
	d.tagLock.RUnlock()

	return &DapperishTraceContext{
		TraceID:   d.TraceID,
		SpanID:    randomID(),
		Sampled:   d.Sampled,
		traceTags: newTags,
	}, opentracing.Tags{"parent_span_id": d.SpanID}
}

// SetTraceTag complies with the opentracing.TraceContext interface.
func (d *DapperishTraceContext) SetTraceTag(key, val string) opentracing.TraceContext {
	d.tagLock.Lock()
	defer d.tagLock.Unlock()

	d.traceTags[key] = val
	return d
}

// TraceTag complies with the opentracing.TraceContext interface.
func (d *DapperishTraceContext) TraceTag(key string) string {
	d.tagLock.RLock()
	defer d.tagLock.RUnlock()

	return d.traceTags[key]
}

// DapperishTraceContextSource is an implementation of
// opentracing.TraceContextSource.
type DapperishTraceContextSource struct{}

// NewTraceContextSource returns a dapperish opentracing.TraceContextSource
// implementation.
func NewTraceContextSource() *DapperishTraceContextSource {
	return &DapperishTraceContextSource{}
}

// NewRootTraceContext complies with the opentracing.TraceContextSource interface.
func (d *DapperishTraceContextSource) NewRootTraceContext() opentracing.TraceContext {
	return &DapperishTraceContext{
		TraceID:   randomID(),
		SpanID:    randomID(),
		Sampled:   randomID()%1024 == 0,
		traceTags: make(map[string]string),
	}
}

// MarshalTraceContextStringMap complies with the
// opentracing.TraceContextSource interface.
func (d *DapperishTraceContextSource) MarshalTraceContextStringMap(
	ctx opentracing.TraceContext,
) (contextIDMap map[string]string, tagsMap map[string]string) {
	dctx := ctx.(*DapperishTraceContext)
	contextIDMap = map[string]string{
		fieldNameTraceID: strconv.FormatInt(dctx.TraceID, 10),
		fieldNameSpanID:  strconv.FormatInt(dctx.SpanID, 10),
		fieldNameSampled: strconv.FormatBool(dctx.Sampled),
	}
	dctx.tagLock.RLock()
	tagsMap = make(map[string]string, len(dctx.traceTags))
	for k, v := range dctx.traceTags {
		tagsMap[k] = v
	}
	dctx.tagLock.RUnlock()
	return contextIDMap, tagsMap
}

// UnmarshalTraceContextStringMap complies with the
// opentracing.TraceContextSource interface.
func (d *DapperishTraceContextSource) UnmarshalTraceContextStringMap(
	contextIDMap map[string]string,
	tagsMap map[string]string,
) (opentracing.TraceContext, error) {
	requiredFieldCount := 0
	var traceID, spanID int64
	var sampled bool
	var err error
	for k, v := range contextIDMap {
		switch k {
		case fieldNameTraceID:
			traceID, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			requiredFieldCount++
		case fieldNameSpanID:
			spanID, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			requiredFieldCount++
		case fieldNameSampled:
			sampled, err = strconv.ParseBool(v)
			if err != nil {
				return nil, err
			}
			requiredFieldCount++
		default:
			return nil, fmt.Errorf("Unknown contextIDMap field: %v", k)
		}
	}
	if requiredFieldCount < 3 {
		return nil, fmt.Errorf("Only found %v of 3 required fields", requiredFieldCount)
	}

	return &DapperishTraceContext{
		TraceID:   traceID,
		SpanID:    spanID,
		Sampled:   sampled,
		traceTags: tagsMap,
	}, nil
}

// MarshalTraceContextBinary complies with the opentracing.TraceContextSource
// interface.
func (d *DapperishTraceContextSource) MarshalTraceContextBinary(ctx opentracing.TraceContext) (contextID []byte, traceTags []byte) {
	dtc := ctx.(*DapperishTraceContext)
	var err error
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, dtc.TraceID)
	if err != nil {
		panic(err)
	}
	err = binary.Write(buf, binary.BigEndian, dtc.SpanID)
	if err != nil {
		panic(err)
	}
	var sampledByte byte
	if dtc.Sampled {
		sampledByte = 1
	}
	err = binary.Write(buf, binary.BigEndian, sampledByte)
	if err != nil {
		panic(err)
	}
	// XXX: support tags
	return buf.Bytes(), []byte{}
}

// UnmarshalTraceContextBinary complies with the opentracing.TraceContextSource
// interface.
func (d *DapperishTraceContextSource) UnmarshalTraceContextBinary(
	contextID []byte,
	traceTags []byte,
) (opentracing.TraceContext, error) {
	var err error
	reader := bytes.NewReader(contextID)
	var traceID, spanID int64
	var sampledByte byte

	err = binary.Read(reader, binary.BigEndian, &traceID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &spanID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &sampledByte)
	if err != nil {
		return nil, err
	}
	// XXX: support tags
	return &DapperishTraceContext{
		TraceID:   traceID,
		SpanID:    spanID,
		Sampled:   sampledByte != 0,
		traceTags: make(map[string]string),
	}, nil
}
