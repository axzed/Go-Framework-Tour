package main

import (
	"context"
	"fmt"
	gen2 "gitee.com/geektime-geekbang/geektime-go/micro/demo1/example/proto/gen"
)

type UserService struct {
	name string
	gen2.UnimplementedUserServiceServer
}

func (u *UserService) GetById(ctx context.Context, req *gen2.GetByIdReq) (*gen2.GetByIdResp, error) {
	fmt.Printf("server %s, get user id: %d \n", u.name, req.Id)
	return &gen2.GetByIdResp{
		User: &gen2.User{
			Id:     req.Id,
			Status: 123,
		},
	}, nil
}
