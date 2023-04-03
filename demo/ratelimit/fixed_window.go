package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type FixWindowLimiter struct {
	latestWindowStartTimestamp int64
	interval int64
	max int64
	onReject rejectStrategy
	cnt int64
}

// interval => 窗口多大
// max 这个窗口内，能够执行多少个请求
func NewFixWindowLimiter(interval time.Duration, max int64) *FixWindowLimiter {
	return &FixWindowLimiter{
		interval: interval.Nanoseconds(),
		max: max,
		onReject: defaultRejection,
	}
}

func (t *FixWindowLimiter) OnReject(onReject rejectStrategy) *FixWindowLimiter {
	t.onReject = onReject
	return t
}

func (t *FixWindowLimiter) BuildUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
		current := time.Now().Nanosecond()
		window := atomic.LoadInt64(&t.latestWindowStartTimestamp)
		// 如果要是最近窗口的起始时间 + 窗口大小 < 当前时间戳
		// 说明换窗口了
		if window + t.interval < int64(current) {
			// 换窗口了
			// 重置了 latestWindowStartTimestamp
			if atomic.CompareAndSwapInt64(&t.latestWindowStartTimestamp, window, 0) {
				atomic.StoreInt64(&t.cnt, 0)
			}
		}

		// 检查这个窗口还能不能处理新请求
		// 我先取号
		cnt := atomic.AddInt64(&t.cnt, 1)
		// 超过上限了
		if cnt > t.max {
			return t.onReject(ctx, info, req, handler)
		}
		return handler(ctx, req)
	}
}
