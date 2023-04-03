package demo

import (
	"context"
	"errors"
	"time"
)

type BloomCache struct {
	BloomFilter
	Cache
	LoadFunc func(ctx context.Context, key string) (any, error)
}

func (b *BloomCache) Get(ctx context.Context, key string) (any, error) {
	val, err := b.Cache.Get(ctx, key)
	if err!= nil && err != errKeyNotFound {
		return nil, err
	}
	if err == errKeyNotFound {
		exist := b.BloomFilter.Exist(key)
		if exist {
			val, err = b.LoadFunc(ctx, key)
			b.Cache.Set(ctx, key, val, time.Minute)
		}
	}
	return val, err
}

type BloomCacheV1 struct {
	*ReadThroughCache
}

func NewBloomCacheV1(c Cache, b BloomFilter, lf func(ctx context.Context, key string) (any, error) ) *BloomCacheV1 {
	return &BloomCacheV1{
		ReadThroughCache: &ReadThroughCache{
			Cache: c,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				exist := b.Exist(key)
				if exist {
					return lf(ctx, key)
				}
				return nil, errors.New("数据不存在")
			},
		},
	}
}



type BloomFilter interface {
	Exist(key string) bool
}


// 加了限流的实现
type LimitCache struct {

}