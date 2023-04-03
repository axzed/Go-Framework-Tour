package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	LimitUnary() grpc.UnaryServerInterceptor
}

type rejectStrategy func(ctx context.Context,
	info *grpc.UnaryServerInfo, req interface{}, handler grpc.UnaryHandler) (interface{}, error)

var defaultRejection rejectStrategy = func(ctx context.Context, info *grpc.UnaryServerInfo,
	req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	return nil, status.Errorf(codes.ResourceExhausted, "触发限流 %s", info.FullMethod)
}

var markLimitedRejection rejectStrategy = func(ctx context.Context, info *grpc.UnaryServerInfo,
	req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = context.WithValue(ctx, "limited", true)
	return handler(ctx, req)
}
