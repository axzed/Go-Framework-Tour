package cache

import (
	"context"
	"errors"
	"golang.org/x/sync/singleflight"
	"log"
	"time"
)


var (
	cache Cache
	KeyNotFound = errors.New("key not found")
)

var group = &singleflight.Group{}
func Biz(key string) {
	val, err := cache.Get(context.Background(), key)
	if err == KeyNotFound {
		val, err, _ = group.Do(key, func() (interface{}, error) {
			newVal, err := QueryFromDB(key)
			if err != nil {
				return nil, err
			}
			err = cache.Set(context.Background(), key, newVal, time.Minute)
			return newVal, err
		})
	}
	println(val)
}

func QueryFromDB(key string) (any, error) {
	panic("implement")
}

// SingleflightCacheV1 也是装饰器模式
// 进一步封装 ReadThroughCache
// 在加载数据并且刷新缓存的时候应用了 singleflight 模式
type SingleflightCacheV1 struct {
	ReadThroughCache
}

func NewSingleflightCacheV1(cache Cache,
	loadFunc func(ctx context.Context, key string)(any, error)) Cache {
	g := &singleflight.Group{}
	return &SingleflightCacheV1{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				defer func() {
					g.Forget(key)
				}()
				// 多个 goroutine 进来这里
				// 只有一个 goroutine 会真的去执行
				val, err, _ := g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				return val, err
			},
		},
	}
}

// SingleflightCacheV2 也是装饰器模式
// 进一步封装 ReadThroughCache
// 在加载数据并且刷新缓存的时候应用了 singleflight 模式
type SingleflightCacheV2 struct {
	ReadThroughCache
	group *singleflight.Group
}

func NewSingleflightCacheV2(cache Cache,
	loadFunc func(ctx context.Context, key string)(any, error)) Cache {
	return &SingleflightCacheV2{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: loadFunc,
		},
		group: &singleflight.Group{},
	}
}

func (s *SingleflightCacheV2) Get(ctx context.Context, key string) (any, error) {
	val, err := s.Cache.Get(ctx, key)
	if err == KeyNotFound {
		defer func() {
			s.group.Forget(key)
		}()
		val, err, _ = s.group.Do(key, func() (interface{}, error) {
			v, er := s.LoadFunc(ctx, key)
			if er == nil {
				if e := s.Set(ctx, key, val, time.Minute); e != nil {
					log.Fatalf("刷新缓存失败, err: %v", err)
				}
			}
			return v, er
		})
	}
	return val, err
}

type BloomFilter interface {
	HasKey(ctx context.Context, key string) (bool, error)
}

type BloomFilterCache struct {
	ReadThroughCache
	bf BloomFilter
}

func NewBloomFilterCache(cache Cache,
	bf BloomFilter,
	loadFunc func(ctx context.Context, key string)(any, error)) Cache {
	return &BloomFilterCache{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				ok, _ := bf.HasKey(ctx, key)
				if ok {
					return loadFunc(ctx, key)
				}
				return nil, errors.New("invalid key")
			},
		},
	}
}

func (s *BloomFilterCache) Get(ctx context.Context, key string) (any, error) {
	val, err := s.Cache.Get(ctx, key)
	if err == KeyNotFound {
		found, _ := s.bf.HasKey(ctx, key)
		if found {
			val, err = s.LoadFunc(ctx, key)
			if err == nil {
				if e := s.Set(ctx, key, val, time.Minute); e != nil {
					log.Fatalf("刷新缓存失败, err: %v", err)
				}
			}
		}
	}
	return val, err
}


type RandomExpirationCache struct {
	Cache
	offset func() time.Duration
}

func NewRandomExpirationCache(cache Cache, offset func()time.Duration) Cache {
	return &RandomExpirationCache{
		Cache: cache,
		offset: offset,
	}
}

func (r *RandomExpirationCache) Set(ctx context.Context,
	key string, val any, expiration time.Duration) error {

	expiration = expiration + r.offset()
	return r.Cache.Set(ctx, key, val, expiration)
}
