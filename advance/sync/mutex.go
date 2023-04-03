
package sync

import (
	"sync"
)

// PublicResource 你永远不知道你的用户拿了它会干啥
// 他即便不用 PublicResourceLock 你也毫无办法
// 如果你用这个resource，一定要用锁
var PublicResource interface{}
var PublicResourceLock sync.Mutex

// privateResource 要好一点，祈祷你的同事会来看你的注释，知道要用锁
// 很多库都是这么写的，我也写了很多类似的代码=。=
var privateResource interface{}
var privateResourceLock sync.Mutex

// safeResource 很棒，所有的期望对资源的操作都只能通过定义在上 safeResource 上的方法来进行
type safeResource struct {
	resource interface{}
	lock     sync.Mutex
}

func (s *safeResource) DoSomethingToResource() {
	s.lock.Lock()
	defer s.lock.Unlock()
}

var l = sync.RWMutex{}

func RecursiveA() {
	l.Lock()
	defer l.Unlock()
	RecursiveB()
}

func RecursiveB() {
	RecursiveC()
}

func RecursiveC() {
	l.Lock()
	defer l.Unlock()
	RecursiveA()
}

// 锁的伪代码
// type Lock struct {
// 	state int
// }
//
// func (l *Lock) Lock() {
//
// 	i = 0
// 	for locked = CAS(UN_LOCK, LOCKED); !locked && i < 10 {
// 		i ++
// 	}
//
// 	if locked {
// 		return
// 	}
//
// 	// 将自己的线程加入阻塞队列
// 	enqueue()
// }

type List[T any] interface {
	Get(index int) T
	Set(index int, t T)
	DeleteAt(index int) T
	Append(t T)
}

type ArrayList[T any] struct {
	vals []T
}

func (a *ArrayList[T]) Get(index int) T {
	return a.vals[index]
}

func (a *ArrayList[T]) Set(index int, t T) {
	if index >= len(a.vals) || index < 0 {
		panic("index 超出范围")
	}
	a.vals[index] = t
}

func (a *ArrayList[T]) DeleteAt(index int) T {
	if index >= len(a.vals) || index < 0 {
		panic("index 超出范围")
	}
	res := a.vals[index]
	a.vals = append(a.vals[:index], a.vals[index+1:]...)
	return res
}

func (a *ArrayList[T]) Append(t T) {
	a.vals = append(a.vals, t)
}

func NewArrayList[T any](initCap int) *ArrayList[T] {
	return &ArrayList[T]{vals: make([]T, 0, initCap)}
}

type SafeListDecorator[T any] struct {
	l List[T]
	mutex sync.RWMutex
}

func (s *SafeListDecorator[T]) Get(index int) T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.l.Get(index)
}
func (s *SafeListDecorator[T]) Set(index int, t T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.l.Set(index, t)
}
func (s *SafeListDecorator[T]) DeleteAt(index int) T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.l.DeleteAt(index)
}
func (s *SafeListDecorator[T]) Append(t T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.l.Append(t)
}



