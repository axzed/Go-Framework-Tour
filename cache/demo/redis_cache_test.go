package demo

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/cache/demo/mocks"
	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct{
		name string
		mock func()redis.Cmdable
		key string
		val string
		expiration time.Duration

		wantErr error
	} {
		{
			name: "return OK",
			mock: func() redis.Cmdable {
				res := mocks.NewMockCmdable(ctrl)
				cmd := redis.NewStatusCmd(nil)
				cmd.SetVal("OK")
				res.EXPECT().Set(gomock.Any(), "key1", "value1", time.Minute).
					Return(cmd)
				return res
			},
			key: "key1",
			val: "value1",
			expiration: time.Minute,
		},

		{
			name: "timeout",
			mock: func() redis.Cmdable {
				res := mocks.NewMockCmdable(ctrl)
				cmd := redis.NewStatusCmd(nil)
				cmd.SetErr(context.DeadlineExceeded)
				res.EXPECT().Set(gomock.Any(), "key1", "value1", time.Minute).
					Return(cmd)
				return res
			},
			key: "key1",
			val: "value1",
			expiration: time.Minute,
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmdable := tc.mock()
			client := NewRedisCache(cmdable)
			err := client.Set(context.Background(), tc.key, tc.val, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
