package broadcast

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo/registry"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type ClusterBuilder struct {
	registry registry.Registry
	service  string
	opts     []grpc.DialOption

	// 还可以考虑设计成这种样子，然后和注册中心解耦
	// 不过在一个框架内部，耦合也没啥关系
	//	findServes func(ctx) []ServiceInstance
}

func NewClusterBuilder(r registry.Registry, service string, dialOpts ...grpc.DialOption) *ClusterBuilder {
	return &ClusterBuilder{
		registry: r,
		service:  service,
		opts:     dialOpts,
	}
}

type key struct{}

func UsingBroadCast(ctx context.Context) context.Context {
	return context.WithValue(ctx, key{}, true)
}

func isBroadCast(ctx context.Context) bool {
	val := ctx.Value(key{})
	if val != nil {
		res, ok := val.(bool)
		return ok && res
	}
	return false
}

func (b ClusterBuilder) BuildUnary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !isBroadCast(ctx) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		// 怎么广播?
		instances, err := b.registry.ListService(ctx, b.service)

		// 在这里之前，引入一些过滤的逻辑，就可以近似实现组播
		if err != nil {
			return err
		}
		// 遍历所有的实例
		var eg errgroup.Group
		for _, instance := range instances {
			ins := instance
			eg.Go(func() error {
				insCC, err := grpc.Dial(ins.Address, b.opts...)
				if err != nil {
					return err
				}
				err = invoker(ctx, method, req, reply, insCC, opts...)
				return err
			})
		}
		return eg.Wait()
	}
}
