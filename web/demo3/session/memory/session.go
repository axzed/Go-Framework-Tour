package memory

import (
	"context"
	"errors"
	web "gitee.com/geektime-geekbang/geektime-go/web/demo3/session"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Store struct {
	// 如果难以确保同一个 id 不会被多个 goroutine 来操作，就加上这个
	mutex      sync.RWMutex
	cache      *cache.Cache
	expiration time.Duration
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		cache: cache.New(expiration, time.Second),
	}
}

func (s *Store) Generate(ctx context.Context, id string) (web.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	res := &Session{
		id:     id,
		values: make(map[string]string),
	}
	s.cache.Set(id, res, s.expiration)
	return res, nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (web.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	res, ok := s.cache.Get(id)
	if ok {
		return res.(web.Session), nil
	}
	return nil, errors.New("web: Session 未找到")
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	res, ok := s.cache.Get(id)
	if ok {
		return errors.New("web: Session 未找到")
	}
	s.cache.Set(id, res, s.expiration)
	return nil
}

type Session struct {
	id     string
	mutex  sync.RWMutex
	values map[string]string
}

func (s *Session) Get(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	val, ok := s.values[key]
	if ok {
		return val, nil
	}
	return "", errors.New("web: key 未找到")
}

func (s *Session) Set(ctx context.Context, key string, val string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values[key] = val
	return nil
}

func (s *Session) ID() string {
	return s.id
}
