package broadcastall

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/registry"
	"google.golang.org/grpc"
	"reflect"
	"sync"
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

func UsingBroadCast(ctx context.Context) (context.Context, <-chan Resp) {
	ch := make(chan Resp)
	return context.WithValue(ctx, key{}, ch), ch
}

// Resp 没有办法用泛型
type Resp struct {
	Val any
	Err error
}

func isBroadCast(ctx context.Context) (bool, chan Resp) {
	val := ctx.Value(key{})
	if val != nil {
		res, ok := val.(chan Resp)
		if ok {
			return ok, res
		}
	}
	return false, nil
}

func (b ClusterBuilder) BuildUnary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ok, ch := isBroadCast(ctx)
		if !ok {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		ins, err := b.registry.ListServices(ctx, b.service)
		if err != nil {
			ch <- Resp{Err: err}
			return nil
		}
		typ := reflect.TypeOf(reply).Elem()
		var wg sync.WaitGroup
		wg.Add(len(ins))
		for _, instance := range ins {
			in := instance
			go func() {
				conn, er := grpc.Dial(in.Address, b.opts...)
				if er != nil {
					ch <- Resp{Err: er}
					return
				}
				r := reflect.New(typ)
				val := r.Interface()
				err = conn.Invoke(ctx, method, req, val, opts...)
				// 这种写法的风险在于，如果用户没有接收响应，
				// 那么这里会阻塞导致 goroutine 泄露
				ch <- Resp{Err: err, Val: val}
				wg.Done()
			}()
		}
		go func() {
			wg.Wait()
			// 要记得 close 掉，不然用户不知道还有没有数据
			// 用户在调用的时候是不知道有多少个实例还存活着
			close(ch)
		}()
		return nil
	}
}
