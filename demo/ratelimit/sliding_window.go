package ratelimit

import (
	"container/list"
	"context"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	interval time.Duration
	// 你需要一个 queue 来缓存住你窗口内每一个请求的时间戳
	queue      *list.List
	// 上限
	max int

	onReject rejectStrategy

	mutex sync.RWMutex
}


func (t *SlidingWindowLimiter) BuildUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		// info.FullMethod = "/user" 就限流
		// user 的和 order 的分开限流

		// user, ok :=req.(*GetUserReq)
		// if ok

		// 假如说现在是 3:17 interval 是一分钟
		current := time.Now()
		t.mutex.Lock()
		cnt := t.queue.Len()
		if cnt < t.max {
			t.queue.PushBack(current)
			t.mutex.Unlock()
			return handler(ctx, req)
		}

		// 慢路径

		// 往前回溯（所以是减号），起始时间是 2:17
		windowStartTime := current.Add(-t.interval)

		// 假如说 reqTime 是 2:12，代表它其实已经不在这个窗口里面了
		reqTime := t.queue.Front()
		for reqTime != nil && reqTime.Value.(time.Time).Before(windowStartTime) {
			// 说明这个请求不在这个窗口范围内，移除窗口
			t.queue.Remove(reqTime)
			reqTime = t.queue.Front()
		}

		cnt = t.queue.Len()
		if cnt >= t.max {
			t.mutex.Unlock()
			return t.onReject(ctx, info, req, handler)
		}
		t.queue.PushBack(current)
		t.mutex.Unlock()
		return handler(ctx, req)
	}
}
