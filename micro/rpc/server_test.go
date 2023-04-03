package rpc

import (
	"context"
	"encoding/json"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/message"
	json2 "gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"testing"
)

func TestServer_handleConnection(t *testing.T) {
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
	serializer := json2.Serializer{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer()
			server.RegisterSerializer(serializer)
			server.RegisterService(tc.service)
			server.handleConnection(tc.conn)
			resp := message.DecodeResp(tc.conn.writeData)
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
	req := &message.Request{
		ServiceName: service,
		Method:      method,
		Data:        data,
		// 固定用 json
		Serializer: 1,
		BodyLength: uint32(len(data)),
	}
	req.SetHeadLength()
	return message.EncodeReq(req)
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
