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
	HTTPStatusCode = uint16Tag("htttp.status_code")
)

// ---

// SpanKindEnum represents common span types
type SpanKindEnum string

type spanKindTag string

// Add adds a string tag to the `span`
func (tag spanKindTag) Set(span opentracing.Span, value SpanKindEnum) {
	span.SetTag(string(tag), value)
}

// ---

type stringTag string

// Add adds a string tag to the `span`
func (tag stringTag) Set(span opentracing.Span, value string) {
	span.SetTag(string(tag), value)
}

// ---

type uint32Tag string

// Add adds a uint32 tag to the `span`
func (tag uint32Tag) Set(span opentracing.Span, value uint32) {
	span.SetTag(string(tag), value)
}

// ---

type uint16Tag string

// Add adds a uint16 tag to the `span`
func (tag uint16Tag) Set(span opentracing.Span, value uint16) {
	span.SetTag(string(tag), value)
}
