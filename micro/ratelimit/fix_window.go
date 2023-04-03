package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type FixWindowLimiter struct {
	// 在 interval 内最多允许 rate 个请求
	rate     int32
	interval int64

	count      int32
	timestamp  int64
	onRejected rejectStrategy
}

func NewFixWindowLimiter(rate int32, interval time.Duration) *FixWindowLimiter {
	return &FixWindowLimiter{
		interval:   interval.Nanoseconds(),
		rate:       rate,
		onRejected: defaultRejection,
	}
}

func (l *FixWindowLimiter) LimitUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		current := time.Now().UnixNano()
		timestamp := atomic.LoadInt64(&l.timestamp)
		// interval 是只读的，所以我们不需要使用原子操作
		if timestamp+l.interval < current {
			cnt := atomic.LoadInt32(&l.count)
			// 重置，注意，这里任何一步 CAS 操作失败，都意味着有别的 goroutine 重置了
			// 所以我们失败了就直接忽略
			if atomic.CompareAndSwapInt64(&l.timestamp, timestamp, current) {
				atomic.CompareAndSwapInt32(&l.count, cnt, 0)
			}
		}
		cnt := atomic.AddInt32(&l.count, 1)
		defer atomic.AddInt32(&l.count, -1)
		// this operation is thread-safe, but count + 1 may be overflow
		if cnt <= l.rate {
			return l.onRejected(ctx, info, req, handler)
		}
		return handler(ctx, req)
	}
}
