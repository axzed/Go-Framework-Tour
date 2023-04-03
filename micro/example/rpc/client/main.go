package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/example/rpc/proto/gen"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc"
)

func main() {
	c, err := rpc.NewClient("0.0.0.0:8081")
	if err != nil {
		panic(err)
	}
	us := &UserService{}
	err = c.InitService(us)
	usOneway := &UserService{}
	err = c.InitService(usOneway)
	if err != nil {
		panic(err)
	}
	resp, err := us.GetById(context.Background(), &FindByUserIdReq{
		Id: 12,
	})
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(resp)
	fmt.Printf("收到响应: %s \n", data)

	_, err = us.AlwaysError(context.Background(), &FindByUserIdReq{
		Id: 12,
	})

	fmt.Printf("收到错误信息: %s \n", err.Error())

	usProto := &UserServiceProto{}
	err = c.InitService(usProto)
	presp, err := usProto.GetById(context.Background(), &gen.GetByIdReq{})
	if err != nil {
		panic(err)
	}
	data, _ = json.Marshal(presp)
	fmt.Printf("收到响应: %s \n", data)
}

// UserService 声明服务，反射会把 GetById 转成一个 RPC 调用
type UserService struct {
	// 需要注意，因为我们没有办法修改方法的实现，所以我们只能考虑使用这种形态
	GetById func(ctx context.Context, req *FindByUserIdReq) (*FindByUserIdResp, error)

	AlwaysError func(ctx context.Context, req *FindByUserIdReq) (*FindByUserIdResp, error)
}

func (u *UserService) ServiceName() string {
	return "user"
}

// UserServiceProto 声明服务，反射会把 GetById 转成一个 RPC 调用
type UserServiceProto struct {
	// 需要注意，因为我们没有办法修改方法的实现，所以我们只能考虑使用这种形态
	GetById func(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error)
}

func (u *UserServiceProto) ServiceName() string {
	return "user-service-proto"
}
