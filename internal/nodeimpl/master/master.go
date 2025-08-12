package master

import (
	"encoding/json"
	"fmt"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/BullionBear/sequex/internal/nodeimpl/utils"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type MasterConfig struct {
	TargetNodes []string `json:"target_nodes"`
}

type MasterNode struct {
	*node.BaseNode
	cfg MasterConfig
}

func init() {
	node.RegisterNode("master", NewMasterNode)
}

func NewMasterNode(name string, nc *nats.Conn, config node.NodeConfig, logger *log.Logger) (node.Node, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %v", err)
	}
	var cfg MasterConfig
	if err := json.Unmarshal(jsonBytes, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return &MasterNode{
		BaseNode: node.NewBaseNode(name, nc, *logger),
		cfg:      cfg,
	}, nil
}

func (n *MasterNode) Start() error {
	return nil
}

func (n *MasterNode) Shutdown() error {
	return nil
}

func (n *MasterNode) OnMessage(msg *nats.Msg) {
}

func (n *MasterNode) OnRPC(req *nats.Msg) *nats.Msg {
	contentType := req.Header.Get("Content-Type")
	messageType := req.Header.Get("Message-Type")

	n.Logger().Debug("Received RPC request",
		log.String("content_type", contentType),
		log.String("message_type", messageType),
	)
	switch {
	case contentType == "application/protobuf" && messageType == "common.EmptyRequest":
		var content pbCommon.EmptyRequest
		if err := proto.Unmarshal(req.Data, &content); err != nil {
			n.Logger().Error("Error unmarshalling EmptyRequest",
				log.Error(err),
			)
			return utils.MakeErrorMessage(0, utils.ErrorProtobufDeserialization, err)
		}
		switch content.Type {
		case pbCommon.RequestType_REQUEST_TYPE_CONFIG:
			return n.onConfig(&content)
		default:
			n.Logger().Warn("Unknown request type",
				log.Int("request_type", int(content.Type)),
			)
			return utils.MakeErrorMessage(content.Id, utils.ErrorUnknownMessageType, fmt.Errorf("unknown message type: %s", content.Type.String()))
		}
	default:
		n.Logger().Warn("Unknown message type",
			log.String("content_type", contentType),
			log.String("message_type", messageType),
		)
		return utils.MakeErrorMessage(0, utils.ErrorUnknownMessageType, fmt.Errorf("unknown message type: %s", messageType))
	}
}

func (n *MasterNode) onConfig(content *pbCommon.EmptyRequest) *nats.Msg {
	// Return a simple success response
	responseBytes, err := json.Marshal(&n.cfg)
	if err != nil {
		n.Logger().Error("Error marshalling config response",
			log.Error(err),
		)
		return utils.MakeErrorMessage(content.Id, utils.ErrorJSONSerialization, err)
	}
	msg := nats.Msg{}
	msg.Header.Set("Content-Type", "application/json")
	msg.Header.Set("Message-Type", "MasterConfig")
	msg.Data = responseBytes
	return &msg
}
