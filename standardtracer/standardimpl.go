package standardtracer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"

	"golang.org/x/net/context"
)

// New creates and returns a standard Tracer which defers to `recorder` and
// `source` as appropriate.
func New(recorder Recorder) opentracing.Tracer {
	return &standardTracer{
		recorder: recorder,
	}
}

type StandardContext struct {
	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID int64

	// A probabilistically unique identifier for a span.
	SpanID int64

	// Whether the trace is sampled.
	Sampled bool

	// `tagLock` protects the `traceAttrs` map, which in turn supports
	// `SetTraceAttribute` and `TraceAttribute`.
	tagLock    sync.RWMutex
	traceAttrs map[string]string
}

func NewRootStandardContext() *StandardContext {
	return &StandardContext{
		TraceID:    randomID(),
		SpanID:     randomID(),
		Sampled:    randomID()%64 == 0,
		traceAttrs: make(map[string]string),
	}
}

// Implements the `Span` interface. Created via standardTracer (see
// `NewStandardTracer()`).
type standardSpan struct {
	lock     sync.Mutex
	tracer   *standardTracer
	recorder Recorder
	raw      RawSpan
}

func newChildContext(
	stdCtx *StandardContext,
) (*StandardContext, opentracing.Tags) {
	stdCtx.tagLock.RLock()
	newTags := make(map[string]string, len(stdCtx.traceAttrs))
	for k, v := range stdCtx.traceAttrs {
		newTags[k] = v
	}
	stdCtx.tagLock.RUnlock()

	return &StandardContext{
		TraceID:    stdCtx.TraceID,
		SpanID:     randomID(),
		Sampled:    stdCtx.Sampled,
		traceAttrs: newTags,
	}, opentracing.Tags{"parent_span_id": stdCtx.SpanID}
}

func (s *standardSpan) StartChild(operationName string) opentracing.Span {
	childCtx, childTags := newChildContext(s.raw.StandardContext)
	return s.tracer.startSpanGeneric(operationName, childCtx, childTags)
}

func (s *standardSpan) SetTag(key string, value interface{}) opentracing.Span {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.raw.Tags[key] = value
	return s
}

func (s *standardSpan) LogEvent(event string) {
	s.Log(opentracing.LogData{
		Event: event,
	})
}

func (s *standardSpan) LogEventWithPayload(event string, payload interface{}) {
	s.Log(opentracing.LogData{
		Event:   event,
		Payload: payload,
	})
}

func (s *standardSpan) Log(ld opentracing.LogData) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if ld.Timestamp.IsZero() {
		ld.Timestamp = time.Now()
	}
	s.raw.Logs = append(s.raw.Logs, &ld)
}

func (s *standardSpan) Finish() {
	duration := time.Since(s.raw.Start)
	s.lock.Lock()
	defer s.lock.Unlock()

	s.raw.Duration = duration
	s.recorder.RecordSpan(&s.raw)
}

// Implements the `Tracer` interface.
type standardTracer struct {
	recorder Recorder
}

func (s *standardTracer) StartTrace(
	operationName string,
) opentracing.Span {
	return s.startSpanGeneric(
		operationName,
		NewRootStandardContext(),
		nil,
	)
}

func (s *standardTracer) JoinTrace(
	operationName string,
	parent interface{},
) opentracing.Span {
	if goCtx, ok := parent.(context.Context); ok {
		return s.startSpanWithGoContextParent(operationName, goCtx)
	} else if span, ok := parent.(opentracing.Span); ok {
		return s.startSpanWithSpanParent(operationName, span)
	} else {
		panic(fmt.Errorf("Invalid parent type: %v", reflect.TypeOf(parent)))
	}
}

func (s *standardTracer) startSpanWithGoContextParent(
	operationName string,
	parent context.Context,
) opentracing.Span {
	if oldSpan := opentracing.SpanFromGoContext(parent); oldSpan != nil {
		// XXX: unchecked cast
		stdSpan := oldSpan.(*standardSpan)
		childCtx, tags := newChildContext(stdSpan.raw.StandardContext)
		return s.startSpanGeneric(
			operationName,
			childCtx,
			tags,
		)
	}

	return s.startSpanGeneric(
		operationName,
		NewRootStandardContext(),
		nil,
	)
}

