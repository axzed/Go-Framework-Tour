package leastactive

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync/atomic"
)

type Balancer struct {
	conns []*conn
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
	for con := range info.ReadySCs {
		conns = append(conns, &conn{
			SubConn: con,
		})
	}
	return &Balancer{
		conns: conns,
	}
}

type Builder struct {
}

type conn struct {
	balancer.SubConn
	active uint32
}
