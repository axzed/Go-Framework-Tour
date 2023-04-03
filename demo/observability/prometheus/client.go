package prometheus

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo/observability"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"time"
)

type ClientInterceptorBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Port string
}

func (b ClientInterceptorBuilder) BuildUnary() grpc.UnaryClientInterceptor {
	// 你也可以考虑使用服务注册的地址
	ip := observability.GetOutboundIP()
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name: b.Name,
		Help: b.Help,
		ConstLabels: map[string]string{
			"component": "client",
			"address": ip + b.Port,
		},
	}, []string{"method"})
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			summaryVec.WithLabelValues(method).Observe(float64(duration.Milliseconds()))
		}()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
