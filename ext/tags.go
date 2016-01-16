package ext

import (
	"github.com/opentracing/opentracing-go"
)

var (
	// PeerXXX tags can be emitted by either client-side of server-side to describe
	// the other side/service in a peer-to-peer communications, like an RPC call.

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
)

type stringTag string

// Add adds a string tag to the `span`
func (tag stringTag) Add(span opentracing.Span, value string) {
	span.SetTag(string(tag), value)
}

type uint32Tag string

// Add adds a uint32 tag to the `span`
func (tag uint32Tag) Add(span opentracing.Span, value uint32) {
	span.SetTag(string(tag), value)
}

type uint16Tag string

// Add adds a uint16 tag to the `span`
func (tag uint16Tag) Add(span opentracing.Span, value uint16) {
	span.SetTag(string(tag), value)
}
