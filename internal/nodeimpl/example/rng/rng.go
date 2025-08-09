package rng

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/BullionBear/sequex/internal/model"
	rngpb "github.com/BullionBear/sequex/internal/model/protobuf/example/rng"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type RNGConfig struct {
	low        int     `json:"low"`
	high       int     `json:"high"`
	interval   float64 `json:"interval"`
	rpcSubject string  `json:"rpc_subject"`
	pubSubject string  `json:"pub_subject"`
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

func (c *RNGConfig) RpcSubject() string {
	return c.rpcSubject
}

func (c *RNGConfig) PubSubject() string {
	return c.pubSubject
}

type RNGNode struct {
	*node.BaseNode
	// Configurable parameters
	cfg RNGConfig

	// State variables
	rand     *rand.Rand
	mutex    sync.Mutex
	nSuccess int64
	nFailure int64
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
		BaseNode: node.NewBaseNode(name),
		cfg:      cfg,
		rand:     rand,
		nSuccess: 0,
		nFailure: 0,
		mutex:    sync.Mutex{},
	}, nil
}

func (n *RNGNode) Start() error {
	go n.publishRng()
	return nil
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
				n.nFailure++
				continue
			}
			if err := n.GetNATSConnection().Publish(n.cfg.PubSubject(), msgBytes); err != nil {
				log.Printf("Error publishing message: %v", err)
				n.nFailure++
				continue
			}
			n.nSuccess++
			log.Printf("Published random number: %d", rand)
		}
	}
}

func (n *RNGNode) Register(nc *nats.Conn) error {
	nc.Subscribe(n.cfg.RpcSubject(), func(m *nats.Msg) {
		contentType := m.Header.Get("Content-Type")
		messageType := m.Header.Get("Message-Type")

		switch contentType {
		case "application/protobuf":
			switch messageType {
			case "rng.RngCountRequest":
				var content rngpb.RngCountRequest
				if err := proto.Unmarshal(m.Data, &content); err != nil {
					log.Printf("Error unmarshalling RngCountRequest: %v", err)
					return
				}
				// Handle count request
				n.mutex.Lock()
				response := &rngpb.RngCountResponse{
					NSuccess: n.nSuccess,
					NFailure: n.nFailure,
				}
				n.mutex.Unlock()

				responseBytes, err := model.MarshallProtobuf(response)
				if err != nil {
					log.Printf("Error marshalling response: %v", err)
					return
				}

				if err := nc.Publish(m.Reply, responseBytes); err != nil {
					log.Printf("Error publishing response: %v", err)
				}

			case "rng.RngMessage":
				var content rngpb.RngMessage
				if err := proto.Unmarshal(m.Data, &content); err != nil {
					log.Printf("Error unmarshalling RngMessage: %v", err)
					return
				}
				n.mutex.Lock()
				log.Printf("Received random number: %d", content.GetRandom())
				n.mutex.Unlock()

			default:
				log.Printf("Unknown message type: %s", messageType)
			}
		default:
			log.Printf("Unknown content type: %s", contentType)
		}
	})
	return nil
}

func (n *RNGNode) WaitForShutdown() {
	n.BaseNode.WaitForShutdown()
}
