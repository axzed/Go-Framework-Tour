package opentelemetry

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientInterceptorBuilder struct {
	Tracer trace.Tracer
}

func (b *ClientInterceptorBuilder) BuildUnary() grpc.UnaryClientInterceptor {
	address := observability.GetOutboundIP()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		ctx, span := b.Tracer.Start(ctx, method, trace.WithSpanKind(trace.SpanKindClient))
		span.SetAttributes(attribute.String("address", address))
		ctx = b.inject(ctx)
		defer func() {
			if err != nil {
				span.SetStatus(codes.Error, "client failed")
				span.RecordError(err)
			}
			span.End()
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}

func (b *ClientInterceptorBuilder) inject(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	// 这个也可以做成可配置的
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(md))
	return metadata.NewOutgoingContext(ctx, md)
}
