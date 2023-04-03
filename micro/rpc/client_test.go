package rpc

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/compress"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/message"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	serializer := json.Serializer{}
	testCases := []struct {
		name     string
		s        *mockService
		proxy    *mockProxy
		wantResp any
		wantErr  error
	}{
		{
			name: "user service",
			s: func() *mockService {
				s := &UserServiceClient{}
				return &mockService{
					s: s,
					do: func() (any, error) {
						return s.GetById(context.Background(), &AnyRequest{Msg: "这是GetById"})
					},
				}
			}(),
			proxy: &mockProxy{
				t: t,
				req: &message.Request{
					HeadLength:  36,
					BodyLength:  23,
					MessageId:   1,
					ServiceName: "user-service",
					Method:      "GetById",
					Serializer:  serializer.Code(),
					Data:        []byte(`{"msg":"这是GetById"}`),
				},
				resp: &message.Response{
					Data: []byte(`{"msg":"这是GetById的响应"}`),
				},
			},
			wantResp: &AnyResponse{
				Msg: "这是GetById的响应",
			},
		},
	}

	s := json.Serializer{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := setFuncField(s, compress.DoNothingCompressor{}, tc.s.s, tc.proxy)
			require.NoError(t, err)
			resp, err := tc.s.do()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}

}

type mockService struct {
	s  Service
	do func() (any, error)
}

type UserServiceClient struct {
	GetById func(ctx context.Context, req *AnyRequest) (*AnyResponse, error)
}

func (u *UserServiceClient) ServiceName() string {
	return "user-service"
}

type AnyRequest struct {
	Msg string `json:"msg"`
}

type AnyResponse struct {
	Msg string `json:"msg"`
}

// mockProxy
// 这里我们不用 mock 工具来生成，手写比较简单
type mockProxy struct {
	t    *testing.T
	req  *message.Request
	resp *message.Response
	err  error
}

func (m *mockProxy) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	assert.Equal(m.t, m.req, req)
	return m.resp, m.err
}
