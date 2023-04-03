package demo1

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo1/registry"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ServerOption func(server *Server)

type Server struct {
	name string
	*grpc.Server
	r        registry.Registry
	listener net.Listener
}

func NewServer(name string, opts...ServerOption) *Server {
	res := &Server{
		name: name,
		Server: grpc.NewServer(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithRegistry(r registry.Registry) ServerOption {
	return func(server *Server) {
		server.r = r
	}
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 这边开始注册
	// 一定是先启动端口再注册
	// 严格地来说，是服务都启动了，才注册
	if s.r != nil {
		// defer s.r.Unregister()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
		err = s.r.Register(ctx, registry.ServiceInstance{
			ServiceName: s.name,
			Address: listener.Addr().String(),
		})
		cancel()
		if err != nil {
			return err
		}
	}
	return s.Serve(listener)
}

func (s *Server) Close() error {
	// 可以在这里 Unregister
	// s.r.Unregister()
	// 这里可以插入你的优雅退出逻辑
	// s.listener.Close()
	return nil
}