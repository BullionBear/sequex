package node

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

type Node interface {
	// Name returns the name of the node
	Name() string

	// SetName sets the name of the node
	SetName(name string)

	// SetNATSConnection sets the NATS connection for the node
	SetNATSConnection(nc *nats.Conn)

	// Context returns the context for the node
	Context() context.Context

	// Start begins the node's operation
	Start() error

	// Stop gracefully shuts down the node
	Stop() error

	// GetNATSConnection returns the NATS connection
	GetNATSConnection() *nats.Conn

	// WaitForShutdown waits for shutdown signal
	WaitForShutdown()
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	name   string
	nc     *nats.Conn
	config map[string]interface{}
	ctx    context.Context
	cancel context.CancelFunc
}

// NewBaseNode creates a new base node
func NewBaseNode(name string) *BaseNode {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseNode{
		name:   name,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Name returns the node name
func (bn *BaseNode) Name() string {
	return bn.name
}

// SetName sets the name of the node
func (bn *BaseNode) SetName(name string) {
	bn.name = name
}

// SetNATSConnection sets the NATS connection for the node
func (bn *BaseNode) SetNATSConnection(nc *nats.Conn) {
	bn.nc = nc
}

// GetNATSConnection returns the NATS connection
func (bn *BaseNode) GetNATSConnection() *nats.Conn {
	return bn.nc
}

// Stop gracefully shuts down the base node
func (bn *BaseNode) Stop() error {
	bn.cancel()
	return nil
}

// GetContext returns the context for the base node
func (bn *BaseNode) Context() context.Context {
	return bn.ctx
}

// WaitForShutdown waits for shutdown signal
func (bn *BaseNode) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("[%s] Waiting for shutdown signal...", bn.name)
	<-sigChan
	log.Printf("[%s] Shutdown signal received", bn.name)
}
