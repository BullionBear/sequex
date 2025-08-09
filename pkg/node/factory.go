package node

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

type NodeConfig = map[string]any

type NewNodeFunc func(name string, nc *nats.Conn, config NodeConfig) (Node, error)

var (
	nodes map[string]NewNodeFunc
	mu    sync.RWMutex
)

func init() {
	nodes = make(map[string]NewNodeFunc)
}

func RegisterNode(name string, fn NewNodeFunc) {
	mu.Lock()
	defer mu.Unlock()
	nodes[name] = fn
}

func CreateNode(name string, nc *nats.Conn, config NodeConfig) (Node, error) {
	mu.RLock()
	defer mu.RUnlock()
	fn, ok := nodes[name]
	if !ok {
		return nil, fmt.Errorf("node %s not found in factory", name)
	}
	return fn(name, nc, config)
}
