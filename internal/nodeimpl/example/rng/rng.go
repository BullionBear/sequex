package rng

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	pb "github.com/BullionBear/sequex/internal/model/protobuf/example"
	"github.com/BullionBear/sequex/internal/nodeimpl/utils"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type RNGConfig struct {
	Low      int     `json:"low"`
	High     int     `json:"high"`
	Interval float64 `json:"interval"`
}

type RNGNode struct {
	*node.BaseNode
	// Configurable parameters
	cfg RNGConfig

	// State variables
	rand   *rand.Rand
	mutex  sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	count  int64
}

func init() {
	node.RegisterNode("rng", NewRNGNode)
}

func NewRNGNode(name string, nc *nats.Conn, config node.NodeConfig) (node.Node, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	var cfg RNGConfig
	if err := json.Unmarshal(jsonBytes, &cfg); err != nil {
		return nil, err
	}
	source := rand.NewPCG(uint64(time.Now().UnixNano()), 1)
	rand := rand.New(source)

	return &RNGNode{
		BaseNode: node.NewBaseNode(name, nc, 100),
		cfg:      cfg,
		rand:     rand,
		count:    0,
		ctx:      context.Background(),
		cancel: func() {
			log.Printf("ctx is not set")
		},
	}, nil
}

func (n *RNGNode) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	n.ctx = ctx
	n.cancel = cancel
	go n.publishRng(ctx)
	return nil
}

func (n *RNGNode) Shutdown() {
	n.cancel()
}

func (n *RNGNode) OnMessage(msg *nats.Msg) {
}

func (n *RNGNode) OnRPC(req *nats.Msg) *nats.Msg {
	contentType := req.Header.Get("Content-Type")
	messageType := req.Header.Get("Message-Type")

	switch {
	case contentType == "application/protobuf" && messageType == "rng.RngCountRequest":
		var content pb.RngCountRequest
		if err := proto.Unmarshal(req.Data, &content); err != nil {
			log.Printf("Error unmarshalling RngCountRequest: %v", err)
			return utils.MakeErrorMessage(utils.ErrorProtobufDeserialization, err)
		}
		return n.onRngCountRequest(&content)
	case contentType == "application/json" && messageType == "Config":
		return n.onConfig(&n.cfg)
	default:
		return utils.MakeErrorMessage(utils.ErrorUnknownMessageType, fmt.Errorf("unknown content-type: %s, message-type: %s", contentType, messageType))
	}
}

func (n *RNGNode) onRngCountRequest(req *pb.RngCountRequest) *nats.Msg {
	fmt.Println("onRngCountRequest", req.String())
	response := &pb.RngCountResponse{
		NCount: n.count,
	}
	responseBytes, err := utils.MarshallProtobuf(response)
	if err != nil {
		return utils.MakeErrorMessage(utils.ErrorProtobufSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Message-Type", "rng.RngCountResponse")
	msg.Data = responseBytes
	return &msg
}

func (n *RNGNode) onConfig(content *RNGConfig) *nats.Msg {
	responseBytes, err := json.Marshal(content)
	if err != nil {
		return utils.MakeErrorMessage(utils.ErrorJSONSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/json")
	msg.Header.Set("Message-Type", "RngConfig")
	msg.Data = responseBytes
	return &msg
}

func (n *RNGNode) publishRng(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(n.cfg.Interval * float64(time.Second)))
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			n.waitForShutdown()
			return
		case <-ticker.C:
			rand := n.rand.IntN(n.cfg.High-n.cfg.Low+1) + n.cfg.Low
			content := &pb.RngMessage{
				Random: int64(rand),
			}
			msgBytes, err := utils.MarshallProtobuf(content)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			// Create message with proper headers
			msg := &nats.Msg{
				Header: map[string][]string{
					"Content-Type": {"application/protobuf"},
					"Message-Type": {"rng.RngMessage"},
				},
				Data: msgBytes,
			}

			// Publish to the topic with the node name
			topic := fmt.Sprintf("%s.rng.RngMessage", n.Name())
			msg.Subject = topic
			if err := n.NATSConnection().PublishMsg(msg); err != nil {
				log.Printf("Error publishing message: %v", err)
				continue
			}
			n.count++
			log.Printf("Published random number: %d", rand)
		}
	}
}

func (n *RNGNode) waitForShutdown() {
	log.Printf("%s Waiting for shutdown", n.Name())
}
