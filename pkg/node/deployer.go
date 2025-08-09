package node

import (
	"fmt"
	"log"
	"sync"
)

// Deployer manages multiple microservice nodes
type Deployer struct {
	nodes map[string]Node
	mutex sync.RWMutex
}

// NewDeployer creates a new deployer instance
func NewDeployer() *Deployer {
	return &Deployer{
		nodes: make(map[string]Node),
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
