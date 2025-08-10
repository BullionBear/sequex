package sum

import (
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/BullionBear/sequex/internal/model/protobuf/example"
	"github.com/BullionBear/sequex/internal/nodeimpl/utils"
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

	sum   int64
	count int64
}

func init() {
	node.RegisterNode("sum", NewSumNode)
}

func NewSumNode(name string, nc *nats.Conn, config node.NodeConfig) (node.Node, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	var cfg SumConfig
	if err := json.Unmarshal(jsonBytes, &cfg); err != nil {
		return nil, err
	}
	return &SumNode{
		BaseNode: node.NewBaseNode(name, nc, 100),
		cfg:      cfg,
		sum:      cfg.InitSum,
		count:    0,
	}, nil
}

func (n *SumNode) Start() error {
	return nil
}

func (n *SumNode) Shutdown() {
}

func (n *SumNode) OnMessage(msg *nats.Msg) {
	contentType := msg.Header.Get("Content-Type")
	messageType := msg.Header.Get("Message-Type")

	switch {
	case contentType == "application/protobuf" && messageType == "rng.RngMessage":
		var content pb.RngMessage
		if err := proto.Unmarshal(msg.Data, &content); err != nil {
			log.Printf("Error unmarshalling RngMessage: %v", err)
		}
		n.onRngMessage(&content)
	default:
		log.Printf("Unknown message type: %s", messageType)
	}
}

func (n *SumNode) onRngMessage(msg *pb.RngMessage) {
	n.sum += msg.Random
	if n.sum > n.cfg.UpperLimit || n.sum < n.cfg.LowerLimit {
		n.sum = n.cfg.InitSum
	}
	n.count++
	log.Printf("Sum: %d, Count: %d", n.sum, n.count)
}

func (n *SumNode) OnRPC(req *nats.Msg) *nats.Msg {
	contentType := req.Header.Get("Content-Type")
	messageType := req.Header.Get("Message-Type")

	switch {
	case contentType == "application/protobuf" && messageType == "sum.SumRequest":
		var content pb.SumRequest
		if err := proto.Unmarshal(req.Data, &content); err != nil {
			log.Printf("Error unmarshalling SumRequest: %v", err)
		}
		return n.onSumRequest(&content)
	case contentType == "application/json" && messageType == "Config":
		return n.onConfig(&n.cfg)
	default:
		log.Printf("Unknown message type: %s", messageType)
		return utils.MakeErrorMessage(utils.ErrorUnknownMessageType, fmt.Errorf("unknown message type: %s", messageType))
	}
}

func (n *SumNode) onSumRequest(req *pb.SumRequest) *nats.Msg {
	response := &pb.SumResponse{
		NSum:   n.sum + req.Offset,
		NCount: n.count,
	}
	responseBytes, err := utils.MarshallProtobuf(response)
	if err != nil {
		log.Printf("Error marshalling SumRequest: %v", err)
		return utils.MakeErrorMessage(utils.ErrorProtobufSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Message-Type", "sum.SumResponse")
	msg.Data = responseBytes
	return &msg
}

func (n *SumNode) onConfig(content *SumConfig) *nats.Msg {
	fmt.Println("onConfig", content)
	responseBytes, err := json.Marshal(content)
	if err != nil {
		return utils.MakeErrorMessage(utils.ErrorJSONSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/json")
	msg.Header.Set("Message-Type", "SumConfig")
	msg.Data = responseBytes
	return &msg
}
