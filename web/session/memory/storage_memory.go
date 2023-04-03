package memory

import (
	"context"
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/web/session"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Store struct {
	// 如果难以确保同一个 id 不会被多个 goroutine 来操作，就加上这个
	mutex sync.RWMutex
	// 利用一个内存缓存来帮助我们管理过期时间
	c          *cache.Cache
	expiration time.Duration
}

// NewStore 创建一个 Store 的实例
// 实际上，这里也可以考虑使用 Option 设计模式，允许用户控制过期检查的间隔
func NewStore(expiration time.Duration) *Store {
	return &Store{
		c:          cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

func (m *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	sess := &memorySession{
		id:   id,
		data: make(map[string]string),
	}
	m.c.Set(sess.ID(), sess, m.expiration)
	return sess, nil
}

func (m *Store) Refresh(ctx context.Context, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	sess, ok := m.c.Get(id)
	if !ok {
		return errors.New("session not found")
	}
	m.c.Set(sess.(*memorySession).ID(), sess, m.expiration)
	return nil
}

func (m *Store) Remove(ctx context.Context, id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.c.Delete(id)
	return nil
}

func (m *Store) Get(ctx context.Context, id string) (session.Session, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	sess, ok := m.c.Get(id)
	if !ok {
		return nil, errors.New("session not found")
	}
	return sess.(*memorySession), nil
}

type memorySession struct {
	mutex      sync.RWMutex
	id         string
	data       map[string]string
}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return "", errors.New("找不到这个 key")
	}
	return val, nil
}

func (m *memorySession) Set(ctx context.Context, key string, val string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = val
	return nil
}

func (m *memorySession) ID() string {
	return m.id
}
