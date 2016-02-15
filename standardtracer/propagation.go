package standardtracer

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

type splitTextPropagator struct {
	tracer *tracerImpl
}
type splitBinaryPropagator struct {
	tracer *tracerImpl
}
type goHTTPPropagator struct {
	*splitBinaryPropagator
}

const (
	fieldNameTraceID = "traceid"
	fieldNameSpanID  = "spanid"
	fieldNameSampled = "sampled"
)

func (p *splitTextPropagator) InjectSpan(
	sp opentracing.Span,
	carrier interface{},
) error {
	sc := sp.(*spanImpl)
	splitTextCarrier, ok := carrier.(*opentracing.SplitTextCarrier)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}
	splitTextCarrier.TracerState = map[string]string{
		fieldNameTraceID: strconv.FormatInt(sc.raw.TraceID, 10),
		fieldNameSpanID:  strconv.FormatInt(sc.raw.SpanID, 10),
		fieldNameSampled: strconv.FormatBool(sc.raw.Sampled),
	}
	sc.Lock()
	if l := len(sc.raw.Attributes); l > 0 {
		splitTextCarrier.TraceAttributes = make(map[string]string, l)
		for k, v := range sc.raw.Attributes {
			splitTextCarrier.TraceAttributes[k] = v
		}
	}
	sc.Unlock()
	return nil
}

func (p *splitTextPropagator) JoinTrace(
	operationName string,
	carrier interface{},
) (opentracing.Span, error) {
	splitTextCarrier, ok := carrier.(*opentracing.SplitTextCarrier)
	if !ok {
		return nil, opentracing.ErrInvalidCarrier
	}
	requiredFieldCount := 0
	var traceID, propagatedSpanID int64
	var sampled bool
	var err error
	for k, v := range splitTextCarrier.TracerState {
		switch strings.ToLower(k) {
		case fieldNameTraceID:
			traceID, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, opentracing.ErrTraceCorrupted
			}
		case fieldNameSpanID:
			propagatedSpanID, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, opentracing.ErrTraceCorrupted
			}
		case fieldNameSampled:
			sampled, err = strconv.ParseBool(v)
			if err != nil {
				return nil, opentracing.ErrTraceCorrupted
			}
		default:
			continue
		}
		requiredFieldCount++
	}
	if requiredFieldCount < 3 {
		if len(splitTextCarrier.TracerState) == 0 {
			return nil, opentracing.ErrTraceNotFound
		}
		return nil, opentracing.ErrTraceCorrupted
	}

	sp := p.tracer.getSpan()
	sp.raw = RawSpan{
		StandardContext: StandardContext{
			TraceID:      traceID,
			SpanID:       randomID(),
			ParentSpanID: propagatedSpanID,
			Sampled:      sampled,
		},
	}
	sp.raw.Attributes = splitTextCarrier.TraceAttributes

	return p.tracer.startSpanInternal(
		sp,
		operationName,
		time.Now(),
		nil,
	), nil
}

func (p *splitBinaryPropagator) InjectSpan(
	sp opentracing.Span,
	carrier interface{},
) error {
	sc := sp.(*spanImpl)
	splitBinaryCarrier, ok := carrier.(*opentracing.SplitBinaryCarrier)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}
	var err error
	var sampledByte byte
	if sc.raw.Sampled {
		sampledByte = 1
	}

	// Handle the trace and span ids, and sampled status.
	contextBuf := bytes.NewBuffer(splitBinaryCarrier.TracerState[:0])
	err = binary.Write(contextBuf, binary.BigEndian, sc.raw.TraceID)
	if err != nil {
		return err
	}

	err = binary.Write(contextBuf, binary.BigEndian, sc.raw.SpanID)
	if err != nil {
		return err
	}

	err = binary.Write(contextBuf, binary.BigEndian, sampledByte)
	if err != nil {
		return err
	}

	// Handle the attributes.
	attrsBuf := bytes.NewBuffer(splitBinaryCarrier.TraceAttributes[:0])
	err = binary.Write(attrsBuf, binary.BigEndian, int32(len(sc.raw.Attributes)))
	if err != nil {
		return err
	}
	for k, v := range sc.raw.Attributes {
		if err = binary.Write(attrsBuf, binary.BigEndian, int32(len(k))); err != nil {
			return err
		}
		attrsBuf.WriteString(k)
		if err = binary.Write(attrsBuf, binary.BigEndian, int32(len(v))); err != nil {
			return err
		}
		attrsBuf.WriteString(v)
	}

	splitBinaryCarrier.TracerState = contextBuf.Bytes()
	splitBinaryCarrier.TraceAttributes = attrsBuf.Bytes()
	return nil
}

