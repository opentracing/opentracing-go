package main

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"sync"

	"github.com/opentracing/api-golang/opentracing"
)

// An implementation of opentracing.TraceContext.
type DapperishTraceContext struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID int64

	// A probabilistically unique identifier for a span.
	SpanID int64

	// Whether the trace is sampled.
	Sampled bool

	// XXX: comment
	tagLock   sync.RWMutex
	traceTags map[string]string
}

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
	return map[string]string{
		// NOTE: silly capitalization can be blamed on silly HTTP Header
		// conventions.
		"Traceid": strconv.FormatInt(dctx.TraceID, 10),
		"Spanid":  strconv.FormatInt(dctx.SpanID, 10),
		"Sampled": strconv.FormatBool(dctx.Sampled),
	}
}

func (d *DapperishTraceContextSource) UnmarshalTraceContextStringMap(
	encoded map[string]string,
) (opentracing.TraceContext, error) {
	traceID, err := strconv.ParseInt(encoded["Traceid"], 10, 64)
	if err != nil {
		return nil, err
	}
	spanID, err := strconv.ParseInt(encoded["Spanid"], 10, 64)
	if err != nil {
		return nil, err
	}
	sampled, err := strconv.ParseBool(encoded["Sampled"])
	if err != nil {
		return nil, err
	}
	// XXX: support tags
	return &DapperishTraceContext{
		TraceID:   traceID,
		SpanID:    spanID,
		Sampled:   sampled,
		traceTags: make(map[string]string),
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
