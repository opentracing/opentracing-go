package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/opentracing/api-golang/opentracing"
	"github.com/opentracing/api-golang/opentracing/standardtracer"
)

func NewDapperishTracer(processName string) opentracing.OpenTracer {
	return standardtracer.New(
		NewTrivialRecorder(processName),
		NewDapperishTraceContextSource())
}

// An implementation of opentracing.TraceContext.
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
	fieldNameTraceId   = "Traceid"
	fieldNameSpanId    = "Spanid"
	fieldNameSampled   = "Sampled"
	fieldNameTagPrefix = "Tag-"
)

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

// An implementation of opentracing.TraceContextSource.
type DapperishTraceContextSource struct{}

func NewDapperishTraceContextSource() *DapperishTraceContextSource {
	return &DapperishTraceContextSource{}
}

func (m *DapperishTraceContextSource) NewRootTraceContext() opentracing.TraceContext {
	return &DapperishTraceContext{
		TraceID:   randomID(),
		SpanID:    randomID(),
		Sampled:   randomID()%1024 == 0,
		traceTags: make(map[string]string),
	}
}

func (d *DapperishTraceContext) SetTraceTag(key, val string) opentracing.TraceContext {
	d.tagLock.Lock()
	defer d.tagLock.Unlock()

	d.traceTags[key] = val
	return d
}

func (d *DapperishTraceContext) TraceTag(key string) string {
	d.tagLock.RLock()
	defer d.tagLock.RUnlock()

	return d.traceTags[key]
}

func (d *DapperishTraceContextSource) MarshalTraceContextStringMap(
	ctx opentracing.TraceContext,
) map[string]string {
	dctx := ctx.(*DapperishTraceContext)
	rval := map[string]string{
		fieldNameTraceId: strconv.FormatInt(dctx.TraceID, 10),
		fieldNameSpanId:  strconv.FormatInt(dctx.SpanID, 10),
		fieldNameSampled: strconv.FormatBool(dctx.Sampled),
	}
	dctx.tagLock.RLock()
	for k, v := range dctx.traceTags {
		rval[fieldNameTagPrefix+k] = v
	}
	dctx.tagLock.RUnlock()
	return rval
}

func (d *DapperishTraceContextSource) UnmarshalTraceContextStringMap(
	encoded map[string]string,
) (opentracing.TraceContext, error) {
	traceTags := make(map[string]string)
	requiredFieldCount := 0
	var traceID, spanID int64
	var sampled bool
	var err error
	for k, v := range encoded {
		switch k {
		case fieldNameTraceId:
			traceID, err = strconv.ParseInt(encoded[fieldNameTraceId], 10, 64)
			if err != nil {
				return nil, err
			}
			requiredFieldCount++
		case fieldNameSpanId:
			spanID, err = strconv.ParseInt(encoded[fieldNameSpanId], 10, 64)
			if err != nil {
				return nil, err
			}
			requiredFieldCount++
		case fieldNameSampled:
			sampled, err = strconv.ParseBool(encoded[fieldNameSampled])
			if err != nil {
				return nil, err
			}
			requiredFieldCount++
		default:
			if strings.HasPrefix(k, fieldNameTagPrefix) {
				traceTags[strings.TrimPrefix(k, fieldNameTagPrefix)] = v
			} else {
				return nil, fmt.Errorf("Unknown string map field: %v", k)
			}
		}
	}
	if requiredFieldCount < 3 {
		return nil, fmt.Errorf("Only found %v of 3 required fields", requiredFieldCount)
	}

	return &DapperishTraceContext{
		TraceID:   traceID,
		SpanID:    spanID,
		Sampled:   sampled,
		traceTags: traceTags,
	}, nil
}

func (d *DapperishTraceContextSource) MarshalTraceContextBinary(ctx opentracing.TraceContext) []byte {
	dtc := ctx.(*DapperishTraceContext)
	// XXX: support tags
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
	var sampledByte byte = 0
	if dtc.Sampled {
		sampledByte = 1
	}
	err = binary.Write(buf, binary.BigEndian, sampledByte)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (d *DapperishTraceContextSource) UnmarshalTraceContextBinary(
	encoded []byte,
) (opentracing.TraceContext, error) {
	// XXX: support tags
	var err error
	reader := bytes.NewReader(encoded)
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
	return &DapperishTraceContext{
		TraceID:   traceID,
		SpanID:    spanID,
		Sampled:   sampledByte != 0,
		traceTags: make(map[string]string),
	}, nil
}
