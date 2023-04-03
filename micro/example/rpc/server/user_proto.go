package main

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/example/rpc/proto/gen"
)

// UserServiceProto 用来测试 protobuf 协议
type UserServiceProto struct {
}

func (u *UserServiceProto) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	return &gen.GetByIdResp{
		User: &gen.User{
			Id: 123,
		},
	}, nil
}

func (u UserServiceProto) ServiceName() string {
	return "user-service-proto"
}
