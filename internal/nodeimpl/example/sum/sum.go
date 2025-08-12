package sum

import (
	"encoding/json"
	"fmt"
	"sync"

	pb "github.com/BullionBear/sequex/internal/model/protobuf/example"
	"github.com/BullionBear/sequex/internal/nodeimpl/utils"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type SumConfig struct {
	InitSum    int64 `json:"init_sum"`
	UpperLimit int64 `json:"upper_limit"`
	LowerLimit int64 `json:"lower_limit"`
}

type SumNode struct {
	*node.BaseNode
	cfg SumConfig

	sum    int64
	count  int64
	mutex  sync.Mutex
	logger log.Logger
}

func init() {
	node.RegisterNode("sum", NewSumNode)
}

func NewSumNode(name string, nc *nats.Conn, config node.NodeConfig, logger *log.Logger) (node.Node, error) {
	// Parse configuration
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %v", err)
	}
	var cfg SumConfig
	if err := json.Unmarshal(jsonBytes, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	baseNode := node.NewBaseNode(name, nc, 100, *logger)

	sumNode := &SumNode{
		BaseNode: baseNode,
		cfg:      cfg,
		sum:      cfg.InitSum,
		count:    0,
	}

	baseNode.Logger().Info("SUM node created",
		log.Int64("init_sum", cfg.InitSum),
		log.Int64("upper_limit", cfg.UpperLimit),
		log.Int64("lower_limit", cfg.LowerLimit),
	)

	return sumNode, nil
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

func (n *SumNode) OnMessage(msg *nats.Msg) {
	contentType := msg.Header.Get("Content-Type")
	messageType := msg.Header.Get("Message-Type")

	n.Logger().Debug("Received message",
		log.String("content_type", contentType),
		log.String("message_type", messageType),
	)

	switch {
	case contentType == "application/protobuf" && messageType == "rng.RngMessage":
		var content pb.RngMessage
		if err := proto.Unmarshal(msg.Data, &content); err != nil {
			n.Logger().Error("Error unmarshalling RngMessage",
				log.Error(err),
			)
			return
		}
		n.onRngMessage(&content)
	default:
		n.Logger().Warn("Unknown message type",
			log.String("message_type", messageType),
		)
	}
}

func (n *SumNode) onRngMessage(msg *pb.RngMessage) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	oldSum := n.sum
	n.sum += msg.Random
	if n.sum > n.cfg.UpperLimit || n.sum < n.cfg.LowerLimit {
		n.sum = n.cfg.InitSum
		n.Logger().Info("Sum reset to initial value",
			log.Int64("old_sum", oldSum),
			log.Int64("random_value", msg.Random),
			log.Int64("new_sum", n.sum),
			log.String("reason", "limit_exceeded"),
		)
	}
	n.count++

	n.Logger().Info("Sum updated",
		log.Int64("old_sum", oldSum),
		log.Int64("random_value", msg.Random),
		log.Int64("new_sum", n.sum),
		log.Int64("count", n.count),
	)
}

func (n *SumNode) OnRPC(req *nats.Msg) *nats.Msg {
	contentType := req.Header.Get("Content-Type")
	messageType := req.Header.Get("Message-Type")

	n.Logger().Debug("Received RPC request",
		log.String("content_type", contentType),
		log.String("message_type", messageType),
	)

	switch {
	case contentType == "application/protobuf" && messageType == "sum.SumRequest":
		var content pb.SumRequest
		if err := proto.Unmarshal(req.Data, &content); err != nil {
			n.Logger().Error("Error unmarshalling SumRequest",
				log.Error(err),
			)
			return utils.MakeErrorMessage(utils.ErrorProtobufDeserialization, err)
		}
		return n.onSumRequest(&content)
	case contentType == "application/json" && messageType == "Config":
		return n.onConfig(&n.cfg)
	default:
		n.Logger().Warn("Unknown message type",
			log.String("message_type", messageType),
		)
		return utils.MakeErrorMessage(utils.ErrorUnknownMessageType, fmt.Errorf("unknown message type: %s", messageType))
	}
}

func (n *SumNode) onSumRequest(req *pb.SumRequest) *nats.Msg {
	n.Logger().Info("Processing SumRequest",
		log.Int64("offset", req.Offset),
	)

	response := &pb.SumResponse{
		NSum:   n.sum + req.Offset,
		NCount: n.count,
	}
	responseBytes, err := utils.MarshallProtobuf(response)
	if err != nil {
		n.Logger().Error("Error marshalling SumResponse",
			log.Error(err),
		)
		return utils.MakeErrorMessage(utils.ErrorProtobufSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Message-Type", "sum.SumResponse")
	msg.Data = responseBytes
	return &msg
}

func (n *SumNode) onConfig(content *SumConfig) *nats.Msg {
	n.Logger().Debug("Processing config request")
	responseBytes, err := json.Marshal(content)
	if err != nil {
		n.Logger().Error("Error marshalling config response",
			log.Error(err),
		)
		return utils.MakeErrorMessage(utils.ErrorJSONSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/json")
	msg.Header.Set("Message-Type", "SumConfig")
	msg.Data = responseBytes
	return &msg
}