func (p *splitBinaryPropagator) JoinTrace(
	operationName string,
	carrier interface{},
) (opentracing.Span, error) {
	var err error
	splitBinaryCarrier, ok := carrier.(*opentracing.SplitBinaryCarrier)
	if !ok {
		return nil, opentracing.ErrInvalidCarrier
	}
	if len(splitBinaryCarrier.TracerState) == 0 {
		return nil, opentracing.ErrTraceNotFound
	}
	// Handle the trace, span ids, and sampled status.
	contextReader := bytes.NewReader(splitBinaryCarrier.TracerState)
	var traceID, propagatedSpanID int64
	var sampledByte byte

	err = binary.Read(contextReader, binary.BigEndian, &traceID)
	if err != nil {
		return nil, opentracing.ErrTraceCorrupted
	}
	err = binary.Read(contextReader, binary.BigEndian, &propagatedSpanID)
	if err != nil {
		return nil, opentracing.ErrTraceCorrupted
	}
	err = binary.Read(contextReader, binary.BigEndian, &sampledByte)
	if err != nil {
		return nil, opentracing.ErrTraceCorrupted
	}

	// Handle the attributes.
	attrsReader := bytes.NewReader(splitBinaryCarrier.TraceAttributes)
	var numAttrs int32
	err = binary.Read(attrsReader, binary.BigEndian, &numAttrs)
	if err != nil {
		return nil, opentracing.ErrTraceCorrupted
	}
	iNumAttrs := int(numAttrs)
	var attrMap map[string]string
	if iNumAttrs > 0 {
		var buf bytes.Buffer // TODO(tschottdorf): candidate for sync.Pool
		attrMap = make(map[string]string, iNumAttrs)
		var keyLen, valLen int32
		for i := 0; i < iNumAttrs; i++ {
			err = binary.Read(attrsReader, binary.BigEndian, &keyLen)
			if err != nil {
				return nil, opentracing.ErrTraceCorrupted
			}
			buf.Grow(int(keyLen))
			if n, err := io.CopyN(&buf, attrsReader, int64(keyLen)); err != nil || int32(n) != keyLen {
				return nil, opentracing.ErrTraceCorrupted
			}
			key := buf.String()
			buf.Reset()

			err = binary.Read(attrsReader, binary.BigEndian, &valLen)
			if err != nil {
				return nil, opentracing.ErrTraceCorrupted
			}
			if n, err := io.CopyN(&buf, attrsReader, int64(valLen)); err != nil || int32(n) != valLen {
				return nil, opentracing.ErrTraceCorrupted
			}
			attrMap[key] = buf.String()
			buf.Reset()
		}
	}

	sp := p.tracer.getSpan()
	sp.raw = RawSpan{
		StandardContext: StandardContext{
			TraceID:      traceID,
			SpanID:       randomID(),
			ParentSpanID: propagatedSpanID,
			Sampled:      sampledByte != 0,
		},
	}
	sp.raw.Attributes = attrMap

	return p.tracer.startSpanInternal(
		sp,
		operationName,
		time.Now(),
		nil,
	), nil
}

const (
	tracerStateHeaderName = "Tracer-State"
	traceAttrsHeaderName  = "Trace-Attributes"
)

func (p *goHTTPPropagator) InjectSpan(
	sp opentracing.Span,
	carrier interface{},
) error {
	// Defer to SplitBinary for the real work.
	splitBinaryCarrier := opentracing.NewSplitBinaryCarrier()
	if err := p.splitBinaryPropagator.InjectSpan(sp, splitBinaryCarrier); err != nil {
		return err
	}

	// Encode into the HTTP header as two base64 strings.
	header := carrier.(http.Header)
	header.Add(tracerStateHeaderName, base64.StdEncoding.EncodeToString(
		splitBinaryCarrier.TracerState))
	header.Add(traceAttrsHeaderName, base64.StdEncoding.EncodeToString(
		splitBinaryCarrier.TraceAttributes))

	return nil
}

func (p *goHTTPPropagator) JoinTrace(
	operationName string,
	carrier interface{},
) (opentracing.Span, error) {
	// Decode the two base64-encoded data blobs from the HTTP header.
	header := carrier.(http.Header)
	tracerStateBase64, found := header[http.CanonicalHeaderKey(tracerStateHeaderName)]
	if !found || len(tracerStateBase64) == 0 {
		return nil, opentracing.ErrTraceNotFound
	}
	traceAttrsBase64, found := header[http.CanonicalHeaderKey(traceAttrsHeaderName)]
	if !found || len(traceAttrsBase64) == 0 {
		return nil, opentracing.ErrTraceNotFound
	}
	tracerStateBinary, err := base64.StdEncoding.DecodeString(tracerStateBase64[0])
	if err != nil {
		return nil, opentracing.ErrTraceCorrupted
	}
	traceAttrsBinary, err := base64.StdEncoding.DecodeString(traceAttrsBase64[0])
	if err != nil {
		return nil, opentracing.ErrTraceCorrupted
	}

	// Defer to SplitBinary for the real work.
	splitBinaryCarrier := &opentracing.SplitBinaryCarrier{
		TracerState:     tracerStateBinary,
		TraceAttributes: traceAttrsBinary,
	}
	return p.splitBinaryPropagator.JoinTrace(operationName, splitBinaryCarrier)
}
