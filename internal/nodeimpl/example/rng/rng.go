package rng

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	pb "github.com/BullionBear/sequex/internal/model/protobuf/example"
	"github.com/BullionBear/sequex/internal/nodeimpl/utils"
	"github.com/BullionBear/sequex/pkg/log"
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
	logger log.Logger
}

func init() {
	node.RegisterNode("rng", NewRNGNode)
}

func NewRNGNode(name string, nc *nats.Conn, config node.NodeConfig) (node.Node, error) {
	// Initialize structured logger
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithTimeRotation("./logs", fmt.Sprintf("%s.log", name), 24*time.Hour, 7),
	)

	// Parse configuration
	var cfg RNGConfig
	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %v", err)
	}
	if err := json.Unmarshal(configBytes, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Create base node
	baseNode := node.NewBaseNode(name, nc, 100)

	// Create RNG node
	rngNode := &RNGNode{
		BaseNode: baseNode,
		cfg:      cfg,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:   logger.With(log.String("node_name", name), log.String("node_type", "rng")),
	}

	rngNode.ctx, rngNode.cancel = context.WithCancel(context.Background())

	rngNode.logger.Info("RNG node created",
		log.Int("low", cfg.Low),
		log.Int("high", cfg.High),
		log.Float64("interval", cfg.Interval),
	)

	return rngNode, nil
}

func (n *RNGNode) Start() error {
	n.logger.Info("Starting RNG node")
	go n.publishRng(n.ctx)
	return nil
}

func (n *RNGNode) Shutdown() error {
	n.logger.Info("Shutting down RNG node")
	n.cancel()
	return nil
}

func (n *RNGNode) OnMessage(msg *nats.Msg) {
}

func (n *RNGNode) OnRPC(req *nats.Msg) *nats.Msg {
	contentType := req.Header.Get("Content-Type")
	messageType := req.Header.Get("Message-Type")

	n.logger.Debug("Received RPC request",
		log.String("content_type", contentType),
		log.String("message_type", messageType),
	)

	switch {
	case contentType == "application/protobuf" && messageType == "rng.RngCountRequest":
		var content pb.RngCountRequest
		if err := proto.Unmarshal(req.Data, &content); err != nil {
			n.logger.Error("Error unmarshalling RngCountRequest",
				log.Error(err),
			)
			return utils.MakeErrorMessage(utils.ErrorProtobufDeserialization, err)
		}
		return n.onRngCountRequest(&content)
	case contentType == "application/json" && messageType == "Config":
		return n.onConfig(&n.cfg)
	default:
		n.logger.Warn("Unknown message type",
			log.String("content_type", contentType),
			log.String("message_type", messageType),
		)
		return utils.MakeErrorMessage(utils.ErrorUnknownMessageType, fmt.Errorf("unknown content-type: %s, message-type: %s", contentType, messageType))
	}
}

func (n *RNGNode) onRngCountRequest(req *pb.RngCountRequest) *nats.Msg {
	n.logger.Debug("Processing RngCountRequest",
		log.String("request", req.String()),
	)

	response := &pb.RngCountResponse{
		NCount: n.count,
	}
	responseBytes, err := utils.MarshallProtobuf(response)
	if err != nil {
		n.logger.Error("Error marshalling RngCountResponse",
			log.Error(err),
		)
		return utils.MakeErrorMessage(utils.ErrorProtobufSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Message-Type", "rng.RngCountResponse")
	msg.Data = responseBytes
	return &msg
}

func (n *RNGNode) onConfig(content *RNGConfig) *nats.Msg {
	n.logger.Debug("Processing config request")
	responseBytes, err := json.Marshal(content)
	if err != nil {
		n.logger.Error("Error marshalling config response",
			log.Error(err),
		)
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

	n.logger.Info("Starting RNG publishing loop",
		log.Float64("interval_seconds", n.cfg.Interval),
	)

	for {
		select {
		case <-ctx.Done():
			n.waitForShutdown()
			return
		case <-ticker.C:
			rand := n.rand.Intn(n.cfg.High-n.cfg.Low+1) + n.cfg.Low
			content := &pb.RngMessage{
				Random: int64(rand),
			}
			msgBytes, err := utils.MarshallProtobuf(content)
			if err != nil {
				n.logger.Error("Error marshalling RNG message",
					log.Int("random_value", rand),
					log.Error(err),
				)
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
				n.logger.Error("Error publishing RNG message",
					log.String("topic", topic),
					log.Int("random_value", rand),
					log.Error(err),
				)
				continue
			}
			n.count++
			n.logger.Info("Published random number",
				log.Int("random_value", rand),
				log.String("topic", topic),
				log.Int64("message_count", n.count),
			)
		}
	}
}

func (n *RNGNode) waitForShutdown() {
	n.logger.Info("RNG node shutdown complete",
		log.String("node_name", n.Name()),
		log.Int64("total_messages", n.count),
	)
}
