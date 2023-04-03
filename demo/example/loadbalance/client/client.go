package main

import (
	"context"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/demo"
	gen2 "gitee.com/geektime-geekbang/geektime-go/demo/example/loadbalance/proto/gen"
	"gitee.com/geektime-geekbang/geektime-go/demo/loadbalance"
	"gitee.com/geektime-geekbang/geektime-go/demo/loadbalance/roundrobin"
	"gitee.com/geektime-geekbang/geektime-go/demo/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"log"
	"time"
)

func main() {
	// 注册中心
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		panic(err)
	}
	r, err := etcd.NewRegistry(etcdClient)
	if err != nil {
		panic(err)
	}
	// 注册你的负载均衡策略
	pickerBuilder := &roundrobin.PickerBuilder{
		Filter: loadbalance.GroupFilter,
	}
	builder := base.NewBalancerBuilder(pickerBuilder.Name(), pickerBuilder, base.Config{HealthCheck: true})
	balancer.Register(builder)

	cc, err := grpc.Dial("registry:///user-service",
		grpc.WithInsecure(),
		grpc.WithResolvers(demo.NewResolverBuilder(r, time.Second * 3)),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`,
			pickerBuilder.Name())))
	if err != nil {
		panic(err)
	}
	client := gen2.NewUserServiceClient(cc)
	for i := 0; i < 100; i++ {
		ctx := context.WithValue(context.Background(), "group", "b")
		resp, err := client.GetById(ctx, &gen2.GetByIdReq{})
		if err != nil {
			panic(err)
		}
		log.Println(resp.User)
	}
}
