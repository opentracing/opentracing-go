package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"

	"github.com/opentracing/api-golang/opentracing"
)

// An implementation of opentracing.TraceContextID.
type DapperishTraceContextID struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID int64

	// A probabilistically unique identifier for a span.
	SpanID int64

	// Whether the trace is sampled.
	Sampled bool
}

func (d *DapperishTraceContextID) NewChild() (opentracing.TraceContextID, opentracing.Tags) {
	return &DapperishTraceContextID{
		TraceID: d.TraceID,
		SpanID:  randomID(),
		Sampled: d.Sampled,
	}, opentracing.Tags{"parent_span_id": d.SpanID}
}

func (d *DapperishTraceContextID) SerializeBinary() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, d.TraceID)
	if err != nil {
		panic(err)
	}
	err = binary.Write(buf, binary.BigEndian, d.SpanID)
	if err != nil {
		panic(err)
	}
	var sampledByte byte = 0
	if d.Sampled {
		sampledByte = 1
	}
	err = binary.Write(buf, binary.BigEndian, sampledByte)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (d *DapperishTraceContextID) SerializeASCII() string {
	return base64.StdEncoding.EncodeToString(d.SerializeBinary())
}

// An implementation of opentracing.TraceContextIDSource.
type DapperishTraceContextIDSource struct{}

func NewDapperishTraceContextIDSource() *DapperishTraceContextIDSource {
	return &DapperishTraceContextIDSource{}
}

func (m *DapperishTraceContextIDSource) NewRootTraceContextID() opentracing.TraceContextID {
	return &DapperishTraceContextID{
		TraceID: randomID(),
		SpanID:  randomID(),
		Sampled: randomID()%1024 == 0,
	}
}

func (m *DapperishTraceContextIDSource) DeserializeBinaryTraceContextID(encoded []byte) (opentracing.TraceContextID, error) {
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
	return &DapperishTraceContextID{
		TraceID: traceID,
		SpanID:  spanID,
		Sampled: sampledByte != 0,
	}, nil
}

func (d *DapperishTraceContextIDSource) DeserializeASCIITraceContextID(encoded string) (opentracing.TraceContextID, error) {
	ctxIdBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return d.DeserializeBinaryTraceContextID(ctxIdBytes)
}
