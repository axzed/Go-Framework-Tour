package cache

import (
	"context"
	"errors"
	redis "github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

var (
	luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)

	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)

	ErrFailedToRefresh = errors.New("续约失败")
	ErrFailedToRelease = errors.New("释放锁失败")
)

type LockClientDemo struct {
	client *redis.Client
}

func (c *LockClientDemo) Lock(ctx context.Context, key string, expiration time.Duration) (*LockDemo, error) {
	token := uuid.New().String()
	ok, err := c.client.SetNX(ctx, key, token, expiration).Result()
	if ok {
		return &LockDemo{
			key: key,
			token: token,
			client: c.client,
			close: make(chan struct{}, 1),
		}, nil
	}
	return nil, err
}

type LockDemo struct {
	key string
	token string
	client *redis.Client
	close chan struct{}
}

func (l *LockDemo) AutoRefresh(
	timeout time.Duration,
	interval time.Duration, newExpire time.Duration) <- chan error {
	ch := make(chan error, 1)
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <- ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				status, err := luaRefresh.Run(ctx,
					l.client, []string{l.key}, l.token, newExpire.Milliseconds()).Result()
				cancel()
				if err != nil {
					ch <- err
					// 这里咩有关闭 l.close，其实问题也不大，因为最终都是垃圾回收
					close(ch)
					return
				}
				if status != int64(1) {
					ch <- ErrFailedToRefresh
					close(ch)
				}
			case <- l.close:
				close(ch)
				close(l.close)
				return
			}
		}
	}()
	return ch
}
func (l *LockDemo) Refresh(ctx context.Context, newExpire time.Duration) error {
	status, err := luaRefresh.Run(ctx, l.client,
		[]string{l.key}, l.token, newExpire.Milliseconds()).Result()
	if err != nil {
		return err
	}
	if status == int64(1) {
		return nil
	}
	return ErrFailedToRefresh
}

func (l *LockDemo) Release(ctx context.Context) error {
	defer func() {
		l.close <- struct{}{}
	}()
	status, err := luaRelease.Run(ctx, l.client,
		[]string{l.key}, l.token).Result()
	if err != nil {
		return err
	}
	if status == int64(1) {
		return nil
	}
	return ErrFailedToRelease
}