package node

import (
	"sync"

	"github.com/BullionBear/sequex/pkg/log"

	"github.com/nats-io/nats.go"
)

type Node interface {
	// Name returns the name of the node
	Name() string

	Start() error
	Shutdown() error
	// WaitForShutdown waits for shutdown signal

	// Digesting messages and publish to the next
	AddSubscription(string)
	Subscriptions() []string
	OnMessage(msg *nats.Msg)

	// OnRPC is called when an RPC is received
	OnRPC(req *nats.Msg) *nats.Msg

	// NATSConnection returns the NATS connection
	NATSConnection() *nats.Conn
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	name   string
	mutex  sync.Mutex
	nc     *nats.Conn
	msgCh  chan *nats.Msg
	subs   []string
	logger log.Logger
}

func NewBaseNode(name string, nc *nats.Conn, sz int, logger log.Logger) *BaseNode {
	return &BaseNode{
		name:   name,
		nc:     nc,
		msgCh:  make(chan *nats.Msg, sz),
		subs:   make([]string, 0),
		mutex:  sync.Mutex{},
		logger: logger,
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

func (bn *BaseNode) Info(msg string, fields ...log.Field) {
	bn.logger.Info(msg, fields...)
}

func (bn *BaseNode) Infof(format string, v ...any) {
	bn.logger.Infof(format, v...)
}

func (bn *BaseNode) Error(msg string, fields ...log.Field) {
	bn.logger.Error(msg, fields...)
}

func (bn *BaseNode) Errorf(format string, v ...any) {
	bn.logger.Errorf(format, v...)
}

func (bn *BaseNode) Debug(msg string, fields ...log.Field) {
	bn.logger.Debug(msg, fields...)
}

func (bn *BaseNode) Debugf(format string, v ...any) {
	bn.logger.Debugf(format, v...)
}

func (bn *BaseNode) Warn(msg string, fields ...log.Field) {
	bn.logger.Warn(msg, fields...)
}

func (bn *BaseNode) Warnf(format string, v ...any) {
	bn.logger.Warnf(format, v...)
}

func (bn *BaseNode) Fatal(msg string, fields ...log.Field) {
	bn.logger.Fatal(msg, fields...)
}

func (bn *BaseNode) Fatalf(format string, v ...any) {
	bn.logger.Fatalf(format, v...)
}
