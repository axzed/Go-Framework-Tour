package proto

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

type Serializer struct {

}

func (s Serializer) Code() byte {
	return 1
}

func (s Serializer) Encode(val any) ([]byte, error) {
	msg, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("micro: 使用 proto 序列化协议必须使用 protoc 编译的类型")
	}
	return proto.Marshal(msg)
}

func (s Serializer) Decode(data []byte, val any) error {
	msg, ok := val.(proto.Message)
	if !ok {
		return errors.New("micro: 使用 proto 序列化协议必须使用 protoc 编译的类型")
	}
	return proto.Unmarshal(data, msg)
}

