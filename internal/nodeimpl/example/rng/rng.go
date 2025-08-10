package rng

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/BullionBear/sequex/internal/model"
	errpb "github.com/BullionBear/sequex/internal/model/protobuf/error"
	rngpb "github.com/BullionBear/sequex/internal/model/protobuf/example/rng"
	"github.com/BullionBear/sequex/internal/nodeimpl/utils"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type RNGConfig struct {
	low      int     `json:"low"`
	high     int     `json:"high"`
	interval float64 `json:"interval"`
}

func (c *RNGConfig) Low() int {
	return c.low
}

func (c *RNGConfig) High() int {
	return c.high
}

func (c *RNGConfig) Interval() time.Duration {
	return time.Duration(c.interval * float64(time.Second))
}

type RNGNode struct {
	*node.BaseNode
	// Configurable parameters
	cfg RNGConfig

	// State variables
	rand  *rand.Rand
	mutex sync.Mutex
	nData int64
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
		nData:    0,
	}, nil
}

func (n *RNGNode) OnMessage(msg *nats.Msg) {
}

func (n *RNGNode) OnRPC(req *nats.Msg) *nats.Msg {
	contentType := req.Header.Get("Content-Type")
	messageType := req.Header.Get("Message-Type")

	switch {
	case contentType == "application/protobuf" && messageType == "rng.RngCountRequest":
		var content rngpb.RngCountRequest
		if err := proto.Unmarshal(req.Data, &content); err != nil {
			log.Printf("Error unmarshalling RngCountRequest: %v", err)
			return utils.MakeErrorMessage(utils.ErrorCodeProtobufDeserialization, err)
		}
		return n.onRngCountRequest(&content)
	case contentType == "application/json" && messageType == "Config":
		var content RNGConfig
		if err := json.Unmarshal(req.Data, &content); err != nil {
			log.Printf("Error unmarshalling Config: %v", err)
			return utils.MakeErrorMessage(utils.ErrorCodeJSONDeserialization, err)
		}
		n.rpcMethodConfig(req)
	}
	return nil
}

func (n *RNGNode) onRngCountRequest(req *rngpb.RngCountRequest) *nats.Msg {
	response := &rngpb.RngCountResponse{
		Count: n.nData,
	}
	responseBytes, err := model.MarshallProtobuf(response)
	if err != nil {
		return utils.MakeErrorMessage(utils.ErrorProtobufSerialization, err)
	}
	msg := nats.NewMsg(req.Reply)
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Message-Type", "rng.RngCountResponse")
	msg.Data = responseBytes
	return msg
}

func (n *RNGNode) publishRng() {
	ticker := time.NewTicker(n.cfg.Interval())
	defer ticker.Stop()
	for {
		select {
		case <-n.Context().Done():
			n.WaitForShutdown()
			return
		case <-ticker.C:
			rand := n.rand.IntN(n.cfg.High()-n.cfg.Low()+1) + n.cfg.Low()
			content := &rngpb.RngMessage{
				Random: int64(rand),
			}
			msgBytes, err := model.MarshallProtobuf(content)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}
			if err := n.GetNATSConnection().Publish(n.cfg.PubSubject(), msgBytes); err != nil {
				log.Printf("Error publishing message: %v", err)
				continue
			}
			n.nData++
			log.Printf("Published random number: %d", rand)
		}
	}
}

func (n *RNGNode) WaitForShutdown() {
	return
}

func (n *RNGNode) rpcMethods(m *nats.Msg) {
	contentType := m.Header.Get("Content-Type")
	messageType := m.Header.Get("Message-Type")

	switch {
	case contentType == "application/protobuf" && messageType == "rng.RngCountRequest":
		n.rpcMethodRngCountRequest(m)
	case contentType == "application/json" && messageType == "Config":
		n.rpcMethodConfig(m)
	}
}

func (n *RNGNode) rpcMethodRngCountRequest(m *nats.Msg) {
	var content rngpb.RngCountRequest
	if err := proto.Unmarshal(m.Data, &content); err != nil {
		log.Printf("Error unmarshalling RngCountRequest: %v", err)
		response := &errpb.ErrorResponse{
			Code:    -1,
			Message: "Error unmarshalling RngCountRequest",
		}
		responseBytes, _ := model.MarshallProtobuf(response)
		responseMsg := nats.NewMsg(m.Reply)
		responseMsg.Header.Set("Content-Type", "application/protobuf")
		responseMsg.Header.Set("Message-Type", "error.ErrorResponse")
		responseMsg.Data = responseBytes
		if err := m.RespondMsg(responseMsg); err != nil {
			log.Printf("Error responding to RngCountRequest: %v", err)
		}
		return
	}
	response := &rngpb.RngCountResponse{
		Count: n.nData,
	}
	responseBytes, _ := model.MarshallProtobuf(response)
	responseMsg := nats.NewMsg(m.Reply)
	m.RespondMsg(responseMsg)
}

func (n *RNGNode) rpcMethodConfig(m *nats.Msg) {
	responseBytes, err := json.Marshal(n.cfg)
	if err != nil {
		log.Printf("Error marshalling config: %v", err)
		return
	}
	responseMsg := nats.NewMsg(m.Reply)
	responseMsg.Data = responseBytes
	responseMsg.Header.Set("Content-Type", "application/json")
	responseMsg.Header.Set("Message-Type", "rng.RngConfig")
	if err := m.RespondMsg(responseMsg); err != nil {
		log.Printf("Error responding to Config: %v", err)
	}
}
