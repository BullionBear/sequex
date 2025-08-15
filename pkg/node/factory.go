package node

import (
	"fmt"
	"sync"

	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
)

type NodeConfig = map[string]any

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

	// Extract the node name from config or generate one
	nodeName, ok := (*config)["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name not found in config")
	}

	return fn(nodeName, eb, config, logger)
}
