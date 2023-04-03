package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

type LeakyBucketLimiter struct {
	producer *time.Ticker
	close chan struct{}
}

func NewLeakyBucketLimiter(interval time.Duration) *LeakyBucketLimiter {
	res := &LeakyBucketLimiter{
		close: make(chan struct{}),
		producer: time.NewTicker(interval),
	}
	return res
}

func (t LeakyBucketLimiter) BuildUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		select {
		case <- ctx.Done():
			return nil, ctx.Err()
		case <- t.producer.C:
			return handler(ctx, req)
		}
	}
}

func (t LeakyBucketLimiter) Close() error {
	t.producer.Stop()
	return nil
}