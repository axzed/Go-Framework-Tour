package v1

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func TestServer_handleConnection(t *testing.T) {
	// 用的是 json 来作为数据传输格式
	testCases := []struct {
		name     string
		conn     *mockConn
		service  Service
		wantResp []byte
	}{
		{
			name:    "user service",
			service: &UserService{},
			conn: &mockConn{
				readData: newRequestBytes(t, "user-service", "GetById", &AnyRequest{}),
			},
			wantResp: []byte(`{"msg":"这是GetById的响应"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer()
			server.RegisterService(tc.service)
			server.handleConnection(tc.conn)
			// 比较写入的数据，去掉长度字段
			data := tc.conn.writeData[8:]
			var resp Response
			err := json.Unmarshal(data, &resp)
			require.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp.Data)
		})
	}
}

type UserService struct {
}

func (u *UserService) ServiceName() string {
	return "user-service"
}

func (u *UserService) GetById(ctx context.Context, request *AnyRequest) (*AnyResponse, error) {
	return &AnyResponse{
		Msg: "这是GetById的响应",
	}, nil
}

func newRequestBytes(t *testing.T, service string, method string, input any) []byte {
	data, err := json.Marshal(input)
	require.NoError(t, err)
	req := &Request{
		ServiceName: service,
		Method:      method,
		Data:        data,
	}
	data, err = json.Marshal(req)
	require.NoError(t, err)
	return EncodeMsg(data)
}

type mockConn struct {
	net.Conn
	readData  []byte
	readIndex int
	readErr   error

	writeData []byte
	writeErr  error
}

func (m *mockConn) Read(bs []byte) (int, error) {
	if m.readIndex+len(bs) > len(m.readData) {
		return 0, io.EOF
	}
	for i := 0; i < len(bs); i++ {
		bs[i] = m.readData[m.readIndex+i]
	}
	m.readIndex = m.readIndex + len(bs)
	return len(bs), m.readErr
}

func (m *mockConn) Write(bs []byte) (int, error) {
	m.writeData = append(m.writeData, bs...)
	return len(bs), m.writeErr
}
