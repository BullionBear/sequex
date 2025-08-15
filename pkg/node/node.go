package node

import (
	"sync"

	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
)

type Node interface {
	Name() string
	Start() error
	Shutdown() error
	EventBus() *eventbus.EventBus
}

type BaseNode struct {
	name   string
	mutex  sync.Mutex
	eb     *eventbus.EventBus
	logger log.Logger
}

func NewBaseNode(name string, eb *eventbus.EventBus, logger log.Logger) *BaseNode {
	return &BaseNode{
		name:   name,
		eb:     eb,
		mutex:  sync.Mutex{},
		logger: logger,
	}
}

// Name returns the node name
func (bn *BaseNode) Name() string {
	return bn.name
}

func (bn *BaseNode) Logger() log.Logger {
	return bn.logger
}
