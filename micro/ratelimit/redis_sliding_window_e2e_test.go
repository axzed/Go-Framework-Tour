//go:build e2e

package ratelimit

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisSlidingWindow_limit(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	testCases := []struct {
		name     string
		key      string
		rate     int
		interval time.Duration

		before func(t *testing.T)
		after  func(t *testing.T)

		wantLimit bool
		wantErr   error
	}{
		{
			name:   "init",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				res, err := rdb.ZCount(context.Background(), "redis-sliding", "-inf", "+inf").Result()
				require.NoError(t, err)
				assert.Equal(t, int64(1), res)
				rdb.Del(context.Background(), "redis-sliding")
			},
			key:      "redis-sliding",
			interval: time.Minute,
			rate:     10,
		},
		{
			name: "limit",
			before: func(t *testing.T) {
				now := time.Now().UnixMilli()
				cnt, err := rdb.ZAdd(context.Background(), "redis-sliding",
					redis.Z{Score: float64(now), Member: now + 1},
					redis.Z{Score: float64(now), Member: now + 2},
					redis.Z{Score: float64(now), Member: now + 3}).Result()
				require.NoError(t, err)
				assert.Equal(t, int64(3), cnt)
			},
			after: func(t *testing.T) {
				res, err := rdb.ZCount(context.Background(), "redis-sliding", "-inf", "+inf").Result()
				require.NoError(t, err)
				assert.Equal(t, int64(3), res)
				rdb.Del(context.Background(), "redis-sliding")
			},
			key:       "redis-sliding",
			interval:  time.Minute,
			rate:      3,
			wantLimit: true,
		},
		{
			name: "window shift",
			before: func(t *testing.T) {
				now := time.Now().UnixMilli()
				cnt, err := rdb.ZAdd(context.Background(), "redis-sliding",
					redis.Z{Score: float64(now), Member: now + 1},
					redis.Z{Score: float64(now), Member: now + 2},
					redis.Z{Score: float64(now), Member: now + 3}).Result()
				require.NoError(t, err)
				assert.Equal(t, int64(3), cnt)
				// 确保窗口滑动了
				time.Sleep(time.Second + time.Millisecond*100)
			},
			after: func(t *testing.T) {
				res, err := rdb.ZCount(context.Background(), "redis-sliding", "-inf", "+inf").Result()
				require.NoError(t, err)
				assert.Equal(t, int64(1), res)
				rdb.Del(context.Background(), "redis-sliding")
			},
			key:      "redis-sliding",
			interval: time.Second,
			rate:     3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			l := NewRedisSlidingWindow(rdb, tc.key, tc.rate, tc.interval)
			limit, err := l.limit(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLimit, limit)
		})
	}
}
