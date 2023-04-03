package random

import (
	"gitee.com/geektime-geekbang/geektime-go/micro/loadbalance"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const Random = "RANDOM"

type Balancer struct {
	conns  []conn
	filter loadbalance.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	candidates := make([]conn, 0, len(b.conns))
	for _, c := range b.conns {
		if !b.filter(info, c.address) {
			continue
		}
		candidates = append(candidates, c)
	}
	index := rand.Intn(len(candidates))
	return balancer.PickResult{
		SubConn: candidates[index].SubConn,
	}, nil
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]conn, 0, len(info.ReadySCs))
	for con, val := range info.ReadySCs {
		conns = append(conns, conn{
			SubConn: con,
			address: val.Address,
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

type Builder struct {
	Filter loadbalance.Filter
}

type conn struct {
	balancer.SubConn
	address resolver.Address
}
