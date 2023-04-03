package demo

import (
	"context"
	"time"
)

type WriteThroughCache struct {
	Cache
	StoreFunc func(ctx context.Context, key string, val any) error
}

func (w *WriteThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// 直接在这里开 goroutine，就是全异步

	err := w.StoreFunc(ctx, key, val)
	if err != nil {
		return err
	}

	// 这里开就是半异步

	return w.Cache.Set(ctx, key, val, expiration)

	// err :=
	// if err != nil {
	// 	return err
	// }


	// 万一我这里失败了呢？我要不要把缓存删掉？
	// err =  w.StoreFunc(ctx, key, val)
	// if err != nil {
	// 	w.Delete(ctx, key)
	// }
}
