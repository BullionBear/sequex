package rng

import (
	"encoding/json"
	"math/rand"
	"time"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	pbExample "github.com/BullionBear/sequex/internal/model/protobuf/example"
	"github.com/BullionBear/sequex/internal/nodeimpl/base"
	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
	"google.golang.org/protobuf/proto"
)

const (
	EmitRandomKey = "emit_random"

	RpcReqMetadataKey   = "req_metadata"
	RpcReqParametersKey = "req_parameters"
	RpcReqStatusKey     = "req_status"
)

type RNGParams struct {
	Low      int           `json:"low"`
	High     int           `json:"high"`
	Interval time.Duration `json:"interval"`
	Seed     int           `json:"seed"`
}

type RNGNode struct {
	*base.BaseNode
	// Configurable parameters
	cfg RNGParams

	// State variables
	rand      *rand.Rand
	shutdownC chan struct{}
	doneC     chan struct{}
	count     int64
}

func init() {
	node.RegisterNode("rng", NewRNGNode)
}

func NewRNGNode(name string, eb *eventbus.EventBus, config *node.NodeConfig, logger log.Logger) (node.Node, error) {
	baseNode := base.NewBaseNode(name, eb, config, logger)

	cfg := RNGParams{
		Low:      config.Params["low"].(int),
		High:     config.Params["high"].(int),
		Interval: time.Duration(config.Params["interval"].(int)) * time.Second,
		Seed:     config.Params["seed"].(int),
	}

	return &RNGNode{
		BaseNode:  baseNode,
		cfg:       cfg,
		rand:      rand.New(rand.NewSource(int64(cfg.Seed))),
		shutdownC: make(chan struct{}),
		doneC:     make(chan struct{}),
	}, nil
}

func (n *RNGNode) Start() error {
	n.Logger().Info("Starting RNG node")

	go n.emitRandom(n.shutdownC, n.doneC)
	if metadata, err := n.GetRpc(RpcReqMetadataKey); err != nil {
		return err
	} else {
		n.EventBus().RegisterRPC(metadata, func(event proto.Message) proto.Message {
			if req, ok := event.(*pbCommon.MetadataRequest); ok {
				return n.RequestMetadata(req)
			}
			return &pbCommon.MetadataResponse{
				Id:      -1,
				Code:    pbCommon.ErrorCode_ERROR_CODE_INVALID_REQUEST,
				Message: "Invalid request",
			}
		})
	}
	if parameters, err := n.GetRpc(RpcReqParametersKey); err != nil {
		return err
	} else {
		n.EventBus().RegisterRPC(parameters, func(event proto.Message) proto.Message {
			if req, ok := event.(*pbCommon.ParametersRequest); ok {
				return n.RequestParameters(req)
			}
			return &pbCommon.ParametersResponse{
				Id:      -1,
				Code:    pbCommon.ErrorCode_ERROR_CODE_INVALID_REQUEST,
				Message: "Invalid request",
			}
		})
	}
	if status, err := n.GetRpc(RpcReqStatusKey); err != nil {
		return err
	} else {
		n.EventBus().RegisterRPC(status, func(event proto.Message) proto.Message {
			if req, ok := event.(*pbCommon.StatusRequest); ok {
				return n.RequestStatus(req)
			}
			return &pbCommon.StatusResponse{
				Id:      -1,
				Code:    pbCommon.ErrorCode_ERROR_CODE_INVALID_REQUEST,
				Message: "Invalid request",
			}
		})
	}
	return nil
}

func (n *RNGNode) Shutdown() error {
	n.Logger().Info("Shutting down RNG node")
	close(n.shutdownC)
	<-n.doneC
	return nil
}
func (n *RNGNode) emitRandom(shutdownC chan struct{}, doneC chan struct{}) {
	ticker := time.NewTicker(n.cfg.Interval)
	defer ticker.Stop()
	subject, err := n.GetEmit(EmitRandomKey)
	if err != nil {
		n.Logger().Fatal("Failed to get emit subject", log.Error(err))
		return
	}
	for {
		select {
		case <-shutdownC:
			n.Logger().Info("Stopping emitRandom")
			close(doneC)
			return
		case <-ticker.C:
			random := n.rand.Intn(n.cfg.High-n.cfg.Low+1) + n.cfg.Low
			n.EventBus().Emit(subject, &pbExample.RandomNumberMessage{
				Id:    n.count,
				Value: int64(random),
			})
			n.Logger().Info("Emitting random number", log.Int64("random", int64(random)), log.Int64("count", n.count))
			n.count++
		}
	}
}

func (n *RNGNode) RequestParameters(req *pbCommon.ParametersRequest) *pbCommon.ParametersResponse {
	jsonBytes, err := json.Marshal(map[string]any{
		"low":      n.cfg.Low,
		"high":     n.cfg.High,
		"interval": n.cfg.Interval.Seconds(),
		"seed":     n.cfg.Seed,
	})
	if err != nil {
		n.Logger().Error("Failed to marshal parameters", log.Error(err))
		return &pbCommon.ParametersResponse{
			Id:      -1,
			Code:    pbCommon.ErrorCode_ERROR_CODE_SERIALIZATION_ERROR,
			Message: "Failed to json marshal parameters",
		}
	}
	return &pbCommon.ParametersResponse{
		Id:         req.Id,
		Code:       pbCommon.ErrorCode_ERROR_CODE_OK,
		Message:    "",
		Parameters: jsonBytes,
	}
}

func (n *RNGNode) RequestStatus(req *pbCommon.StatusRequest) *pbCommon.StatusResponse {
	jsonBytes, err := json.Marshal(map[string]any{
		"count": n.count,
	})
	if err != nil {
		n.Logger().Error("Failed to marshal status", log.Error(err))
		return &pbCommon.StatusResponse{
			Id:      req.Id,
			Code:    pbCommon.ErrorCode_ERROR_CODE_SERIALIZATION_ERROR,
			Message: "Failed to json marshal status",
		}
	}
	return &pbCommon.StatusResponse{
		Id:      req.Id,
		Code:    pbCommon.ErrorCode_ERROR_CODE_OK,
		Message: "",
		Status:  jsonBytes,
	}
}
