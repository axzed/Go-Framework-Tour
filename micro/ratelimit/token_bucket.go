package ratelimit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

// TokenBucketLimiter 基于令牌桶的限流
// 大多数时候我们不需要自己手写算法，直接使用
// golang.org/x/time/rate
// 这里我们还是会手写一个
type TokenBucketLimiter struct {
	tokens chan struct{}
	close  chan struct{}
}

// NewTokenBucketLimiter buffer 最多能缓存住多少 token
// interval 多久产生一个令牌
func NewTokenBucketLimiter(buffer int, interval time.Duration) *TokenBucketLimiter {
	res := &TokenBucketLimiter{
		tokens: make(chan struct{}, buffer),
		close:  make(chan struct{}),
	}
	go func() {
		producer := time.NewTicker(interval)
		defer producer.Stop()
		for {
			select {
			case <-res.close:
				// 关闭
				return
			case <-producer.C:
				select {
				//case <- res.close:
				// 关闭。在这里其实可以没有这个分支
				//return
				case res.tokens <- struct{}{}:
				default:
					// 加 default 分支防止一直没有人取令牌，我们这里不能正常退出
				}
			}
		}
	}()
	return res
}

func (l *TokenBucketLimiter) LimitUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-l.close:
			// 已经关闭了
			// 这里你可以决策，如果认为限流器被关了，就代表不用限流，那么就直接发起调用。
			// 这种情况下，还要考虑提供 Start 方法重启限流器
			// 我这里采用另外一种语义，就是我认为限流器被关了，其实代表的是整个应用关了，所以我这里退出
			return nil, errors.New("micro: 系统未被保护")
		case <-l.tokens:
			return handler(ctx, req)
		}
	}
}

func (l *TokenBucketLimiter) Close() error {
	// 直接关闭就可以
	// 多次关闭的情况我们就不处理了，用户需要自己来保证
	close(l.close)
	return nil
}
