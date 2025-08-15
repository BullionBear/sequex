package node

import (
	"github.com/BullionBear/sequex/pkg/eventbus"
)

type Node interface {
	Name() string
	Start() error
	Shutdown() error
	EventBus() *eventbus.EventBus
}
