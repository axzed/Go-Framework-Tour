package ratelimit

import (
	"context"
	_ "embed"
	"github.com/go-redis/redis/v9"
	"google.golang.org/grpc"
	"time"
)

//go:embed lua/fix_window.lua
var luaFixWindow string

// RedisFixWindowLimiter 基于 Redis 实现的限流有很多搞法
// 这是其中一种比较简单的做法
type RedisFixWindowLimiter struct {
	// key 的前缀
	key string

	// 这个 key 在 interval 内只允许 rate 个请求
	rate     int
	interval time.Duration

	onRejected rejectStrategy

	// 理论上来说，我们可以抽象出来一个统一的接口
	// 然后 redis 作为其中的一个实现
	// 但是，有点没必要，因为基本上逻辑都是在 redis 的 lua 脚本里面
	// 抽象出来接口并不能提高系统的复用程度
	client redis.Cmdable
}

func NewRedisFixWindowLimiter(client redis.Cmdable, key string, rate int, interval time.Duration) *RedisFixWindowLimiter {
	return &RedisFixWindowLimiter{
		client:     client,
		key:        key,
		rate:       rate,
		interval:   interval,
		onRejected: defaultRejection,
	}
}

func (l *RedisFixWindowLimiter) LimitUnary() grpc.UnaryServerInterceptor {
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

func (l *RedisFixWindowLimiter) limit(ctx context.Context) (bool, error) {
	return l.client.Eval(ctx, luaFixWindow, []string{l.key}, l.interval.Milliseconds(), l.rate).Bool()
}
