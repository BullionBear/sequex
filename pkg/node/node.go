package node

import (
	"github.com/nats-io/nats.go"
)

type Node interface {
	// Name returns the name of the node
	Name() string

	// Subscriptions
	Subscriptions() []string

	// OnMessage is called when a message is received
	OnMessage(msg *nats.Msg)

	// OnRPC is called when an RPC is received
	OnRPC(req *nats.Msg) *nats.Msg

	// WaitForShutdown waits for shutdown signal
	WaitForShutdown()

	// NATSConnection returns the NATS connection
	NATSConnection() *nats.Conn
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	name  string
	nc    *nats.Conn
	msgCh chan *nats.Msg
}

func NewBaseNode(name string, nc *nats.Conn, sz int) *BaseNode {
	return &BaseNode{
		name:  name,
		nc:    nc,
		msgCh: make(chan *nats.Msg, sz),
	}
}

// Name returns the node name
func (bn *BaseNode) Name() string {
	return bn.name
}

func (bn *BaseNode) GetNATSConnection() *nats.Conn {
	return bn.nc
}
