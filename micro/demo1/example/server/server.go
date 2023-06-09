package main

import (
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo1"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo1/example/proto/gen"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo1/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"strconv"
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
	var eg errgroup.Group

	for i := 0; i < 3; i++ {
		idx := i
		eg.Go(func() error {
			server := demo1.NewServer("user-service", demo1.ServerWithRegistry(r))
			us := &UserService{
				name: fmt.Sprintf("server-%d", idx),
			}
			gen.RegisterUserServiceServer(server, us)
			fmt.Println("启动服务器: " + us.name)
			// 端口分别是 8081, 8082, 8083
			return server.Start(":" + strconv.Itoa(8081+idx))
		})
	}
	// 正常或者异常退出都会返回
	err = eg.Wait()
}
