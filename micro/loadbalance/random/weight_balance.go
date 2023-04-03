package random

import (
	"gitee.com/geektime-geekbang/geektime-go/micro/loadbalance"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"math/rand"
)

const WeightRandom = "WEIGHT_RANDOM"

type WeightBalancer struct {
	conns  []*weightConn
	filter loadbalance.Filter
}

func (b *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var totalWeight uint32
	for _, con := range b.conns {
		// 加上去之后，
		if !b.filter(info, con.address) {
			continue
		}
		totalWeight += con.weight
	}
	val := rand.Intn(int(totalWeight))
	for _, con := range b.conns {
		// 加上去之后，
		if !b.filter(info, con.address) {
			continue
		}
		val = val - int(con.weight)
		if val <= 0 {
			return balancer.PickResult{
				SubConn: con.SubConn,
				Done: func(info balancer.DoneInfo) {
					// 实际上在这里我们也可以考虑根据调用结果来调整权重，
					// 类似于 roubin 里面
				},
			}, nil
		}
	}
	// 实际上不可能运行到这里，因为在前面我们必然能找到一个值
	return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
}

type WeightBalancerBuilder struct {
	Filter loadbalance.Filter
}

func (b *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for con, conInfo := range info.ReadySCs {
		// 这里你可以考虑容错，例如服务器没有配置权重，给一个默认的权重
		// 但是我认为这种容错会让用户不经意间出 BUG，所以我这里不会校验，而是直接让它 panic
		// 这是因为 gRPC 确实没有设计 error 作为返回值
		weight := conInfo.Address.Attributes.Value("weight").(int)
		conns = append(conns, &weightConn{
			SubConn: con,
			weight:  uint32(weight),
			address: conInfo.Address,
		})
	}

	filter := b.Filter
	if filter == nil {
		filter = func(info balancer.PickInfo, address resolver.Address) bool {
			return true
		}
	}
	return &WeightBalancer{
		conns:  conns,
		filter: filter,
	}
}

type weightConn struct {
	// 初始权重
	weight uint32
	balancer.SubConn
	address resolver.Address
}
