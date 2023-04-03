package opentelemetry

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo/observability"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"go.opentelemetry.io/otel/attribute"
)

type ClientInterceptorBuilder struct {
	Tracer trace.Tracer
}

func (b *ClientInterceptorBuilder) BuildUnary() grpc.UnaryClientInterceptor {
	address := observability.GetOutboundIP()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		spanCtx, span := b.Tracer.Start(ctx, method, trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		span.SetAttributes(attribute.String("address", address))
		err := invoker(spanCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.RecordError(err)
		}
		return err
	}
}