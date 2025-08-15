package node

import (
	"fmt"
	"sync"

	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
)

type NodeConfig struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Params map[string]any    `yaml:"params,omitempty"`
	On     map[string]string `yaml:"on,omitempty"`
	Emit   map[string]string `yaml:"emit,omitempty"`
	Rpc    map[string]string `yaml:"rpc,omitempty"`
}

type NewNodeFunc func(name string, eb *eventbus.EventBus, config *NodeConfig, logger log.Logger) (Node, error)

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

func CreateNode(nodeType string, eb *eventbus.EventBus, config *NodeConfig, logger log.Logger) (Node, error) {
	mu.RLock()
	defer mu.RUnlock()
	fn, ok := nodes[nodeType]
	if !ok {
		return nil, fmt.Errorf("node type %s not found in factory", nodeType)
	}

	return fn(config.Name, eb, config, logger)
}
