package base

import (
	"time"

	"github.com/BullionBear/sequex/pkg/log"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/BullionBear/sequex/pkg/eventbus"
)

type BaseNode struct {
	name      string
	eb        *eventbus.EventBus
	logger    log.Logger
	createdAt int64
}

func NewBaseNode(name string, eb *eventbus.EventBus, logger log.Logger) *BaseNode {
	return &BaseNode{
		name:      name,
		eb:        eb,
		logger:    logger,
		createdAt: time.Now().Unix(),
	}
}

// Name returns the node name
func (bn *BaseNode) Name() string {
	return bn.name
}

func (bn *BaseNode) Logger() log.Logger {
	return bn.logger
}

func (bn *BaseNode) EventBus() *eventbus.EventBus {
	return bn.eb
}

func (bn *BaseNode) CreatedAt() int64 {
	return bn.createdAt
}

func (bn *BaseNode) GetMetadata(pb *pbCommon.MetadataRequest) *pbCommon {
	return &pbCommon.MetadataResponse{
		Id:        bn.name,
		CreatedAt: bn.createdAt,
	}
}
