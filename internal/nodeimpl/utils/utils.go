package utils

import (
	"github.com/BullionBear/sequex/internal/model"
	pb "github.com/BullionBear/sequex/internal/model/protobuf/error"
	"github.com/nats-io/nats.go"
)

type ErrorCode int

const (
	ErrorProtobufDeserialization ErrorCode = 1000
	ErrorProtobufSerialization   ErrorCode = 1001
	ErrorJSONDeserialization     ErrorCode = 1002
	ErrorJSONSerialization       ErrorCode = 1003
	ErrorInternal                ErrorCode = 1004
)

func (e ErrorCode) Int() int {
	return int(e)
}

func MakeErrorMessage(code ErrorCode, err error) *nats.Msg {
	contentType := "application/protobuf"
	messageType := "error"
	pbError := &pb.ErrorResponse{
		Message: err.Error(),
		Code:    int64(code.Int()),
	}
	data, err := model.MarshallProtobuf(pbError)
	if err != nil {
		return nil
	}
	msg := &nats.Msg{
		Header: map[string][]string{
			"Content-Type": {contentType},
			"Message-Type": {messageType},
		},
		Data: data,
	}
	return msg
}
