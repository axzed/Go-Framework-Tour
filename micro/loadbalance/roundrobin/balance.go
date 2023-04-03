package roundrobin

import (
	"gitee.com/geektime-geekbang/geektime-go/micro/loadbalance"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"sync"
)

const RoundRobin = "ROUND_ROBIN"

type Balancer struct {
	mutex  sync.Mutex
	cnt    uint32
	conns  []conn
	filter loadbalance.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {

	// 使用原子操作而不是锁，理论上是可行的，但是最终效果就不是一个严格的轮询，而是一个大致的轮询
	// 这种情况下，为什么不直接使用随机呢？
	// cnt := atomic.AddUint32(&b.cnt, 1)
	// index := cnt % b.length
	// atomic.StoreUint32(&b.cnt, index)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	candidates := make([]conn, 0, len(b.conns))
	for _, c := range b.conns {
		if !b.filter(info, c.address) {
			continue
		}
		candidates = append(candidates, c)
	}

	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	index := b.cnt % uint32(len(candidates))
	b.cnt = index + 1
	return balancer.PickResult{
		SubConn: candidates[index].SubConn,
		Done: func(info balancer.DoneInfo) {
			// 实际上，这里你是要考虑如果调用失败，
			// 会不会是客户端和服务端的网络不通，
			// 按照道理来说，是需要将这个连不通的节点删除的
			// 但是删除之后又要考虑一段时间之后加回来
		},
	}, nil
}

type Builder struct {
	Filter loadbalance.Filter
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]conn, 0, len(info.ReadySCs))
	for con, conInfo := range info.ReadySCs {
		conns = append(conns, conn{
			SubConn: con,
			address: conInfo.Address,
		})
	}
	filter := b.Filter
	if filter == nil {
		filter = func(info balancer.PickInfo, address resolver.Address) bool {
			return true
		}
	}
	return &Balancer{
		conns:  conns,
		filter: filter,
	}
}

func (b *Builder) Name() string {
	return RoundRobin
}

type conn struct {
	balancer.SubConn
	address resolver.Address
}
