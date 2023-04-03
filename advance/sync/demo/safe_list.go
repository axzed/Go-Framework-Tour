package demo

import "sync"

type SafeList[T any] struct {
	List[T]
	lock sync.RWMutex
}

func (s *SafeList[T]) Get(index int) (T, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.List.Get(index)
}

func (s *SafeList[T]) Append(t T) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.List.Append(t)
}
