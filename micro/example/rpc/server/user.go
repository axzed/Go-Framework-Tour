package main

import (
	"context"
	"errors"
	"time"
)

type UserService struct {
}

func (u *UserService) GetById(ctx context.Context, req *FindByUserIdReq) (*FindByUserIdResp, error) {
	return &FindByUserIdResp{
		User: &User{
			Id:         12,
			Name:       "Tom",
			Avatar:     "http://my-avatar",
			Email:      "xxx@xxx.com",
			Password:   "123456",
			CreateTime: time.Now().Second(),
		},
	}, nil
}

func (u *UserService) AlwaysError(ctx context.Context, req *FindByUserIdReq) (*FindByUserIdResp, error) {
	return nil, errors.New("this is an error")
}

func (u *UserService) ServiceName() string {
	return "user"
}
