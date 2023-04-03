//go:build e2e
package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_Generate(t *testing.T) {
	s := newStore()
	ctx := context.Background()
	id := "sess_test_id"
	sess, err := s.Generate(ctx, id)
	require.NoError(t, err)
	defer s.Remove(ctx, id)
	err = sess.Set(ctx, "key1", "123")
	require.NoError(t, err)
	val, err := sess.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "123", val)
}

func newStore() *Store {
	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "abc",
	})
	return NewStore(rc)
}
