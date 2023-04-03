package leastactive

import (
	"gitee.com/geektime-geekbang/geektime-go/micro/loadbalance"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"math"
	"sync/atomic"
)

type Balancer struct {
	conns  []*conn
	filter loadbalance.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// 使用原子操作的弊端就是不够准确
	// 而如果改用锁，则性能太差
	// 想想这是为什么？
	var leastActive uint32 = math.MaxUint32
	var res *conn
	for _, c := range b.conns {
		if !b.filter(info, c.address) {
			continue
		}
		active := atomic.LoadUint32(&c.active)
		if active < leastActive {
			leastActive = active
			res = c
		}
	}

	atomic.AddUint32(&res.active, 1)
	return balancer.PickResult{
		SubConn: res.SubConn,
		Done: func(info balancer.DoneInfo) {
			atomic.AddUint32(&res.active, -1)
		},
	}, nil
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*conn, 0, len(info.ReadySCs))
	for con, val := range info.ReadySCs {
		conns = append(conns, &conn{
			SubConn: con,
			address: val.Address,
		})
	}
	flt := b.Filter
	if flt == nil {
		flt = func(info balancer.PickInfo, address resolver.Address) bool {
			return true
		}
	}
	return &Balancer{
		conns:  conns,
		filter: flt,
	}
}

type Builder struct {
	Filter loadbalance.Filter
}

type conn struct {
	balancer.SubConn
	active  uint32
	address resolver.Address
}
