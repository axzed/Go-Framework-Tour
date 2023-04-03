package main

import (
	rpc "gitee.com/geektime-geekbang/geektime-go/micro/demo"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/serialize/json"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/serialize/proto"
)

func main() {
	svr := rpc.NewServer()
	svr.MustRegister(&UserService{})
	svr.MustRegister(&UserServiceProto{})
	svr.RegisterSerializer(json.Serializer{})
	svr.RegisterSerializer(proto.Serializer{})
	if err := svr.Start(":8081"); err != nil {
		panic(err)
	}
}
