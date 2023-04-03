package roundrobin

import (
	"gitee.com/geektime-geekbang/geektime-go/demo/loadbalance"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"sync/atomic"
)

type Picker struct {
	ins []instance
	cnt uint64
	filter loadbalance.Filter
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	candidates := make([]instance, 0, len(p.ins))
	for _, sub := range p.ins {
		if !p.filter(info, sub.address) {
			continue
		}
		candidates = append(candidates, sub)
	}
	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	cnt := atomic.AddUint64(&p.cnt, 1)
	index := cnt % uint64(len(candidates))
	return balancer.PickResult{
		SubConn: candidates[index].sub,
		// 用来设计反馈式的负载均衡策略
		Done: func(info balancer.DoneInfo) {
			// 可以打个不健康标签
			// if info.Err !=nil {
			//
			// }
			// 这个地方就是很神奇的地方
			// 效果就是根据调用结果来调整你的负载均衡策略

			// 假如说出错了
			// if info.Err != nil {
				// 尝试把这个 subConn 置为不可用，或者临时移除出去
			// }
		},
	}, nil
}

type PickerBuilder struct {
	Filter loadbalance.Filter
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]instance, 0, len(info.ReadySCs))
	for sub, subInfo := range info.ReadySCs {
		conns = append(conns, instance{
			sub: sub,
			address: subInfo.Address,
		})
	}
	return &Picker{
		filter: p.Filter,
		ins: conns,
	}
}

func (p *PickerBuilder) Name() string {
	return "ROUND_ROBIN"
}

type instance struct {
	sub balancer.SubConn
	address resolver.Address
}

