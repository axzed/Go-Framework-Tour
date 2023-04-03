package demo

import "sync"

var PublicResource map[string]string
var PublicLock sync.RWMutex

var privateResource map[string]string
var privateLock sync.RWMutex

func NewFeature() {
	privateLock.Lock()
	defer privateLock.Unlock()
	privateResource["a"] = "b"
}

var safeResourceInstance safeResource

type safeResource struct {
	resource map[string]string
	lock     sync.RWMutex
}

func (s *safeResource) Add(key string, value string) {
	s.lock.Lock()
	defer s.lock.RUnlock()
	s.resource[key] = value
}

type SafeMap[K comparable, V any] struct {
	values map[K]V
	lock   sync.RWMutex
}

// 已经有 key，返回对应的值，然后 loaded = true
// 没有，则放进去，返回 loaded false
// goroutine 1 => ("key1", 1)
// goroutine 2 => ("key1", 2)

func (s *SafeMap[K, V]) LoadOrStoreV3(key K, newValue V) (V, bool) {
	s.lock.RLock()
	oldVal, ok := s.values[key]
	s.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	oldVal, ok = s.values[key]
	if ok {
		return oldVal, true
	}
	// goroutine1 先进来，那么这里就会变成 key1 => 1
	// goroutine2 进来，那么这里就会变成 key1 => 2
	s.values[key] = newValue
	return newValue, false
}

func (s *SafeMap[K, V]) LoadOrStoreV2(key K, newValue V) (V, bool) {
	s.lock.RLock()
	oldVal, ok := s.values[key]
	s.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	// goroutine1 先进来，那么这里就会变成 key1 => 1
	// goroutine2 进来，那么这里就会变成 key1 => 2
	s.values[key] = newValue
	return newValue, false
}

func (s *SafeMap[K, V]) LoadOrStoreV1(key K, newValue V) (V, bool) {
	s.lock.RLock()
	oldVal, ok := s.values[key]
	defer s.lock.RUnlock()
	if ok {
		return oldVal, true
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	oldVal, ok = s.values[key]
	if ok {
		return oldVal, true
	}
	// goroutine1 先进来，那么这里就会变成 key1 => 1
	// goroutine2 进来，那么这里就会变成 key1 => 2
	s.values[key] = newValue
	return newValue, false
}
