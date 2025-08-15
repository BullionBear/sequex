package base

import (
	"fmt"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/BullionBear/sequex/pkg/eventbus"
)

type BaseNode struct {
	name      string
	eb        *eventbus.EventBus
	logger    log.Logger
	createdAt int64
	params    map[string]any
	emit      map[string]string
	on        map[string]string
	rpc       map[string]string
}

func NewBaseNode(name string, eb *eventbus.EventBus, config *node.NodeConfig, logger log.Logger) *BaseNode {
	return &BaseNode{
		name:      name,
		eb:        eb,
		logger:    logger,
		createdAt: time.Now().Unix(),
		params:    config.Params,
		emit:      config.Emit,
		on:        config.On,
		rpc:       config.Rpc,
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

func (bn *BaseNode) RequestMetadata(req *pbCommon.MetadataRequest) *pbCommon.MetadataResponse {
	return &pbCommon.MetadataResponse{
		Id:        req.Id,
		Code:      pbCommon.ErrorCode_ERROR_CODE_OK,
		Message:   "",
		Name:      bn.name,
		CreatedAt: bn.createdAt,
		Emit:      bn.emit,
		On:        bn.on,
		Rpc:       bn.rpc,
	}
}

func (bn *BaseNode) GetEmit(key string) (string, error) {
	if value, ok := bn.emit[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("emit key %s not found", key)
}

func (bn *BaseNode) GetOn(key string) (string, error) {
	if value, ok := bn.on[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("on key %s not found", key)
}

func (bn *BaseNode) GetRpc(key string) (string, error) {
	if value, ok := bn.rpc[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("rpc key %s not found", key)
}
