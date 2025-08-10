package node

import (
	"sync"

	"github.com/nats-io/nats.go"
)

type Node interface {
	// Name returns the name of the node
	Name() string

	// Add subscription
	AddSubscription(string)

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
	mutex sync.Mutex
	nc    *nats.Conn
	msgCh chan *nats.Msg
	subs  []string
}

func NewBaseNode(name string, nc *nats.Conn, sz int) *BaseNode {
	return &BaseNode{
		name:  name,
		nc:    nc,
		msgCh: make(chan *nats.Msg, sz),
		subs:  make([]string, 0),
		mutex: sync.Mutex{},
	}
}

// Name returns the node name
func (bn *BaseNode) Name() string {
	return bn.name
}

func (bn *BaseNode) NATSConnection() *nats.Conn {
	return bn.nc
}

func (bn *BaseNode) AddSubscription(sub string) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()
	bn.subs = append(bn.subs, sub)
}

func (bn *BaseNode) Subscriptions() []string {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()
	return bn.subs
}
