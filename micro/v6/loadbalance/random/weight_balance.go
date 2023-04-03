package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
)

const WeightRandom = "WEIGHT_RANDOM"

type WeightBalancer struct {
	totalWeight uint32
	conns       []*weightConn
}

func (b *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	val := rand.Intn(int(b.totalWeight))

	for _, con := range b.conns {
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
}

func (b *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	var totalWeight uint32 = 0
	for con, conInfo := range info.ReadySCs {
		// 这里你可以考虑容错，例如服务器没有配置权重，给一个默认的权重
		// 但是我认为这种容错会让用户不经意间出 BUG，所以我这里不会校验，而是直接让它 panic
		// 这是因为 gRPC 确实没有设计 error 作为返回值
		weight := conInfo.Address.Attributes.Value("weight").(int)
		totalWeight += uint32(weight)
		conns = append(conns, &weightConn{
			SubConn: con,
			weight:  uint32(weight),
		})
	}
	return &WeightBalancer{
		totalWeight: totalWeight,
		conns:       conns,
	}
}

type weightConn struct {
	// 初始权重
	weight uint32
	balancer.SubConn
}
