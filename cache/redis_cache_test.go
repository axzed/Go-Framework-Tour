package cache

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/cache/mocks"
	redis "github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 不要使用这种整个测试范围的 mock，这样用例之间有依赖关系
	// mockCmd := mocks.NewMockCmdable(ctrl)
	testCases := []struct{
		name string

		// mock 数据，这样可以做到用例直接互不影响
		mock func() redis.Cmdable

		// 输入
		key string
		val any

		// 输出
		wantErr error
	} {
		{
			name:"success",
			key: "key1",
			val: 123,
			mock: func() redis.Cmdable {
				client := mocks.NewMockCmdable(ctrl)
				cmd := redis.NewStatusCmd(context.Background())
				cmd.SetVal("OK")

				client.EXPECT().
					Set(gomock.Any(), "key1", 123, time.Minute).
					Return(cmd)
				return client
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := RedisCache{
				client: tc.mock(),
			}
			err := c.Set(context.Background(), tc.key, tc.val, time.Minute)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

