package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuildinMapCache_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct{
		name string
		key string
		wantVal any
		wantErr error
	} {
		{
			name: "exist",
			key: "key1",
			wantVal: "value1",
		},
		{
			name: "not exist",
			key: "invalid",
			wantErr: errKeyNotFound,
		},
	}

	c := NewBuildinMapCache()
	err := c.Set(context.Background(), "key1", "value1", 2 * time.Second)
	require.NoError(t, err)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := c.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
	time.Sleep(time.Second * 3)
	_, err = c.Get(context.Background(), "key1")
	assert.Equal(t, errKeyExpired, err)
}

func TestBuildinMapCache_checkCycle(t *testing.T) {
	c := NewBuildinMapCache(BuildinMapWithCycleInterval(time.Second))
	err := c.Set(context.Background(), "key1", "value1", time.Millisecond * 100)
	require.NoError(t, err)
	// 以防万一
	time.Sleep(time.Second * 3)
	_, err = c.Get(context.Background(), "key1")
	assert.Equal(t, errKeyNotFound, err)
}

func TestMaxCntCache_Set(t *testing.T) {
	delegate := NewBuildinMapCache()
	cntCache := NewMaxCntCache(delegate, 2)
	err := cntCache.Set(context.Background(), "key1", 123, time.Second * 10)
	require.NoError(t, err)
	err = cntCache.Set(context.Background(), "key2", 123, time.Second * 10)
	require.NoError(t, err)
	err = cntCache.Set(context.Background(), "key3", 123, time.Second * 10)
	assert.Equal(t, errOverCapacity, err)

	err = cntCache.Delete(context.Background(), "key1")
	require.NoError(t, err)

	// 可以放进去了
	err = cntCache.Set(context.Background(), "key3", 123, time.Second * 10)
	require.NoError(t, err)
}