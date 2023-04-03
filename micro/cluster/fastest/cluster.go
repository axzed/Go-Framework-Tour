package fastest

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/registry"
	"google.golang.org/grpc"
	"reflect"
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
		typ := reflect.TypeOf(reply).Elem()
		ch := make(chan resp)
		for _, instance := range ins {
			in := instance
			go func() {
				conn, er := grpc.Dial(in.Address, b.opts...)
				if er != nil {
					ch <- resp{err: er}
					return
				}
				r := reflect.New(typ)
				val := r.Interface()
				err = conn.Invoke(ctx, method, req, val, opts...)
				select {
				case ch <- resp{err: err, val: r}:
				default:
				}
			}()
		}
		select {
		case r := <-ch:
			if r.err == nil {
				reflect.ValueOf(reply).Elem().Set(r.val.Elem())
			}
			return r.err
		case <-ctx.Done():
			// 实际上这里是否监听 ctx 不重要，因为我们可以预期 grpc 会在超时的时候返回，走到上面的 error 分支
			return ctx.Err()
		}
	}
}

type resp struct {
	val reflect.Value
	err error
}
