package v1

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_setFuncField(t *testing.T) {

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
				req: &Request{
					ServiceName: "user-service",
					Method:      "GetById",
					Data:        []byte(`{"msg":"这是GetById"}`),
				},
				resp: &Response{
					Data: []byte(`{"msg":"这是GetById的响应"}`),
				},
			},
			wantResp: &AnyResponse{
				Msg: "这是GetById的响应",
			},
		},
	}

	// 通过传入 mockProxy，然后我们执行一下方法上面的调用，
	// 确保已经篡改成功了
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setFuncField(tc.s.s, tc.proxy)
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
	req  *Request
	resp *Response
	err  error
}

func (m *mockProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
	assert.Equal(m.t, m.req, req)
	return m.resp, m.err
}
