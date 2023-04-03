package homework2_fastest

import (
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/loadbalance"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"sync"
	"time"
)

type Balancer struct {
	mutex    sync.RWMutex
	conns    []*conn
	filter   loadbalance.Filter
	lastSync time.Time
	endpoint string
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	b.mutex.RLock()
	// 执行过滤，并且挑出响应时间最短的那个节点
	b.mutex.RUnlock()

	return balancer.PickResult{
	}, nil
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	// 构造链接
	conns := make([]*conn, 0, len(info.ReadySCs))
	// 构造 flt
	flt := b.Filter
	res := &Balancer{
		conns:  conns,
		filter: flt,
	}

	// 在这里考虑定时刷新响应时间，也就是从 prometheus 里面查询
	// 基本上就是启动计时器，调用 updateRespTime
	return res
}

func (b *Balancer) updateRespTime(endpoint, query string) {
	// 这里很难容错，即如果刷新响应时间失败该怎么办
	httpResp, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=%s", endpoint, query))
	if err != nil {
		// 这里难处理，可以考虑记录错误，然后等下一次
		// 可以考虑中断
		// 也可以重试一定次数之后中断
		log.Fatalln("查询 prometheus 失败", err)
		return
	}
	fmt.Println(httpResp)
	// 解析你从 prometheus 里面拿到的结果
}

type Builder struct {
	Filter loadbalance.Filter
	// prometheus 的地址
	Endpoint string
	Query    string
	// 刷新响应时间的间隔
	Interval time.Duration
}

type conn struct {
	balancer.SubConn
	address resolver.Address
	// 响应时间
	response time.Duration
}

