package main

import (
	"context"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/demo"
	"gitee.com/geektime-geekbang/geektime-go/demo/cluster/broadcast"
	"gitee.com/geektime-geekbang/geektime-go/demo/example/loadbalance/proto/gen"
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

	cb := broadcast.NewClusterBuilder(r, "user-service")
	cc, err := grpc.Dial("registry:///user-service",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(cb.BuildUnary()),
		grpc.WithResolvers(demo.NewResolverBuilder(r, time.Second * 3)),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`,
			pickerBuilder.Name())))
	if err != nil {
		panic(err)
	}
	client := gen.NewUserServiceClient(cc)
	for i := 0; i < 100; i++ {
		// 指定使用广播
		ctx := context.Background()
		ctx = broadcast.UsingBroadCast(ctx)

		resp, err := client.GetById(ctx, &gen.GetByIdReq{})
		if err != nil {
			panic(err)
		}
		log.Println(resp.User)
	}
}
