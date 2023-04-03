package cache

import (
	"context"
	"log"
	"time"
)

// ReadThroughCache 是一个装饰器
// 在原本 Cache 的功能上添加了 read through 功能
type ReadThroughCache struct {
	Cache
	LoadFunc func(ctx context.Context, key string) (any, error)
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == KeyNotFound {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			if er := r.Set(ctx, key, val, time.Minute); er != nil {
				log.Fatalf("刷新缓存失败, err: %v", err)
			}
		}
	}
	return val, err
}

var _ Cache = &ReadThroughCacheV1[any]{}

// ReadThroughCacheV1 使用泛型会直接报错
// var c Cache= &ReadThroughCacheV1[*User]{} 编译无法通过
type ReadThroughCacheV1[T any] struct {
	Cache
	LoadFunc func(ctx context.Context, key string) (T, error)
}

func (r *ReadThroughCacheV1[T]) Get(ctx context.Context, key string) (T, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == KeyNotFound {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			if er := r.Set(ctx, key, val.(T), time.Minute); er != nil {
				log.Fatalf("刷新缓存失败, err: %v", err)
			}
		}
	}
	var t T
	if val != nil {
		t = val.(T)
	}
	return t, err
}

func (r *ReadThroughCacheV1[T]) Set(ctx context.Context, key string, val T,
expiration time.Duration) error {
	return r.Cache.Set(ctx, key, val, expiration)
}
func (r *ReadThroughCacheV1[T]) Delete(ctx context.Context, key string) error {
	return r.Cache.Delete(ctx, key)
}

func (r *ReadThroughCacheV1[T]) LoadAndDelete(ctx context.Context, key string) (T, error) {
	val, err := r.Cache.LoadAndDelete(ctx, key)
	var t T
	if val != nil {
		t = val.(T)
	}
	return t, err
}