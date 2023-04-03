package integration

import (
	"context"
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRpc(t *testing.T) {
	address := ":8081"
	usServer := &UserServiceServer{}
	server := rpc.NewServer()
	go func() {
		server.RegisterService(usServer)
		server.Start(address)
	}()
	defer server.Close()

	// 确保服务端已经启动完成
	time.Sleep(time.Second * 3)

	usClient := &UserServiceClient{}
	client, err := rpc.NewClient(address)
	require.NoError(t, err)
	err = client.InitService(usClient)
	require.NoError(t, err)

	testCases := []struct {
		name string

		req  *GetByIdRequest
		resp *GetByIdResponse
		err  error
	}{
		{
			name: "both",
			req: &GetByIdRequest{
				Id: 12,
			},
			resp: &GetByIdResponse{
				Msg: "hello, world",
			},
			err: errors.New("this is an error"),
		},
		{
			name: "response",
			req: &GetByIdRequest{
				Id: 12,
			},
			resp: &GetByIdResponse{
				Msg: "hello, world",
			},
		},
		{
			name: "error",
			req: &GetByIdRequest{
				Id: 12,
			},
			err: errors.New("this is an error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 让服务端返回这个数据
			usServer.resp = tc.resp
			usServer.err = tc.err

			// 发起调用
			resp, err := usClient.GetById(context.Background(), tc.req)

			// 确保服务端返回了对应的数据
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.resp, resp)
		})
	}
}

func TestRpcOneway(t *testing.T) {
	address := ":8081"
	usServer := &UserServiceServer{}
	server := rpc.NewServer()
	go func() {
		server.RegisterService(usServer)
		server.Start(address)
	}()
	defer server.Close()

	// 确保服务端已经启动完成
	time.Sleep(time.Second * 3)

	usClient := &UserServiceClient{}
	client, err := rpc.NewClient(address)
	require.NoError(t, err)
	err = client.InitService(usClient)
	require.NoError(t, err)

	_, err = usClient.GetById(rpc.CtxWithOneway(context.Background()), &GetByIdRequest{})
	assert.EqualError(t, err, "client: 这是 oneway 调用", err)
}

type UserServiceClient struct {
	GetById func(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error)
}

func (u *UserServiceClient) ServiceName() string {
	return "user-service"
}

type GetByIdRequest struct {
	Id int
}

type GetByIdResponse struct {
	Msg string
}

type UserServiceServer struct {
	// 接收到的请求
	req *GetByIdRequest

	// 希望返回的数据
	resp *GetByIdResponse
	err  error
}

func (u *UserServiceServer) ServiceName() string {
	return "user-service"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdRequest) (*GetByIdResponse, error) {
	return u.resp, u.err
}
