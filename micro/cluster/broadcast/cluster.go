package broadcast

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/registry"
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
		ins, err := b.registry.ListServices(ctx, b.service)
		if err != nil {
			return err
		}
		var eg errgroup.Group
		for _, instance := range ins {
			in := instance
			// 转变为直连。这里我们预期 Address 是一个真的地址，例如 IP + 端口
			eg.Go(func() error {

				// 可怕的是我们每次进来都需要重新连，除非我们考虑缓存的问题
				// 缓存的问题则在于，我们需要管理它，在必要的时候关掉 conn
				conn, er := grpc.Dial(in.Address, b.opts...)
				if er != nil {
					return er
				}
				// 这里你可以考虑设计接口，允许用户把所有广播响应都拿到
				return conn.Invoke(ctx, method, req, reply, opts...)
			})
		}
		// 这种做法
		// 返回 error，则是第一个返回响应的 error
		// 返回 没有返回，那么 reply 将会是最后一个返回的值
		// 所以实际上存在覆盖的可能性
		return eg.Wait()
	}
}
