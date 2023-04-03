package demo

import (
	"context"
	"log"
	"time"
)

type WriteBackCache struct {
	*LocalCache
}

func NewWriteBackCache(store func(ctx context.Context, key string, val any) error) *WriteBackCache {
	return &WriteBackCache{
		LocalCache: NewLocalCache(func(key string, val any) {
			// 这个地方，context 不好设置
			// error 不好处理
			err := store(context.Background(), key, val)
			if err != nil {
				log.Fatalln(err)
			}
		}),
	}
}

func (w *WriteBackCache) Close() error {
	// 遍历所有的 key，将值刷新到数据库
	return nil
}


// 预加载
type PreloadCache struct {
	Cache
	sentinelCache *LocalCache
}

func NewPreloadCache(c Cache, loadFunc func(ctx context.Context, key string) (any, error)) *PreloadCache {

	// sentinel Cache 上的 key value 过期
	// 就把主 cache 上的数据刷新
	return &PreloadCache{
		Cache: c,
		sentinelCache: NewLocalCache(func(key string, val any) {
			val, err := loadFunc(context.Background(), key)
			if err == nil {
				err = c.Set(context.Background(), key, val, time.Minute)
				if err != nil {
					log.Fatalln(err)
				}
			}
			// 加入 for 循环来重试
			// for {
			//
			// }
		}),
	}
}

func (c *PreloadCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// sentinelExpiration 的设置原则是：
	// 确保 expiration - sentinelExpiration 这段时间内，来得及加载数据刷新缓存
	// 要注意 OnEvicted 的时机，尤其是懒删除，但是轮询删除效果又不是很好的时候
	sentinelExpiration := expiration - time.Second * 3
	err := c.sentinelCache.Set(ctx, key, "", sentinelExpiration)
	if err != nil {
		log.Fatalln(err)
	}
	return c.Cache.Set(ctx, key, val, expiration)
}