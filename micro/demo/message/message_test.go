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
				MethodName:      "GetById",
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
				MethodName:      "GetById",
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
				MethodName:      "GetById",
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
			tc.req.CalHeadLength()
			tc.req.BodyLength = uint32(len(tc.req.Data))

			bs := EncodeReq(tc.req)
			req := DecodeReq(bs)
			assert.Equal(t, tc.req, req)
		})
	}
}
