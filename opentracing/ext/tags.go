package ext

import (
	"github.com/opentracing/api-golang/opentracing"
)

var (
	// RPC tags as symmetrical, each can be emitted by either client-side of server-side

	// RPCService records the service name of the RPC peer
	RPCService = &stringTag{"rpc.service"}

	// RPCHostname records the host name of the RPC peer
	RPCHostname = &stringTag{"rpc.hostname"}

	// RPCHostIPv4 records IP v4 host address of the RPC peer
	RPCHostIPv4 = &uint32Tag{"rpc.ipv4"}

	// RPCHostIPv6 records IP v6 host address of the RPC peer
	RPCHostIPv6 = &stringTag{"rpc.ipv6"}

	// RPCPort records port number of the RPC peer
	RPCPort = &uint16Tag{"rpc.port"}
)

type stringTag struct {
	Key string
}

// Add adds a string tag to the `span`
func (tag *stringTag) Add(span opentracing.Span, value string) {
	span.SetTag(tag.Key, value)
}

type uint32Tag struct {
	Key string
}

// Add adds a uint32 tag to the `span`
func (tag *uint32Tag) Add(span opentracing.Span, value uint32) {
	span.SetTag(tag.Key, value)
}

type uint16Tag struct {
	Key string
}

// Add adds a uint16 tag to the `span`
func (tag *uint16Tag) Add(span opentracing.Span, value uint16) {
	span.SetTag(tag.Key, value)
}
