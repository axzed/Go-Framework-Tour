package cache

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

// Cache 屏蔽不同的缓存中间件的差异
type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any,
		expiration time.Duration) error
	Delete(ctx context.Context, key string) error

	LoadAndDelete(ctx context.Context, key string) (any, error)

	// 作业在这里
	// OnEvicted(ctx context.Context) <- chan KV
}

// type KV struct {
// 	Key string
// 	Val any
// }


type CacheV4 interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any,
		expiration time.Duration) error
	Delete(ctx context.Context, key string) error

	LoadAndDelete(ctx context.Context, key string) (any, error)

	io.Closer
}

type ClosedCache struct {
	Cache
	mutex sync.RWMutex
	closed bool
}

func (c *ClosedCache) Get(ctx context.Context, key string) (any, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.closed {
		return nil, errors.New("closed 了缓存")
	}
	return c.Cache.Get(ctx, key)
}

func (c *ClosedCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.closed {
		return errors.New("closed 了缓存")
	}
	return c.Cache.Set(ctx, key, val, expiration)
}

func (c *ClosedCache) Delete(ctx context.Context, key string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.closed {
		return errors.New("closed 了缓存")
	}
	return c.Cache.Delete(ctx, key)
}

func (c *ClosedCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.closed {
		return nil, errors.New("closed 了缓存")
	}
	return c.Cache.LoadAndDelete(ctx, key)
}

func (c *ClosedCache) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.closed = true
	return nil
}

// func NewCache() {
// 	f := io.ReadAll()
// 	var c Cache
// 	if xxx {
// 		c =  NewBuildinMapCache()
// 	} else {
// 		c =  NewRedisCache()
// 	}

	// 继续解析配置
	// 用装饰器来装饰 c
// }