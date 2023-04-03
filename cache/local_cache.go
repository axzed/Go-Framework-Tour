// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errKeyNotFound      = errors.New("cache: key 不存在")
	errKeyExpired       = errors.New("cache: key 已经过期")
	errOverCapacity     = errors.New("cache: 超过缓存最大容量")
	errFailedToSetCache = errors.New("cache: 设置键值对失败")
)

type BuildinMapCache struct {
	lock          sync.RWMutex
	data          map[string]*item
	close         chan struct{}
	closed        bool
	onEvicted     func(key string, val any)
	cycleInterval time.Duration
}

type BuildinMapCacheOption func(b *BuildinMapCache)

func BuildinMapWithCycleInterval(interval time.Duration) BuildinMapCacheOption {
	return func(b *BuildinMapCache) {
		b.cycleInterval = interval
	}
}

func NewBuildinMapCache(opts ...BuildinMapCacheOption) *BuildinMapCache {
	res := &BuildinMapCache{
		data:          make(map[string]*item),
		cycleInterval: time.Second * 10,
	}
	for _, opt := range opts {
		opt(res)
	}
	res.checkCycle()
	return res
}

func (b *BuildinMapCache) Get(ctx context.Context, key string) (any, error) {
	b.lock.RLock()
	// if b.closed {
	// 	return nil, errors.New("缓存已经被关闭")
	// }
	val, ok := b.data[key]
	b.lock.RUnlock()
	if !ok {
		return nil, errKeyNotFound
	}
	// 别的 goroutine 设置值了
	now := time.Now()
	if val.deadlineBefore(now) {
		b.lock.Lock()
		defer b.lock.Unlock()
		//
		// if b.closed {
		// 	return nil, errors.New("缓存已经被关闭")
		// }

		val, ok = b.data[key]
		if !ok {
			return nil, errKeyNotFound
		}
		if val.deadlineBefore(now) {
			b.delete(key)
			// 要注意，这里可以返回 errKeyNotFound
			return nil, errKeyExpired
		}
	}
	return val.val, nil
}

func (b *BuildinMapCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.closed {
		return errors.New("缓存已经被关闭")
	}
	var dl time.Time
	if expiration > 0 {
		dl = time.Now().Add(expiration)
	}
	b.data[key] = &item{
		val:      val,
		deadline: dl,
	}
	return nil
}

// func (b *BuildinMapCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
// 	b.lock.Lock()
// 	b.data[key] = val
// 	b.lock.Unlock()
// 	if expiration > 0 {
// 		time.AfterFunc(expiration, func() {
// 			delete(b.data, key)
// 		})
// 	}
//
// 	return nil
// }

func (b *BuildinMapCache) Delete(ctx context.Context, key string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.closed {
		return errors.New("缓存已经被关闭")
	}
	b.delete(key)
	return nil
}

func (b *BuildinMapCache) checkCycle() {
	go func() {
		ticker := time.NewTicker(b.cycleInterval)
		for {
			select {
			case now := <-ticker.C:
				b.lock.Lock()
				for key, val := range b.data {
					// 设置了过期时间，并且已经过期
					if !val.deadline.IsZero() &&
						val.deadline.Before(now) {
						b.delete(key)
					}
				}
				b.lock.Unlock()
			case <-b.close:
				close(b.close)
				return
			}
		}
	}()
}

func (b *BuildinMapCache) delete(key string) {
	val, ok := b.data[key]
	if ok {
		delete(b.data, key)
		if b.onEvicted != nil {
			b.onEvicted(key, val.val)
		}
	}
}

func (b *BuildinMapCache) Close() error {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.close <- struct{}{}
	if b.onEvicted != nil {
		for key, val := range b.data {
			b.onEvicted(key, val.val)
		}
	}

	b.data = nil
	return nil
}

func (b *BuildinMapCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	itm, ok := b.data[key]
	if !ok {
		return nil, errKeyNotFound
	}
	b.delete(key)
	return itm.val, nil
}

type item struct {
	val      any
	deadline time.Time
}

func (i *item) deadlineBefore(t time.Time) bool {
	return !i.deadline.IsZero() && i.deadline.Before(t)
}

type MaxCntCache struct {
	*BuildinMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(c *BuildinMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildinMapCache: c,
		maxCnt:          maxCnt,
	}
	origin := c.onEvicted
	c.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}
	return res
}

func (c *MaxCntCache) Set(ctx context.Context,
	key string, val any, expiration time.Duration) error {
	cnt := atomic.AddInt32(&c.cnt, 1)
	if cnt > c.maxCnt {
		atomic.AddInt32(&c.cnt, -1)
		return errOverCapacity
	}
	return c.BuildinMapCache.Set(ctx, key, val, expiration)
}
