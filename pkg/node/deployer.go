package node

import (
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

// Deployer manages multiple microservice nodes
type Deployer struct {
	nodes         map[string]Node
	mutex         sync.RWMutex
	subscriptions map[string][]*nats.Subscription
}

// NewDeployer creates a new deployer instance
func NewDeployer() *Deployer {
	return &Deployer{
		nodes:         make(map[string]Node),
		subscriptions: make(map[string][]*nats.Subscription),
	}
}

// RegisterNode registers a node with the deployer
func (d *Deployer) RegisterNode(n Node) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.nodes[n.Name()]; exists {
		return fmt.Errorf("node %s already registered", n.Name())
	}

	d.nodes[n.Name()] = n
	log.Printf("Registered node: %s", n.Name())
	return nil
}

func (d *Deployer) Start(name string) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	node, ok := d.nodes[name]
	if !ok {
		return fmt.Errorf("node %s not found", name)
	}
	nc := node.NATSConnection()
	for _, sub := range node.Subscriptions() {
		subscription, err := nc.Subscribe(sub, node.OnMessage)
		if err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", sub, err)
		}
		d.subscriptions[name] = append(d.subscriptions[name], subscription)
	}

	rpcSub, err := nc.Subscribe(fmt.Sprintf("rpc.%s", node.Name()), func(msg *nats.Msg) {
		response := node.OnRPC(msg)
		msg.RespondMsg(response)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to rpc.%s: %w", node.Name(), err)
	}
	d.subscriptions[name] = append(d.subscriptions[name], rpcSub)

	if err := node.Start(); err != nil {
		return fmt.Errorf("failed to start node %s: %w", name, err)
	}

	return nil
}

func (d *Deployer) Stop(name string) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	node, ok := d.nodes[name]
	if !ok {
		return fmt.Errorf("node %s not found", name)
	}
	if err := node.Shutdown(); err != nil {
		log.Printf("failed to shutdown node %s: %v", name, err)
	}
	for _, sub := range d.subscriptions[name] {
		sub.Unsubscribe()
	}
	delete(d.subscriptions, name)
	delete(d.nodes, name)

	return nil
}
