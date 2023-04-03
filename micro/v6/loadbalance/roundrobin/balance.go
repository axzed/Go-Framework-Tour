package roundrobin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

const RoundRobin = "ROUND_ROBIN"

type Balancer struct {
	mutex  sync.Mutex
	cnt    uint32
	length uint32
	conns  []balancer.SubConn
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// 使用原子操作而不是锁，理论上是可行的，但是最终效果就不是一个严格的轮询，而是一个大致的轮询
	// 这种情况下，为什么不直接使用随机呢？
	// cnt := atomic.AddUint32(&b.cnt, 1)
	// index := cnt % b.length
	// atomic.StoreUint32(&b.cnt, index)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	index := b.cnt % b.length
	b.cnt = index + 1
	return balancer.PickResult{
		SubConn: b.conns[index],
		Done: func(info balancer.DoneInfo) {
			// 实际上，这里你是要考虑如果调用失败，
			// 会不会是客户端和服务端的网络不通，
			// 按照道理来说，是需要将这个连不通的节点删除的
			// 但是删除之后又要考虑一段时间之后加回来
		},
	}, nil
}

type Builder struct {
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for con := range info.ReadySCs {
		conns = append(conns, con)
	}
	return &Balancer{
		conns:  conns,
		length: uint32(len(conns)),
	}
}

func (b *Builder) Name() string {
	return RoundRobin
}
