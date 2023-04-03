//go:build answer

package sync

import "sync"

type SafeMap[K comparable, V any] struct {
	m     map[K]V
	mutex sync.RWMutex
}

// LoadOrStore loaded 代表是返回老的对象，还是返回了新的对象
func (s *SafeMap[K, V]) LoadOrStore(key K,
	newVale V) (val V, loaded bool) {
	s.mutex.RLock()
	val, ok := s.m[key]
	s.mutex.RUnlock()
	if ok {
		return val, true
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok = s.m[key]
	if ok {
		return val, true
	}
	s.m[key] = newVale
	return newVale, false
}

type valProvider[V any] func() V

func (s *SafeMap[K, V]) LoadOrStoreHeavy(key K, p valProvider[V]) (val interface{}, loaded bool) {
	s.mutex.RLock()
	val, ok := s.m[key]
	s.mutex.RUnlock()
	if ok {
		return val, true
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok = s.m[key]
	if ok {
		return val, true
	}
	newVale := p()
	s.m[key] = newVale
	return newVale, false
}

func (s *SafeMap[K, V]) CheckAndDoSomething() {
	s.mutex.Lock()
	// check and do something
	s.mutex.Unlock()
}

func (s *SafeMap[K, V]) CheckAndDoSomething1() {
	s.mutex.RLock()
	// check 第一次检查
	s.mutex.RUnlock()

	s.mutex.Lock()
	// check and doSomething
	defer s.mutex.Unlock()
}

type Counter struct {
	i int
}

func (c *Counter) Incr() {
	c.i++
}

func (c *Counter) Get() int {
	return c.i
}