package main

import (
	"encoding/json"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro"
	"gitee.com/geektime-geekbang/geektime-go/micro/example/proto/gen"
	"gitee.com/geektime-geekbang/geektime-go/micro/loadbalance/roundrobin"
	"gitee.com/geektime-geekbang/geektime-go/micro/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"time"
)

func main() {
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
	pickerBuilder := &roundrobin.Builder{}
	client := micro.NewClient(micro.ClientWithInsecure(),
		micro.ClientWithRegistry(r, time.Second*3),
		micro.ClientWithPickerBuilder(pickerBuilder.Name(), pickerBuilder))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 设置初始化连接的 ctx
	conn, err := client.Dial(ctx, "user-service")
	cancel()
	if err != nil {
		panic(err)
	}
	userClient := gen.NewUserServiceClient(conn)
	for i := 0; i < 10; i++ {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
		resp, err := userClient.GetById(ctx, &gen.GetByIdReq{
			Id: 12,
		})
		if err != nil {
			panic(err)
		}
		cancel()
		bs, err := json.Marshal(resp.User)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bs))
	}
}
