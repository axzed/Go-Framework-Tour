//go:build answer

package cache

import (
	"context"
	"github.com/gotomicro/ekit/list"
	"sync"
	"time"
)

type MaxMemoryCache struct {
	Cache
	max  int64
	used int64
	// 使用 mutex 性能最差，但是肯定不会有并发问题
	// 可以考虑替换原子操作，但是内存控制就不准确了
	mutex *sync.Mutex
	// 这个地方因为 linked list 删除必须要遍历，所以性能很差
	// 按照正常的设计，这边需要的是一个接近  Java 的 LinkedHashMap 的结构
	keys *list.LinkedList[string]

	// 理论上优秀设计
	// s Strategy
}

// 淘汰策略接口
// type Strategy interface {
//
// }

func NewMaxMemoryCache(max int64, cache Cache) *MaxMemoryCache {
	res := &MaxMemoryCache{
		max:   max,
		Cache: cache,
		mutex: &sync.Mutex{},
		keys:  list.NewLinkedList[string](),
	}
	res.Cache.OnEvicted(res.evicted)
	return res
}

func (m *MaxMemoryCache) Set(ctx context.Context, key string, val []byte,
	expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 也可以用 Get，但是要记得调整 keys 和计算容量变化差值
	_, _ = m.Cache.LoadAndDelete(ctx, key)
	for m.used+int64(len(val)) > m.max {
		k, err := m.keys.Get(0)
		if err != nil {
			return err
		}
		_ = m.Cache.Delete(ctx, k)
	}
	err := m.Cache.Set(ctx, key, val, expiration)
	if err == nil {
		m.used = m.used + int64(len(val))
		_ = m.keys.Append(key)
	}

	return nil
}

func (m *MaxMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	// 加锁是为了防止遇上懒惰删除的情况，触发了删除
	m.mutex.Lock()
	defer m.mutex.Unlock()
	val, err := m.Cache.Get(ctx, key)
	if err == nil {
		// 把原本的删掉
		// 然后将 key 加到末尾
		m.deleteKey(key)
		_ = m.keys.Append(key)
	}
	return val, err
}

func (m *MaxMemoryCache) Delete(ctx context.Context, key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Cache.Delete(ctx, key)
}

func (m *MaxMemoryCache) LoadAndDelete(ctx context.Context, key string) ([]byte, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Cache.LoadAndDelete(ctx, key)
}

func (m *MaxMemoryCache) OnEvicted(fn func(key string, val []byte)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Cache.OnEvicted(func(key string, val []byte) {
		m.evicted(key, val)
		fn(key, val)
	})
}

func (m *MaxMemoryCache) evicted(key string, val []byte) {
	m.used = m.used - int64(len(val))
	m.deleteKey(key)
}

func (m *MaxMemoryCache) deleteKey(key string) {
	for i := 0; i < m.keys.Len(); i++ {
		ele, _ := m.keys.Get(i)
		if ele == key {
			_, _ = m.keys.Delete(i)
			return
		}
	}
}
