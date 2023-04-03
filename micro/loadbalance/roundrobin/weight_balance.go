package roundrobin

import (
	"gitee.com/geektime-geekbang/geektime-go/micro/loadbalance"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"math"
	"sync"
	"sync/atomic"
)

const WeightRoundRobin = "WEIGHT_ROUND_ROBIN"

type WeightBalancer struct {
	mutex  sync.Mutex
	conns  []*weightConn
	filter loadbalance.Filter
}

func (b *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var totalWeight uint32
	var res *weightConn

	b.mutex.Lock()
	for _, node := range b.conns {
		if !b.filter(info, node.address) {
			continue
		}
		totalWeight += node.efficientWeight
		node.currentWeight += node.efficientWeight
		if res == nil || res.currentWeight < node.currentWeight {
			res = node
		}
	}
	res.currentWeight -= totalWeight
	b.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.SubConn,
		Done: func(info balancer.DoneInfo) {
			for {
				// 这里就是一个棘手的地方了
				// 按照算法，如果调用没有问题，那么增加权重
				// 如果调用有问题，减少权重

				// 直接减是很危险的事情，因为你可能 0 - 1 直接就最大值了
				// 也就是说一个节点不断失败不断失败，最终反而权重最大
				// 类似地，如果一个节点不断加不断加，最大值加1反而变最小值
				// if info.Err != nil {
				// 	atomic.AddUint32(&res.weight, -1)
				// } else {
				// 	atomic.AddUint32(&res.weight, 1)
				// }
				// 所以可以考虑 CAS 来，或者在 weightConn 里面设置一个锁
				weight := atomic.LoadUint32(&(res.efficientWeight))
				if info.Err != nil && weight == 0 {
					return
				}
				if info.Err == nil && weight == math.MaxUint32 {
					return
				}
				newWeight := weight
				if info.Err == nil {
					newWeight += 1
				} else {
					newWeight -= 1
				}
				if atomic.CompareAndSwapUint32(&(res.efficientWeight), weight, newWeight) {
					return
				}
			}
		},
	}, nil
}

type WeightBalancerBuilder struct {
}

func (b *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for con, conInfo := range info.ReadySCs {
		// 这里你可以考虑容错，例如服务器没有配置权重，给一个默认的权重
		// 但是我认为这种容错会让用户不经意间出 BUG，所以我这里不会校验，而是直接让它 panic
		// 这是因为 gRPC 确实没有设计 error 作为返回值
		weight := conInfo.Address.Attributes.Value("weight").(int)
		conns = append(conns, &weightConn{
			SubConn:         con,
			weight:          uint32(weight),
			efficientWeight: uint32(weight),
			currentWeight:   uint32(weight),
			address:         conInfo.Address,
		})
	}
	return &WeightBalancer{
		conns: conns,
	}
}

type weightConn struct {
	// 初始权重
	weight uint32
	// 当前权重
	currentWeight uint32
	// 有效权重，在整个过程中我们是会动态调整权重的
	efficientWeight uint32
	balancer.SubConn
	address resolver.Address
}
