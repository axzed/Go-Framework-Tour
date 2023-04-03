package main

import (
	"context"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/example/proto/gen"
)

type UserService struct {
	name string
	gen.UnimplementedUserServiceServer
}

func (u *UserService) GetById(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error) {
	fmt.Printf("server %s, get user id: %d \n", u.name, req.Id)
	return &gen.GetByIdResp{
		User: &gen.User{
			Id:     req.Id,
			Status: 123,
		},
	}, nil
}
