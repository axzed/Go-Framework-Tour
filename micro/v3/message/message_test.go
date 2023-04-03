package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeRequest(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "with meta",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  25,
				Serializer:  17,
				ServiceName: "user-service",
				Method:      "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "b",
					"shadow":   "true",
				},
				Data: []byte("hello, world"),
			},
		},
		{
			name: "no meta",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  25,
				Serializer:  17,
				ServiceName: "user-service",
				Method:      "GetById",
				Data:        []byte("hello, world"),
			},
		},
		{
			name: "empty value",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  25,
				Serializer:  17,
				ServiceName: "user-service",
				Method:      "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "",
					"shadow":   "true",
				},
				Data: []byte("hello, world"),
			},
		},
	}

	for _, tc := range testCases {
		// 这里测试我们利用 encode/decode 过程相反的特性
		t.Run(tc.name, func(t *testing.T) {
			tc.req.SetHeadLength()
			tc.req.BodyLength = uint32(len(tc.req.Data))
			bs := EncodeReq(tc.req)
			req := DecodeReq(bs)
			assert.Equal(t, tc.req, req)
		})
	}
}

func TestEncodeDecodeResponse(t *testing.T) {
	testCases := []struct {
		name string
		resp *Response
	}{
		{
			name: "no error",
			resp: &Response{
				MessageId:  123,
				Version:    12,
				Compresser: 25,
				Serializer: 17,
				Data:       []byte("hello, world"),
				Error:      []byte{},
			},
		},
		{
			name: "with error",
			resp: &Response{
				MessageId:  123,
				Version:    12,
				Compresser: 25,
				Serializer: 17,
				Data:       []byte("hello, world"),
				Error:      []byte("这是错误信息"),
			},
		},
	}

	for _, tc := range testCases {
		// 这里测试我们利用 encode/decode 过程相反的特性
		t.Run(tc.name, func(t *testing.T) {
			tc.resp.SetHeadLength()
			tc.resp.BodyLength = uint32(len(tc.resp.Data))
			bs := EncodeResp(tc.resp)
			req := DecodeResp(bs)
			assert.Equal(t, tc.resp, req)
		})
	}
}
