package ctx

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Cache interface {
	Get(key string) (string, error)
}

type OtherCache interface {
	GetValue(ctx context.Context, key string) (any, error)
}

// CacheAdapter 适配器强调的是不同接口之间进行适配
// 装饰器强调的是添加额外的功能
type CacheAdapter struct {
	Cache
}

func (c *CacheAdapter) GetValue(ctx context.Context, key string) (any, error) {
	return c.Cache.Get(key)
}

// 已有的，不是线程安全的
type memoryMap struct {
	// 如果你这样添加锁，那么就是一种侵入式的写法，
	// 那么你就需要测试这个类
	// 而且有些时候，这个是第三方的依赖，你都改不了
	// lock sync.RWMutex
	m map[string]string
}

func (m *memoryMap) Get(key string) (string, error) {
	return m.m[key], nil
}

var s = &SafeCache{
	Cache: &memoryMap{},
}

// SafeCache 我要改造为线程安全的
// 无侵入式地改造
type SafeCache struct {
	Cache
	lock sync.RWMutex
}

func (s *SafeCache) Get(key string) (string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Cache.Get(key)
}

// type valueCtx struct {
// 	context.Context
// 	vals map[any]any
// }

// func TestSourceCode(t *testing.T) {
// 	ctx := context.WithCancel(context.Background())
// }

func TestErrgroup(t *testing.T) {
	eg, ctx := errgroup.WithContext(context.Background())
	var result int64 = 0
	for i := 0; i < 10; i++ {
		delta := i
		eg.Go(func() error {
			atomic.AddInt64(&result, int64(delta))
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
	ctx.Err()
	fmt.Println(result)
}

func TestBusinessTimeout(t *testing.T) {
	ctx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	end := make(chan struct{}, 1)
	go func() {
		MyBusiness()
		end <- struct{}{}
	}()
	ch := timeoutCtx.Done()
	select {
	case <-ch:
		fmt.Println("timeout")
	case <-end:
		fmt.Println("business end")
	}
}

func MyBusiness() {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("hello, world")
}

func TestParentValueCtx(t *testing.T) {
	ctx := context.Background()
	childCtx := context.WithValue(ctx, "map", map[string]string{})
	ccChild := context.WithValue(childCtx, "key1", "value1")
	m := ccChild.Value("map").(map[string]string)
	m["key1"] = "val1"
	val := childCtx.Value("key1")
	fmt.Println(val)
	val = childCtx.Value("map")
	fmt.Println(val)
}

func TestParentCtx(t *testing.T) {
	ctx := context.Background()
	dlCtx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Minute))
	childCtx := context.WithValue(dlCtx, "key", 123)
	cancel()
	err := childCtx.Err()
	fmt.Println(err)
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	valCtx := context.WithValue(ctx, "abc", 123)
	val := valCtx.Value("abc")
	fmt.Println(val)
}

// func TestContext(t *testing.T) {
// 	ctx := context.Background()
// 	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
// 	defer cancel()
// 	time.Sleep(2 * time.Second)
// 	err := timeoutCtx.Err()
// 	fmt.Println(err)
// }

// func SomeBusiness() {
// 	ctx := context.TODO()
// 	Step1()
// }

//
// func Step1(ctx context.Context) {
// 	var db *sql.DB
// 	db.ExecContext(ctx, "UPDATE XXXX", 1)
// }

type A struct {
	ctx context.Context
}
