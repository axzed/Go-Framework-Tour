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

const instrumentationName = "gitee.com/geektime-geekbang/geektime-go/micro/observability/opentelemetry"

type ServerInterceptorBuilder struct {
	Tracer trace.Tracer
}

func (s ServerInterceptorBuilder) BuildUnary() grpc.UnaryServerInterceptor {
	if s.Tracer == nil {
		s.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	address := observability.GetOutboundIP()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, span := s.Tracer.Start(ctx, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
		ctx = s.extract(ctx)
		// 这里可以记录非常多的数据，一般来说可以考虑机器本身的信息，例如 ip，端口
		// 也可以考虑进一步记录和请求有关的信息，例如业务 ID
		span.SetAttributes(attribute.String("address", address))
		defer func() {
			if err != nil {
				// 在使用 err.String()
				span.SetStatus(codes.Error, "server failed")
				span.RecordError(err)
			}
			span.End()
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func (s *ServerInterceptorBuilder) extract(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(md))
}
