package ratelimit

import (
	"context"
	"google.golang.org/grpc"
)

type MethodLimiter struct {
	Limiter
	FullMethod string
}

func (m *MethodLimiter) LimitUnary() grpc.UnaryServerInterceptor {
	interceptor := m.Limiter.LimitUnary()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if info.FullMethod == m.FullMethod {
			return interceptor(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}
}
