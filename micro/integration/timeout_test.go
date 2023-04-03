package integration

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	address := ":8081"
	usServer := &UserServiceTimeoutServer{}
	server := rpc.NewServer()
	go func() {
		server.RegisterService(usServer)
		server.Start(address)
	}()
	defer server.Close()

	// 确保服务端已经启动完成
	time.Sleep(time.Second * 3)

	usClient := &UserServiceTimeoutClient{}
	client, err := rpc.NewClient(address)
	require.NoError(t, err)
	err = client.InitService(usClient)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		sleep   time.Duration
		timeout time.Duration

		req  *GetByIdRequest
		resp *GetByIdResponse
		err  error
	}{
		{
			name:    "timeout",
			timeout: time.Second,
			sleep:   time.Second * 2,
			err:     context.DeadlineExceeded,
		},
		{
			name:    "response",
			timeout: time.Second * 3,
			sleep:   time.Second,
			req: &GetByIdRequest{
				Id: 12,
			},
			resp: &GetByIdResponse{
				Msg: "hello, world",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 让服务端返回这个数据
			usServer.resp = tc.resp
			usServer.sleep = tc.sleep
			usServer.err = tc.err

			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			// 发起调用
			resp, err := usClient.GetById(ctx, tc.req)
			// 确保服务端返回了对应的数据
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.resp, resp)
			deadline, ok := ctx.Deadline()
			require.True(t, ok)
			assert.Equal(t, deadline.UnixMilli(), usServer.deadline.UnixMilli())
		})
	}
}

type UserServiceTimeoutClient struct {
	GetById func(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error)
}

func (u *UserServiceTimeoutClient) ServiceName() string {
	return "user-timeout-service"
}

type UserServiceTimeoutServer struct {
	// 接收到的请求
	req *GetByIdRequest

	sleep time.Duration
	// 希望返回的数据
	resp *GetByIdResponse
	err  error

	deadline time.Time
}

func (u *UserServiceTimeoutServer) ServiceName() string {
	return "user-timeout-service"
}

func (u *UserServiceTimeoutServer) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	u.deadline, _ = ctx.Deadline()
	time.Sleep(u.sleep)
	return u.resp, u.err
}
