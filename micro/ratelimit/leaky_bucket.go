package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// LeakyBucketLimiter 漏桶算法
// 大多数时候，可以用 uber 的库 https://github.com/uber-go/ratelimit
// 这里我们依旧手写一个演示基本原理
type LeakyBucketLimiter struct {
	producer *time.Ticker
}

// NewLeakyBucketLimiter 隔多久产生一个令牌
func NewLeakyBucketLimiter(interval time.Duration) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{
		producer: time.NewTicker(interval),
	}
}

func (l *LeakyBucketLimiter) LimitUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		select {
		case <-ctx.Done():
			// 等令牌过期了
			// 这里你也可以考虑回调拒绝策略
			return nil, ctx.Err()
		case <-l.producer.C:
			return handler(ctx, req)
		}
	}
}
