package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Acquire(ctx context.Context, req interface{}) (interface{}, error)
	Release(resp interface{}, err error)
}

type rejectStrategy func(ctx context.Context,
	info *grpc.UnaryServerInfo, req interface{}, handler grpc.UnaryHandler) (interface{}, error)

var defaultRejection rejectStrategy = func(ctx context.Context, info *grpc.UnaryServerInfo,
	req interface{}, handler grpc.UnaryHandler) (interface{}, error) {
	return nil, status.Errorf(codes.ResourceExhausted, "触发限流 %s", info.FullMethod)
}

type Guardian interface {
	Allow(ctx context.Context, req interface{}) (cb func(), err error)
	AllowV1(ctx context.Context, req interface{}) (cb func(), resp interface{},  err error)
	OnRejection(ctx context.Context, req interface{}) (interface{}, error)
}

// func Limit() {
// 	var g Guardian
// 	cb, err := g.Allow(xx)
// 	if  err != nil {
// 		return g.OnRejection(ctx, req)
// 	}
// 	cb, resp, err := !g.Allow(xx)
// 	if err != nil {
// 		return resp, err
// 	}
//
// 	// 执行增长的业务逻辑
// }