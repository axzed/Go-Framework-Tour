package ratelimit

import (
	"context"
	_ "embed"
	"github.com/go-redis/redis/v9"
	"google.golang.org/grpc"
	"time"
)

//go:embed lua/sliding_window.lua
var luaSlidingWindow string

type RedisSlidingWindowLimiter struct {
	key string
	// 窗口内的流量阈值
	rate int
	// 窗口大小，毫秒
	interval   int64
	onRejected rejectStrategy
	client     redis.Cmdable
}

func NewRedisSlidingWindow(client redis.Cmdable, key string, rate int, interval time.Duration) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		client:     client,
		rate:       rate,
		key:        key,
		interval:   interval.Milliseconds(),
		onRejected: defaultRejection,
	}
}

func (l *RedisSlidingWindowLimiter) LimitUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 这里直接占了坑
		limit, err := l.limit(ctx)
		if err != nil {
			// 正常来说，遇到 error 表示你也不知道要不要限流
			// 那么你可以选择限流，也可以选择不限流
			return nil, err
		}
		if limit {
			return l.onRejected(ctx, info, req, handler)
		}
		return handler(ctx, req)
	}
}

func (l *RedisSlidingWindowLimiter) limit(ctx context.Context) (bool, error) {
	now := time.Now()
	return l.client.Eval(ctx, luaSlidingWindow, []string{l.key}, l.rate, l.interval, now.UnixMilli()).Bool()
}
