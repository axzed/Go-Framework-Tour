package demo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// type AsyncReadThroughCache struct {
// 	ReadThroughCache
// }


type ReadThroughCache struct {
	mutex sync.RWMutex
	Cache
	Expiration time.Duration
	// 我们把最常见的”捞DB”这种说法抽象为”加载数据”
	LoadFunc func(ctx context.Context, key string) (any, error)
	// Loader

	Async bool
}

func (c *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	// 逻辑：
	// 先捞缓存
	// 再捞 DB
	c.mutex.RLock()
	val, err := c.Cache.Get(ctx, key)
	c.mutex.RUnlock()
	if err != nil && err != errKeyNotFound {
		// 这个代表出问题了，但是不知道哪里出问题了
		return nil, err
	}

	// 缓存没有数据
	if err == errKeyNotFound {

		// 全异步，在这里直接开 goroutine
		// go func() {
		// 	c.mutex.Lock()
		// 	defer c.mutex.Unlock()
		// 	// 加锁问题
		// 	// 两个 goroutine 进来这里
		//
		//
		// 	// 捞DB
		// 	val, err = c.LoadFunc(ctx, key)
		//
		// 	// 第一个 key1=value1
		// 	// 中间有人更新了数据库
		// 	// 第二个 key1=value2
		//
		// 	if err != nil {
		// 		log.Fatalln(err)
		// 		return
		// 	}
		//
		// 	// 这里 err 可以考虑忽略掉，或者输出 warn 日志
		// 	err = c.Cache.Set(ctx, key, val, c.Expiration)
		// 	if err != nil {
		// 		log.Fatalln(err)
		// 	}
		// }()


		c.mutex.Lock()
		defer c.mutex.Unlock()
		// 加锁问题
		// 两个 goroutine 进来这里


		// 捞DB
		val, err = c.LoadFunc(ctx, key)

		// 第一个 key1=value1
		// 中间有人更新了数据库
		// 第二个 key1=value2

		if err != nil {

			// 讨论清楚
			// 1. 你的缓存框架有 Log 抽象，那么你可以打印错误
			// log.Println(err)
			// return nil, errors.New("cache: 无法加载数据")
			// 所以很显然，如果你是公司自研的缓存框架，那么你就爱怎么打就怎么打

			// 2. 你不想引入 Log，而是希望通过返回 error 来暴露信息
			// 这里你就不要丢掉原始 err 信息
			// 你可以 wrap 也可以不 wrap
			// 我个人偏好：只在我确实不希望用户知道我的底层实现的时候，我才会 wrap
			// 但是 LoadFunc 是用户指定的，不关我事，所以直接返回也没啥


			// 这里会暴露 LoadFunc 底层
			// 例如如果 LoadFunc 是数据库查询，这里就会暴露数据库的错误信息（或者 ORM 框架的）
			return nil, fmt.Errorf("cache: 无法加载数据, %w", err)

			// 转新 error 我不建议
			// return nil, errors.New("cache: 无法加载数据")
		}

		// 这里开 goroutine 就是半异步


		// 这里 err 可以考虑忽略掉，或者输出 warn 日志
		err = c.Cache.Set(ctx, key, val, c.Expiration)

		// 可能的结果: goroutine1 先，毫无问题，数据库和缓存是一致的
		// goroutine2 先，那就有问题了, DB 是 value2，但是 cache 是 value1
		if err != nil {
			log.Fatalln(err)
		}
		return val, nil
	}
	return val, nil
}

func (c *ReadThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 你加不加锁，数据都可能不一致
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Cache.Set(ctx, key, val, expiration)
}

type ReadThroughCacheV1[T any] struct {
	Cache
	Expiration time.Duration
	// 我们把最常见的”捞DB”这种说法抽象为”加载数据”
	LoadFunc func(ctx context.Context, key string) (T, error)
}

type ReadThroughCacheV2[T any] struct {
	CacheV2[T]
	Expiration time.Duration
	// 我们把最常见的”捞DB”这种说法抽象为”加载数据”
	LoadFunc func(ctx context.Context, key string) (T, error)
}

type ReadThroughCacheV3 struct {
	mutex sync.RWMutex
	Cache
	Expiration time.Duration
	// 我们把最常见的”捞DB”这种说法抽象为”加载数据”
	// LoadFunc func(ctx context.Context, key string) (any, error)
	Loader
}

type Loader interface {
	Load(ctx context.Context, key string) (any, error)
	// 如果你预期你后面会在这里加方法
}

type LoadFunc func(ctx context.Context, key string) (any, error)

func (l LoadFunc) Load(ctx context.Context, key string) (any, error) {
	return l(ctx, key)
}