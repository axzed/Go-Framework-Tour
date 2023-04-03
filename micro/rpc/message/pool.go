package message

import (
	"sync"
)

var (
	reqPool = &sync.Pool{
		New: func() any {
			return &Request{}
		},
	}
	respPool = &sync.Pool{
		New: func() any {
			return &Response{}
		},
	}
)

func GetResponse() *Response {
	return respPool.Get().(*Response)
}

func PutResponse(resp *Response) {
	resp.HeadLength = 0
	// 协议版本
	resp.BodyLength = 0
	// 消息 ID
	resp.MessageId = 0
	// 版本，一个字节
	resp.Version = 0
	// 压缩算法
	resp.Compresser = 0
	// 序列化协议
	resp.Serializer = 0

	resp.Error = nil

	resp.Data = nil
}

func GetRequest() *Request {
	return reqPool.Get().(*Request)
}

func PutRequest(req *Request) {
	req.HeadLength = 0
	// 协议版本
	req.BodyLength = 0
	// 消息 ID
	req.MessageId = 0
	// 版本，一个字节
	req.Version = 0
	// 压缩算法
	req.Compresser = 0
	// 序列化协议
	req.Serializer = 0

	// 服务名和方法名
	req.ServiceName = ""
	req.Method = ""

	req.Meta = nil

	req.Data = nil
	reqPool.Put(req)
}
