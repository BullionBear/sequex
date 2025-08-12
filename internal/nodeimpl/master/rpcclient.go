package master

import (
	"time"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type RequestType int

const (
	RequestTypeConfig RequestType = iota
	RequestTypeHealth
	RequestTypeVarz
)

type MasterRPCClient struct {
	nc *nats.Conn
}

func NewMasterRPCClient(nc *nats.Conn) *MasterRPCClient {
	return &MasterRPCClient{nc: nc}
}

func (c *MasterRPCClient) Request(subject string, reqType RequestType) (*nats.Msg, error) {
	msg := nats.Msg{
		Subject: subject,
		Header: map[string][]string{
			"Content-Type": {"application/protobuf"},
			"Message-Type": {"common.EmptyRequest"},
		},
	}
	contentBytes, err := proto.Marshal(&pbCommon.EmptyRequest{
		Id:   time.Now().UnixNano(),
		Type: parseRequestType(reqType),
	})
	if err != nil {
		return nil, err
	}
	msg.Data = contentBytes
	return c.nc.RequestMsg(&msg, 10*time.Second)
}

func parseRequestType(reqType RequestType) pbCommon.RequestType {
	switch reqType {
	case RequestTypeConfig:
		return pbCommon.RequestType_REQUEST_TYPE_CONFIG
	case RequestTypeHealth:
		return pbCommon.RequestType_REQUEST_TYPE_HEALTH
	case RequestTypeVarz:
		return pbCommon.RequestType_REQUEST_TYPE_VARZ
	default:
		return pbCommon.RequestType_REQUEST_TYPE_UNDEFINED
	}
}
