package model

import (
	"google.golang.org/protobuf/proto"
)

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
