package demo

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	errKeyNotFound = errors.New("cache: 找不到 key")
)
// 值的问题
// - string: 可以，问题是本地缓存，结构体转化为 string，比如用 json 表达 User
// - []byte: 最通用的表达，可以存储序列化后的数据，也可以存储加密数据，还可以存储压缩数据。用户用起来不方便
// - any: Redis 之类的实现，你要考虑序列化的问题

type Cache interface {
	// val, err  := Get(ctx)
	// str = val.(string)
	Get(ctx context.Context, key string) (any, error)
	// Get(ctx context.Context, key string) AnyValue

	Set(ctx context.Context, key string, val any, expiration time.Duration) error

	// Set(ctx context.Context, key string, val any) AnyValue
	Delete(ctx context.Context, key string) error

	// OnEvicted(func(key string, val any))

	// Incr(ctx context.Context, key string, delta int64) error
	// IncrFloat(ctx context.Context, key string, delta float64) error
}

// type MyCache struct {
// 	expiration time.Duration
// }

// func (*MyCache)Set(ctx context.Context, key string, val any) error {
	// 这里用 MyCache 上的 expiration
// }

type AnyValue struct {
	Val any
	Err error
}

func (a AnyValue) String() (string, error) {
	if a.Err != nil {
		return "", a.Err
	}
	str, ok := a.Val.(string)
	if !ok {
		return "", errors.New("无法转换的类型")
	}
	return str, nil
}


func (a AnyValue) Bytes() ([]byte, error) {
	if a.Err != nil {
		return nil, a.Err
	}
	str, ok := a.Val.([]byte)
	if !ok {
		return nil, errors.New("无法转换的类型")
	}
	return str, nil
}

func (a AnyValue) BindJson(val any) error {
	if a.Err != nil {
		return a.Err
	}
	str, ok := a.Val.([]byte)
	if !ok {
		return errors.New("无法转换的类型")
	}
	return json.Unmarshal(str, val)
}

type CacheV2[T any] interface {
	Get(ctx context.Context, key string) (T, error)

	Set(ctx context.Context, key string, val T, expiration time.Duration) error

	Delete(ctx context.Context, key string) error
}

// type CacheV3 interface {
// 	Get[T any](ctx context.Context, key string) (T, error)
//
// 	Set[T any](ctx context.Context, key string, val T, expiration time.Duration) error
//
// 	Delete(ctx context.Context, key string) error
// }
//
// type CacheV4[E any] interface {
// 	Get[T any](ctx context.Context, key string) (T, error)
//
// 	Set[T any](ctx context.Context, key string, val T, expiration time.Duration) error
//
// 	Delete(ctx context.Context, key string) error
// }