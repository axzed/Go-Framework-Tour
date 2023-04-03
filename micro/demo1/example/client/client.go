package main

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo1"
	gen2 "gitee.com/geektime-geekbang/geektime-go/micro/demo1/example/proto/gen"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo1/registry"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var r registry.Registry
	rsBuilder := demo1.NewResolverBuilder(r, time.Second)
	cc, err := grpc.Dial("registry:///user-service",
		grpc.WithInsecure(),
		grpc.WithResolvers(rsBuilder))
	if err != nil {
		panic(err)
	}
	usClient := gen2.NewUserServiceClient(cc)
	resp, err := usClient.GetById(context.Background(), &gen2.GetByIdReq{
		Id: 12,
	})
	if err != nil {
		panic(err)
	}
	log.Println(resp)
}
