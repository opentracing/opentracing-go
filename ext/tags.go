package ext

import opentracing "github.com/opentracing/opentracing-go"

// These constants define common tag names recommended for better portability across
// tracing systems and languages/platforms.
//
// The tag names are defined as typed strings, so that in addition to the usual use
//
//     span.setTag(TagName, value)
//
// they also support value type validation via this additional syntax:
//
//    TagName.Set(span, value)
//
var (
	//////////////////////////////////////////////////////////////////////
	// SpanKind (client/server)
	//////////////////////////////////////////////////////////////////////

	// SpanKind hints at relationship between spans, e.g. client/server
	SpanKind = spanKindTag("span.kind")

	// SpanKindRPCClient marks a span representing the client-side of an RPC
	// or other remote call
	SpanKindRPCClient = SpanKindEnum("client")

	// SpanKindRPCServer marks a span representing the server-side of an RPC
	// or other remote call
	SpanKindRPCServer = SpanKindEnum("server")

	//////////////////////////////////////////////////////////////////////
	// Component name
	//////////////////////////////////////////////////////////////////////

	// Component is a low-cardinality identifier of the module, library,
	// or package that is generating a span.
	Component = stringTag("component")

	//////////////////////////////////////////////////////////////////////
	// Sampling hint
	//////////////////////////////////////////////////////////////////////

	// SamplingPriority determines the priority of sampling this Span.
	SamplingPriority = uint16Tag("sampling.priority")

	//////////////////////////////////////////////////////////////////////
	// Peer tags. These tags can be emitted by either client-side of
	// server-side to describe the other side/service in a peer-to-peer
	// communications, like an RPC call.
	//////////////////////////////////////////////////////////////////////

	// PeerService records the service name of the peer
	PeerService = stringTag("peer.service")

	// PeerHostname records the host name of the peer
	PeerHostname = stringTag("peer.hostname")

	// PeerHostIPv4 records IP v4 host address of the peer
	PeerHostIPv4 = uint32Tag("peer.ipv4")

	// PeerHostIPv6 records IP v6 host address of the peer
	PeerHostIPv6 = stringTag("peer.ipv6")

	// PeerPort records port number of the peer
	PeerPort = uint16Tag("peer.port")

	//////////////////////////////////////////////////////////////////////
	// HTTP Tags
	//////////////////////////////////////////////////////////////////////

	// HTTPUrl should be the URL of the request being handled in this segment
	// of the trace, in standard URI format. The protocol is optional.
	HTTPUrl = stringTag("http.url")

	// HTTPMethod is the HTTP method of the request, and is case-insensitive.
	HTTPMethod = stringTag("http.method")

	// HTTPStatusCode is the numeric HTTP status code (200, 404, etc) of the
	// HTTP response.
	HTTPStatusCode = uint16Tag("http.status_code")

	//////////////////////////////////////////////////////////////////////
	// Error Tag
	//////////////////////////////////////////////////////////////////////

	// Error indicates that operation represented by the span resulted in an error.
	Error = boolTag("error")
)

// ---

// SpanKindEnum represents common span types
type SpanKindEnum string

type spanKindTag string

// Set adds a string tag to the `span`
func (tag spanKindTag) Set(span opentracing.Span, value SpanKindEnum) {
	span.SetTag(string(tag), value)
}

type rpcServerOption struct {
	clientContext opentracing.SpanContext
}

func (r rpcServerOption) Apply(o *opentracing.StartSpanOptions) {
	opentracing.ChildOf(r.clientContext).Apply(o)
	(opentracing.Tags{string(SpanKind): SpanKindRPCServer}).Apply(o)
}

// RPCServerOption returns a StartSpanOption appropriate for an RPC server span
// with `client` representing the metadata for the remote peer Span.
func RPCServerOption(client opentracing.SpanContext) opentracing.StartSpanOption {
	return rpcServerOption{client}
}

// ---

type stringTag string

// Set adds a string tag to the `span`
func (tag stringTag) Set(span opentracing.Span, value string) {
	span.SetTag(string(tag), value)
}

// ---

type uint32Tag string

// Set adds a uint32 tag to the `span`
func (tag uint32Tag) Set(span opentracing.Span, value uint32) {
	span.SetTag(string(tag), value)
}

// ---

type uint16Tag string

// Set adds a uint16 tag to the `span`
func (tag uint16Tag) Set(span opentracing.Span, value uint16) {
	span.SetTag(string(tag), value)
}

// ---

type boolTag string

// Add adds a bool tag to the `span`
func (tag boolTag) Set(span opentracing.Span, value bool) {
	span.SetTag(string(tag), value)
}
