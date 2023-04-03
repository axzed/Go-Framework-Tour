package leastactive

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
	"sync/atomic"
)

type Picker struct {
	mutex sync.Mutex
	conns []conn
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	res := p.conns[0]
	for i := 1; i < len(p.conns); i++ {
		if res.numReq > p.conns[i].numReq {
			res = p.conns[i]
		}
	}
	res.numReq ++
	return balancer.PickResult{SubConn: res.sub, Done: func(info balancer.DoneInfo) {
		atomic.AddUint64(&res.numReq, -1)
	}}, nil
}

type PickerBuilder struct {

}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]conn, 0, len(info.ReadySCs))
	for subCon := range info.ReadySCs {
		conns = append(conns, conn{sub: subCon, numReq: 0})
	}
	return &Picker{
		conns: conns,
	}
}

type conn struct {
	numReq uint64
	sub balancer.SubConn
}