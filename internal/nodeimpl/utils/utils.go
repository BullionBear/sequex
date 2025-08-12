package utils

import (
	pb "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type ErrorCode int

const (
	ErrorProtobufDeserialization ErrorCode = 1000
	ErrorProtobufSerialization   ErrorCode = 1001
	ErrorJSONDeserialization     ErrorCode = 1002
	ErrorJSONSerialization       ErrorCode = 1003
	ErrorInternal                ErrorCode = 1004
	ErrorUnknownMessageType      ErrorCode = 1005
)

func (e ErrorCode) Int() int {
	return int(e)
}

func MakeErrorMessage(id int64, code ErrorCode, err error) *nats.Msg {
	contentType := "application/protobuf"
	messageType := "error"
	pbError := &pb.ErrorResponse{
		Message: err.Error(),
		Code:    int64(code),
		Id:      id,
	}
	data, err := MarshallProtobuf(pbError)
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

func MarshallProtobuf[T proto.Message](obj T) ([]byte, error) {
	return proto.Marshal(obj)
}

func UnmarshallProtobuf[T proto.Message](data []byte) (T, error) {
	var obj T
	if err := proto.Unmarshal(data, obj); err != nil {
		return obj, err
	}
	return obj, nil
}
