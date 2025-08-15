package sum

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	pbExample "github.com/BullionBear/sequex/internal/model/protobuf/example"
	"github.com/BullionBear/sequex/internal/nodeimpl/base"
	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
)

const (
	OnRandomKey = "on_random"

	RpcReqMetadataKey   = "req_metadata"
	RpcReqParametersKey = "req_parameters"
	RpcReqStatusKey     = "req_status"
	RpcReqSumKey        = "req_sum"
)

type SumConfig struct {
	InitSum    int64 `json:"init_sum"`
	UpperLimit int64 `json:"upper_limit"`
	LowerLimit int64 `json:"lower_limit"`
}

type SumNode struct {
	*base.BaseNode
	cfg SumConfig

	sum   int64
	count int64
	mutex sync.Mutex
}

func init() {
	node.RegisterNode("sum", NewSumNode)
}

func NewSumNode(name string, eb *eventbus.EventBus, config *node.NodeConfig, logger log.Logger) (node.Node, error) {
	baseNode := base.NewBaseNode(name, eb, config, logger)
	initSum, ok := config.Params["init_sum"].(int)
	if !ok {
		return nil, fmt.Errorf("init_sum is not an int")
	}
	upperLimit, ok := config.Params["upper_limit"].(int)
	if !ok {
		return nil, fmt.Errorf("upper_limit is not an int")
	}
	lowerLimit, ok := config.Params["lower_limit"].(int)
	if !ok {
		return nil, fmt.Errorf("lower_limit is not an int")
	}
	cfg := SumConfig{
		InitSum:    int64(initSum),
		UpperLimit: int64(upperLimit),
		LowerLimit: int64(lowerLimit),
	}
	return &SumNode{
		BaseNode: baseNode,
		cfg:      cfg,
		sum:      cfg.InitSum,
		count:    0,
	}, nil
}

func (n *SumNode) Start() error {
	n.Logger().Info("Starting SUM node")
	return nil
}

func (n *SumNode) Shutdown() error {
	n.Logger().Info("Shutting down SUM node",
		log.Int64("final_sum", n.sum),
		log.Int64("total_count", n.count),
	)
	return nil
}

func (n *SumNode) OnRandom(msg *pbExample.RandomNumberMessage) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.sum += int64(msg.Value)
	n.sum = int64(math.Min(float64(n.sum), float64(n.cfg.UpperLimit)))
	n.sum = int64(math.Max(float64(n.sum), float64(n.cfg.LowerLimit)))
	n.count++
	n.Logger().Info("Sum updated",
		log.Int64("sum", n.sum),
		log.Int64("random", msg.Value),
		log.Int64("count", n.count),
	)
}

func (n *SumNode) RequestParameters(req *pbCommon.ParametersRequest) *pbCommon.ParametersResponse {
	jsonBytes, err := json.Marshal(n.cfg)
	if err != nil {
		n.Logger().Error("Failed to marshal parameters", log.Error(err))
		return &pbCommon.ParametersResponse{
			Id:         req.Id,
			Code:       pbCommon.ErrorCode_ERROR_CODE_SERIALIZATION_ERROR,
			Message:    "Failed to json marshal parameters",
			Parameters: jsonBytes,
		}
	}
	return &pbCommon.ParametersResponse{
		Id:         req.Id,
		Code:       pbCommon.ErrorCode_ERROR_CODE_OK,
		Message:    "",
		Parameters: jsonBytes,
	}
}

func (n *SumNode) RequestStatus(req *pbCommon.StatusRequest) *pbCommon.StatusResponse {
	jsonBytes, err := json.Marshal(map[string]any{
		"sum":   n.sum,
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
