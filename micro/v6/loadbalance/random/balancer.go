package random

import (
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Random = "RANDOM"

type Balancer struct {
	length int
	conns  []balancer.SubConn
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	index := rand.Intn(b.length)
	return balancer.PickResult{
		SubConn: b.conns[index],
	}, nil
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for con := range info.ReadySCs {
		conns = append(conns, con)
	}
	return &Balancer{
		conns:  conns,
		length: len(conns),
	}
}

type Builder struct {
}