func (s *standardTracer) startSpanWithSpanParent(
	operationName string,
	parent opentracing.Span,
) opentracing.Span {
	childCtx, tags := newChildContext(parent.(*standardSpan).raw.StandardContext)
	return s.startSpanGeneric(
		operationName,
		childCtx,
		tags,
	)
}

// A helper for standardSpan creation.
func (s *standardTracer) startSpanGeneric(
	operationName string,
	childCtx *StandardContext,
	tags opentracing.Tags,
) opentracing.Span {
	if tags == nil {
		tags = opentracing.Tags{}
	}
	span := &standardSpan{
		tracer:   s,
		recorder: s.recorder,
		raw: RawSpan{
			StandardContext: childCtx,
			Operation:       operationName,
			Start:           time.Now(),
			Duration:        -1,
			Tags:            tags,
			Logs:            []*opentracing.LogData{},
		},
	}
	return span
}

func (s *standardTracer) PropagateSpanAsText(
	sp opentracing.Span,
) (
	contextIDMap map[string]string,
	attrsMap map[string]string,
) {
	sc := sp.(*standardSpan).raw.StandardContext
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

func (s *standardTracer) PropagateSpanAsBinary(
	sp opentracing.Span,
) (
	traceContextID []byte,
	traceAttrs []byte,
) {
	sc := sp.(*standardSpan).raw.StandardContext
	var err error
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, sc.TraceID)
	if err != nil {
		panic(err)
	}
	err = binary.Write(buf, binary.BigEndian, sc.SpanID)
	if err != nil {
		panic(err)
	}
	var sampledByte byte
	if sc.Sampled {
		sampledByte = 1
	}
	err = binary.Write(buf, binary.BigEndian, sampledByte)
	if err != nil {
		panic(err)
	}
	// XXX: support attributes
	return buf.Bytes(), []byte{}
}

func (s *standardTracer) JoinTraceFromBinary(
	operationName string,
	traceContextID []byte,
	traceAttrs []byte,
) (opentracing.Span, error) {
	var err error
	reader := bytes.NewReader(traceContextID)
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
	// XXX: support attributes
	return s.startSpanGeneric(
			operationName,
			&StandardContext{
				TraceID:    traceID,
				SpanID:     spanID,
				Sampled:    sampledByte != 0,
				traceAttrs: make(map[string]string),
			},
			opentracing.Tags{}),
		nil
}

func (s *standardTracer) JoinTraceFromText(
	operationName string,
	contextIDMap map[string]string,
	tagsMap map[string]string,
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

	lowercaseTagsMap := make(map[string]string, len(tagsMap))
	for k, v := range tagsMap {
		lowercaseTagsMap[strings.ToLower(k)] = v
	}

	return s.startSpanGeneric(
			operationName,
			&StandardContext{
				TraceID:    traceID,
				SpanID:     spanID,
				Sampled:    sampled,
				traceAttrs: lowercaseTagsMap,
			},
			opentracing.Tags{}),
		nil
}

const (
	fieldNameTraceID = "traceid"
	fieldNameSpanID  = "spanid"
	fieldNameSampled = "sampled"
)

func (s *standardSpan) SetTraceAttribute(restrictedKey, val string) opentracing.Span {
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.raw.StandardContext.tagLock.Lock()
	defer s.raw.StandardContext.tagLock.Unlock()

	s.raw.StandardContext.traceAttrs[canonicalKey] = val
	return s
}

func (s *standardSpan) TraceAttribute(restrictedKey string) string {
	canonicalKey, valid := opentracing.CanonicalizeTraceAttributeKey(restrictedKey)
	if !valid {
		panic(fmt.Errorf("Invalid key: %q", restrictedKey))
	}

	s.raw.StandardContext.tagLock.RLock()
	defer s.raw.StandardContext.tagLock.RUnlock()

	return s.raw.StandardContext.traceAttrs[canonicalKey]
}

var (
	seededIDGen  = rand.New(rand.NewSource(time.Now().UnixNano()))
	seededIDLock sync.Mutex
)

func randomID() int64 {
	// The golang rand generators are *not* intrinsically thread-safe.
	seededIDLock.Lock()
	defer seededIDLock.Unlock()
	return seededIDGen.Int63()
}
