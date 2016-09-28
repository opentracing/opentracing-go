package ext

import opentracing "github.com/opentracing/opentracing-go"

// These constants define common tag keys recommended for better portability across
// tracing systems and languages/platforms.
//
// The tag keys are defined as typed strings, so that in addition to the usual use
//
//     span.SetTag(string(TagKey), value)
//
// they also support value type validation via this additional syntax:
//
//    TagKey.Set(span, value)
//
var (
	//////////////////////////////////////////////////////////////////////
	// SpanKind (client/server) tag keys
	//////////////////////////////////////////////////////////////////////

	// SpanKind hints at relationship between spans, e.g. client/server
	SpanKind = spanKindTagKey("span.kind")

	// SpanKindRPCClient marks a span representing the client-side of an RPC
	// or other remote call
	SpanKindRPCClientTagValue = SpanKindTagValue("client")

	// SpanKindRPCServer marks a span representing the server-side of an RPC
	// or other remote call
	SpanKindRPCServerTagValue = SpanKindTagValue("server")

	//////////////////////////////////////////////////////////////////////
	// Component tag keys
	//////////////////////////////////////////////////////////////////////

	// Component is a low-cardinality identifier of the module, library,
	// or package that is generating a span.
	Component = stringTagKey("component")

	//////////////////////////////////////////////////////////////////////
	// Sampling hint tag key
	//////////////////////////////////////////////////////////////////////

	// SamplingPriority determines the priority of sampling this Span.
	SamplingPriority = uint16TagKey("sampling.priority")

	//////////////////////////////////////////////////////////////////////
	// Peer tag keys. These tag keys can be emitted by either client-side of
	// server-side to describe the other side/service in a peer-to-peer
	// communications, like an RPC call.
	//////////////////////////////////////////////////////////////////////

	// PeerService records the service name of the peer
	PeerService = stringTagKey("peer.service")

	// PeerHostname records the host key of the peer
	PeerHostname = stringTagKey("peer.hostname")

	// PeerHostIPv4 records IP v4 host address of the peer
	PeerHostIPv4 = uint32TagKey("peer.ipv4")

	// PeerHostIPv6 records IP v6 host address of the peer
	PeerHostIPv6 = stringTagKey("peer.ipv6")

	// PeerPort records port number of the peer
	PeerPort = uint16TagKey("peer.port")

	//////////////////////////////////////////////////////////////////////
	// HTTP Tag keys
	//////////////////////////////////////////////////////////////////////

	// HTTPUrl should be the URL of the request being handled in this segment
	// of the trace, in standard URI format. The protocol is optional.
	HTTPUrl = stringTagKey("http.url")

	// HTTPMethod is the HTTP method of the request, and is case-insensitive.
	HTTPMethod = stringTagKey("http.method")

	// HTTPStatusCode is the numeric HTTP status code (200, 404, etc) of the
	// HTTP response.
	HTTPStatusCode = uint16TagKey("http.status_code")

	//////////////////////////////////////////////////////////////////////
	// Error Tag key
	//////////////////////////////////////////////////////////////////////

	// Error indicates that operation represented by the span resulted in an error.
	Error = boolTagKey("error")
)

var (
	//////////////////////////////////////////////////////////////////////
	// Conventional SpanKind Tags
	//////////////////////////////////////////////////////////////////////

	// SpanKindRPCClient is a tag indicating the span is a RPC client.
	SpanKindRPCClient = opentracing.Tag{Key: string(SpanKind), Value: SpanKindRPCClientTagValue}

	// SpanKindRPCServer is a tag indicating the span is a RPC client.
	SpanKindRPCServer = opentracing.Tag{Key: string(SpanKind), Value: SpanKindRPCServerTagValue}
)

// ---

// SpanKindTagValue represents common span types
type SpanKindTagValue string

type spanKindTagKey string

// Set adds a string tag to the `span`
func (tag spanKindTagKey) Set(span opentracing.Span, value SpanKindTagValue) {
	span.SetTag(string(tag), value)
}

type rpcServerOption struct {
	clientContext opentracing.SpanContext
}

func (r rpcServerOption) Apply(o *opentracing.StartSpanOptions) {
	if r.clientContext != nil {
		opentracing.ChildOf(r.clientContext).Apply(o)
	}
	SpanKindRPCServer.Apply(o)
}

// RPCServerOption returns a StartSpanOption appropriate for an RPC server span
// with `client` representing the metadata for the remote peer Span if available.
// In case client == nil, due to the client not being instrumented, this RPC
// server span will be a root span.
func RPCServerOption(client opentracing.SpanContext) opentracing.StartSpanOption {
	return rpcServerOption{client}
}

// ---

type stringTagKey string

// Set adds a string tag to the `span`
func (tag stringTagKey) Set(span opentracing.Span, value string) {
	span.SetTag(string(tag), value)
}

// ---

type uint32TagKey string

// Set adds a uint32 tag to the `span`
func (tag uint32TagKey) Set(span opentracing.Span, value uint32) {
	span.SetTag(string(tag), value)
}

// ---

type uint16TagKey string

// Set adds a uint16 tag to the `span`
func (tag uint16TagKey) Set(span opentracing.Span, value uint16) {
	span.SetTag(string(tag), value)
}

// ---

type boolTagKey string

// Add adds a bool tag to the `span`
func (tag boolTagKey) Set(span opentracing.Span, value bool) {
	span.SetTag(string(tag), value)
}
