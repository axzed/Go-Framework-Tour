package cache

import (
	"context"
	"log"
	"time"
)

// RateLimitReadThroughCache 是一个装饰器
// 在原本 Cache 的功能上添加了 read through 功能
type RateLimitReadThroughCache struct {
	Cache
	LoadFunc func(ctx context.Context, key string) (any, error)
}

func (r *RateLimitReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == KeyNotFound && ctx.Value("limited") == nil {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			if er := r.Set(ctx, key, val, time.Minute); er != nil {
				log.Fatalf("刷新缓存失败, err: %v", err)
			}
		}
	}
	return val, err
}
