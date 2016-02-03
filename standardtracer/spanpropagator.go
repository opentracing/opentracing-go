package standardtracer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/opentracing/opentracing-go"
)

const (
	fieldNameTraceID = "traceid"
	fieldNameSpanID  = "spanid"
	fieldNameSampled = "sampled"
)

func (s *tracerImpl) PropagateSpanAsText(
	sp opentracing.Span,
) (
	contextIDMap map[string]string,
	attrsMap map[string]string,
) {
	sc := sp.(*spanImpl).raw.StandardContext
	contextIDMap = map[string]string{
		fieldNameTraceID: strconv.FormatInt(sc.TraceID, 10),
		fieldNameSpanID:  strconv.FormatInt(sc.SpanID, 10),
		fieldNameSampled: strconv.FormatBool(sc.Sampled),
	}
	sc.tagLock.RLock()
	attrsMap = make(map[string]string, len(sc.traceAttrs))
	for k, v := range sc.traceAttrs {
		attrsMap[k] = v
	}
	sc.tagLock.RUnlock()
	return contextIDMap, attrsMap
}

func (s *tracerImpl) PropagateSpanAsBinary(
	sp opentracing.Span,
) (
	traceContextID []byte,
	traceAttrs []byte,
) {
	sc := sp.(*spanImpl).raw.StandardContext
	var err error

	// Handle the trace and span ids.
	contextBuf := new(bytes.Buffer)
	err = binary.Write(contextBuf, binary.BigEndian, sc.TraceID)
	if err != nil {
		panic(err)
	}
	err = binary.Write(contextBuf, binary.BigEndian, sc.SpanID)
	if err != nil {
		panic(err)
	}

	// Handle the attributes.
	attrsBuf := new(bytes.Buffer)
	err = binary.Write(attrsBuf, binary.BigEndian, len(sc.traceAttrs))
	if err != nil {
		panic(err)
	}
	for k, v := range sc.traceAttrs {
		keyBytes := []byte(k)
		err = binary.Write(attrsBuf, binary.BigEndian, len(keyBytes))
		err = binary.Write(attrsBuf, binary.BigEndian, keyBytes)
		valBytes := []byte(v)
		err = binary.Write(attrsBuf, binary.BigEndian, len(valBytes))
		err = binary.Write(attrsBuf, binary.BigEndian, valBytes)
	}

	return contextBuf.Bytes(), attrsBuf.Bytes()
}

func (s *tracerImpl) JoinTraceFromBinary(
	operationName string,
	traceContextID []byte,
	traceAttrs []byte,
) (opentracing.Span, error) {
	var err error
	// Handle the trace and span ids.
	contextReader := bytes.NewReader(traceContextID)
	var traceID, propagatedSpanID int64
	var sampledByte byte
	err = binary.Read(contextReader, binary.BigEndian, &traceID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(contextReader, binary.BigEndian, &propagatedSpanID)
	if err != nil {
		return nil, err
	}

	// Handle the attributes.
	attrsReader := bytes.NewReader(traceAttrs)
	var numAttrs int
	err = binary.Read(attrsReader, binary.BigEndian, &numAttrs)
	if err != nil {
		return nil, err
	}
	attrMap := make(map[string]string, numAttrs)
	for i := 0; i < numAttrs; i++ {
		var keyLen int
		err = binary.Read(attrsReader, binary.BigEndian, &keyLen)
		if err != nil {
			return nil, err
		}
		keyBytes := make([]byte, keyLen)
		err = binary.Read(attrsReader, binary.BigEndian, &keyBytes)
		if err != nil {
			return nil, err
		}

		var valLen int
		err = binary.Read(attrsReader, binary.BigEndian, &valLen)
		if err != nil {
			return nil, err
		}
		valBytes := make([]byte, valLen)
		err = binary.Read(attrsReader, binary.BigEndian, &valBytes)
		if err != nil {
			return nil, err
		}

		attrMap[string(keyBytes)] = string(valBytes)
	}

	return s.startSpanGeneric(
			operationName,
			&StandardContext{
				TraceID:      traceID,
				SpanID:       randomID(),
				ParentSpanID: propagatedSpanID,
				Sampled:      sampledByte != 0,
				traceAttrs:   make(map[string]string),
			}),
		nil
}

func (s *tracerImpl) JoinTraceFromText(
	operationName string,
	contextIDMap map[string]string,
	attrsMap map[string]string,
) (opentracing.Span, error) {
	requiredFieldCount := 0
	var traceID, spanID int64
	var sampled bool
	var err error
	for k, v := range contextIDMap {
		switch strings.ToLower(k) {
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

	lcAttrsMap := make(map[string]string, len(attrsMap))
	for k, v := range attrsMap {
		lcAttrsMap[strings.ToLower(k)] = v
	}

	return s.startSpanGeneric(
			operationName,
			&StandardContext{
				TraceID:    traceID,
				SpanID:     spanID,
				Sampled:    sampled,
				traceAttrs: lcAttrsMap,
			}),
		nil
}
