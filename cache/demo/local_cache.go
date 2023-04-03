package demo

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/cache/demo/internal/errs"
	"sync"
	"time"
)



type LocalCache struct {
	data map[string]any
	mutex sync.RWMutex
	close chan struct{}
	closeOnce sync.Once

	// onEvicted func(key string, val any) error
	onEvicted func(key string, val any)
	// onEvicted func(ctx context.Context, key string, val any) error
	// onEvicted func(ctx context.Context, key string, val any)
}

func NewLocalCache(onEvicted func(key string, val any)) *LocalCache {
	res := &LocalCache{
		close: make(chan struct{}),
	}
	// 间隔时间，过长则过期的缓存迟迟得不到删除
	// 过短，则频繁执行，效果不好（过期的 key 很少）
	ticker := time.NewTicker(time.Second)
	go func() {
		// 没有时间间隔，不断遍历
		for {
			select {
			// case now := <-ticker.C:
			case <-ticker.C:
				// 00:01:00
				cnt := 0
				res.mutex.Lock()
				for k, v := range res.data {
					if v.(*item).deadline.Before(time.Now()) {
						res.delete(k, v.(*item).val)
					}
					cnt ++
					if cnt >= 1000 {
						break
					}
				}
				res.mutex.Unlock()
			case <- res.close:
				return
			}
		}
	}()
	return res
}

func (l *LocalCache) delete(key string, val any) {
	delete(l.data, key)
	if l.onEvicted != nil {
		l.onEvicted(key, val)
	}
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {
	l.mutex.RLock()
	val, ok := l.data[key]
	l.mutex.RUnlock()
	if !ok {
		return nil, errs.NewErrKeyNotFound(key)
	}
	// 别人在这里调用 Set
	itm := val.(*item)
	if itm.deadline.Before(time.Now()) {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		val, ok = l.data[key]
		if !ok {
			return nil, errs.NewErrKeyNotFound(key)
		}
		itm = val.(*item)
		if itm.deadline.Before(time.Now()) {
			l.delete(key, itm.val)
		}
		return nil, errs.NewErrKeyNotFound(key)
	}
	return itm.val, nil
}

// Set(ctx, "key1", value1, time.Minute)
// 执行业务三十秒
// Set(ctx, "key1", value2, time.Minute)
// 再三十秒，第一个 time.AfterFunc 就要执行了
func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 这个是有破绽的
	// time.AfterFunc(expiration, func() {
	// 	if l.m.Load(key).expiration
	// 	l.Delete(context.Background(), key)
	// })
	// 如果你想支持永不过期的，expiration = 0 就说明永不过期
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.data[key] = &item{
		val: val,
		deadline: time.Now().Add(expiration),
	}
	return nil
}

func (l *LocalCache) Delete(ctx context.Context, key string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	val, ok := l.data[key]
	if ok {
		return nil
	}
	l.delete(key, val.(*item).val)
	return nil
}


// close 无缓存，调用两次 Close 呢？第二次会阻塞
// close 1 缓存，调用三次就会阻塞
func (l *LocalCache) Close() error {
	// 这种写法，第二次调用会 panic
	// l.close <- struct{}{}
	// close(l.close)


	// 这种写法最好
	l.closeOnce.Do(func() {
		l.close <- struct{}{}
		close(l.close)
	})


	// 使用 select + default 防止多次 close 阻塞调用者
	// select {
	// case l.close<- struct{}{}:
	// 关闭 channel 要小心，发送数据到已经关闭的 channel 会引起 panic
	// 	close(l.close)
	// default:
	// 	// return errors.New("cache: 已经被关闭了")
	// 	return nil
	// }
	return nil
}

// 可以考虑用 sync.Pool 来复用的，是复用，不是缓存
type item struct {
	val any
	deadline time.Time
}