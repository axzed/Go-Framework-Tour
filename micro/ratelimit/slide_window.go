package ratelimit

import (
	"container/list"
	"context"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type SlideWindowLimiter struct {
	rate       int
	interval   int64
	mutex      *sync.Mutex
	queue      *list.List
	onRejected rejectStrategy
}

func NewSlideWindowLimiter(rate int, interval time.Duration) *SlideWindowLimiter {
	return &SlideWindowLimiter{
		rate:       rate,
		interval:   interval.Nanoseconds(),
		mutex:      &sync.Mutex{},
		queue:      list.New(),
		onRejected: defaultRejection,
	}
}

func (l *SlideWindowLimiter) LimitUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if l.limit() {
			return l.onRejected(ctx, info, req, handler)
		}
		return handler(ctx, req)
	}
}

func (l *SlideWindowLimiter) limit() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	// 快路径
	size := l.queue.Len()
	current := time.Now().UnixNano()
	if size < l.rate {
		l.queue.PushBack(current)
		return false
	}
	// 慢路径
	boundary := current - l.interval
	timestamp := l.queue.Front()
	// 删除已经不在窗口内的元素
	for timestamp != nil && timestamp.Value.(int64) < boundary {
		l.queue.Remove(timestamp)
		timestamp = l.queue.Front()
	}
	if l.queue.Len() < l.rate {
		l.queue.PushBack(current)
		return false
	}
	return true
}
