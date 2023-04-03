package prometheus

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/observability"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"time"
)

type ServerInterceptorBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string

	// 这个其实是为了我们 homework2_fastest 负载均衡设计的，因为正常情况下，我们不太可能
	// 一个进程启动多个端口
	Port string
}

func (b ServerInterceptorBuilder) BuildUnary() grpc.UnaryServerInterceptor {
	address := observability.GetOutboundIP()
	if b.Port != "" {
		address = address + ":" + b.Port
	}
	// 这个部分可以简化，比如说用默认值，只需要用户传入一个应用名字
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_response",
		Help:      b.Help,
		ConstLabels: map[string]string{
			"address": address,
			"kind":    "server",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"method"})

	errCntVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_error_cnt",
		Help:      b.Help,
		ConstLabels: map[string]string{
			"address": address,
			"kind":    "server",
		},
	}, []string{"method"})

	reqCntVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_active_req_cnt",
		Help:      b.Help,
		ConstLabels: map[string]string{
			"address": address,
			"kind":    "server",
		},
	}, []string{"method"})
	prometheus.MustRegister(summaryVec, errCntVec, reqCntVec)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		reqCnt := reqCntVec.WithLabelValues(info.FullMethod)
		reqCnt.Add(1)
		startTime := time.Now()
		// 类似于 opentelemetry，这里也可以记录一下业务ID之类的信息
		defer func() {
			if err != nil {
				errCntVec.WithLabelValues(info.FullMethod).Add(1)
			}
			duration := time.Now().Sub(startTime)
			reqCnt.Sub(1)
			summaryVec.WithLabelValues(info.FullMethod).Observe(float64(duration.Milliseconds()))
		}()
		resp, err = handler(ctx, req)
		return
	}
}
