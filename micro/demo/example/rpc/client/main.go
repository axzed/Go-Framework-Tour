package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/example/rpc/proto/gen"
)

func main() {
	c, err := demo.NewClient("0.0.0.0:8081")
	if err != nil {
		panic(err)
	}
	us := &UserService{}
	err = c.InitService(us)
	if err != nil {
		panic(err)
	}
	resp, err := us.GetById(context.Background(), &FindByUserIdReq{
		Id: 12,
	})
	if err != nil {
		panic(err)
	}

	_, _  = us.GetById(demo.CtxWithOneway(context.Background()), &FindByUserIdReq{
		Id: 12,
	})

	// 异步调用
	// go func() {
	// 	resp, err = us.GetById(context.Background(), &FindByUserIdReq{
	// 		Id: 12,
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	// 处理一下 resp，就是回调
	// }()

	// 用户组织的单向调用
	// go func() {
	// 	us.GetById(context.Background(), &FindByUserIdReq{
	// 		Id: 12,
	// 	})
		// 处理一下 resp，就是回调
	// }()

	// var eg errgroup.Group
	// eg.Go(func() error {
	//
	// 	resp, err := us.GetById(context.Background(), &FindByUserIdReq{
	// 		Id: 12,
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// })
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

func (u *UserService) Name() string {
	return "user"
}

// UserServiceProto 声明服务，反射会把 GetById 转成一个 RPC 调用
type UserServiceProto struct {
	// 需要注意，因为我们没有办法修改方法的实现，所以我们只能考虑使用这种形态
	GetById func(ctx context.Context, req *gen.GetByIdReq) (*gen.GetByIdResp, error)
}

func (u *UserServiceProto) Name() string {
	return "user-service-proto"
}
