//go:build !answer

package cache

import (
	"context"
	"time"
)

type MaxMemoryCache struct {
	Cache
	max  int64
	used int64
}

func NewMaxMemoryCache(max int64, cache Cache) *MaxMemoryCache {
	panic("implement me")
}

func (m *MaxMemoryCache) Set(ctx context.Context, key string, val []byte,
	expiration time.Duration) error {
	panic("implement me")
}
