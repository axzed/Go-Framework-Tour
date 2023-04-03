package rpc

import (
	"gitee.com/geektime-geekbang/geektime-go/micro/v6/registry"
	"golang.org/x/net/context"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

type grpcResolverBuilder struct {
	r       registry.Registry
	timeout time.Duration
}

func newResolverBuilder(r registry.Registry, timeout time.Duration) *grpcResolverBuilder {
	return &grpcResolverBuilder{
		r:       r,
		timeout: timeout,
	}
}

func (r *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	res := &grpcResolver{
		target:   target,
		cc:       cc,
		registry: r.r,
		close:    make(chan struct{}, 1),
		timeout:  r.timeout,
	}
	res.resolve()
	go res.watch()
	return res, nil
}

// 伪代码
// func (r *grpcResolverBuilder) Build1(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
// 	res := &grpcResolver{}
// 	// 要在这里初始化连接，然后更新 cc 里面的连接信息
// 	address := make([]resolver.Address, 10)
// 	cc.UpdateState(resolver.State{
// 		Addresses: address,
// 	})
// 	return res, nil
// }

// Scheme 返回一个固定的值，registry 代表的是我们设计的注册中心
func (r *grpcResolverBuilder) Scheme() string {
	return "registry"
}

type grpcResolver struct {
	target   resolver.Target
	cc       resolver.ClientConn
	registry registry.Registry
	close    chan struct{}
	timeout  time.Duration
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	// 重新获取一下所有服务
	g.resolve()
}

func (g *grpcResolver) Close() {
	g.close <- struct{}{}
}

func (g *grpcResolver) watch() {
	events := g.registry.Subscribe(g.target.Endpoint)
	for {
		select {
		case <-events:
			// 一种做法就是我们这边区别处理不同事件类型，然后更新数据
			// switch event.Type {
			//
			//			}
			// 另外一种做法就是我们这里采用的，每次事件发生的时候，就直接刷新整个可用服务列表
			g.resolve()

		case <-g.close:
			return
		}
	}
}

func (g *grpcResolver) resolve() {
	serviceName := g.target.Endpoint
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	instances, err := g.registry.ListServices(ctx, serviceName)
	cancel()
	if err != nil {
		g.cc.ReportError(err)
	}

	address := make([]resolver.Address, 0, len(instances))
	for _, ins := range instances {
		address = append(address, resolver.Address{
			Addr:       ins.Address,
			ServerName: ins.Name,
			Attributes: attributes.New("weight", ins.Weight),
		})
	}
	err = g.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	if err != nil {
		g.cc.ReportError(err)
	}
}
