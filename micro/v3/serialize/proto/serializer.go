package proto

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

type Serializer struct {
}

func (s Serializer) Code() byte {
	return 2
}

func (s Serializer) Encode(val any) ([]byte, error) {
	msg, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("serialize: 必须是 proto.Message 类型")
	}
	return proto.Marshal(msg)
}

func (s Serializer) Decode(data []byte, val any) error {
	msg, ok := val.(proto.Message)
	if !ok {
		return errors.New("serialize: 必须是 proto.Message 类型")
	}
	return proto.Unmarshal(data, msg)
}
