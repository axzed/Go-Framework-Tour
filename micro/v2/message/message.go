package message

import (
	"bytes"
	"encoding/binary"
)

// 头部不定长字段的分隔符
const (
	splitter     = '\n'
	pairSplitter = '\r'

	// 头部长度和 body 长度都是四个字节

	HeadLengthBytes = 4
	BodyLengthBytes = 4
)

type Request struct {
	// 头部
	// 消息长度
	HeadLength uint32
	// 协议版本
	BodyLength uint32
	// 消息 ID
	MessageId uint32
	// 版本，一个字节
	Version uint8
	// 压缩算法
	Compresser uint8
	// 序列化协议
	Serializer uint8

	// 服务名和方法名
	ServiceName string
	Method      string

	// 扩展字段，用于传递自定义元数据
	Meta map[string]string

	// 协议体
	Data []byte
}

func (req *Request) SetHeadLength() {
	// 固定部分
	res := 15
	res += len(req.ServiceName)
	// 分隔符
	res++
	res += len(req.Method)
	// 分隔符
	res++
	for key, value := range req.Meta {
		res += len(key)
		// 键值对分隔符
		res++
		res += len(value)
		// 分隔符
		res++
	}
	req.HeadLength = uint32(res)
}

func EncodeReq(req *Request) []byte {
	bs := make([]byte, req.HeadLength+req.BodyLength)

	cur := bs
	// 1. 写入 HeadLength，四个字节
	binary.BigEndian.PutUint32(cur[:4], req.HeadLength)
	cur = cur[4:]
	// 2. 写入 BodyLength 四个字节
	binary.BigEndian.PutUint32(cur[:4], req.BodyLength)
	cur = cur[4:]

	// 3. 写入 message id, 四个字节
	binary.BigEndian.PutUint32(cur[:4], req.MessageId)
	cur = cur[4:]

	// 4. 写入 version，因为本身就是一个字节，所以不用进行编码了
	cur[0] = req.Version
	cur = cur[1:]

	// 5. 写入压缩算法
	cur[0] = req.Compresser
	cur = cur[1:]

	// 6. 写入序列化协议
	cur[0] = req.Serializer
	cur = cur[1:]

	// 7. 写入服务名
	copy(cur, req.ServiceName)
	cur[len(req.ServiceName)] = splitter
	cur = cur[len(req.ServiceName)+1:]

	// 写入方法名
	copy(cur, req.Method)
	cur[len(req.Method)] = splitter
	cur = cur[len(req.Method)+1:]

	for key, value := range req.Meta {
		copy(cur, key)
		cur[len(key)] = pairSplitter

		cur = cur[len(key)+1:]
		copy(cur, value)
		cur[len(value)] = splitter
		cur = cur[len(value)+1:]
	}
	// 剩下的数据
	copy(cur, req.Data)
	return bs
}

// DecodeReq 解析 Request
func DecodeReq(bs []byte) *Request {
	req := &Request{}
	// 按照 EncodeReq 写下来
	// 1. 读取 HeadLength
	req.HeadLength = binary.BigEndian.Uint32(bs[:4])
	// 2. 读取 BodyLength
	req.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	// 3. 读取 message id
	req.MessageId = binary.BigEndian.Uint32(bs[8:12])
	// 4. 读取 Version
	req.Version = bs[12]
	// 5. 读取压缩算法
	req.Compresser = bs[13]
	// 6. 读取序列化协议
	req.Serializer = bs[14]
	// 7. 拆解服务名和方法名
	meta := bs[15:req.HeadLength]
	index := bytes.IndexByte(meta, splitter)
	req.ServiceName = string(meta[:index])
	meta = meta[index+1:]
	index = bytes.IndexByte(meta, splitter)
	req.Method = string(meta[:index])
	meta = meta[index+1:]

	// 继续拆解 meta 剩下的 key value
	for len(meta) > 0 {
		// 这个地方不好预估容量，但是大部分都很少，我们把现在能够想到的元数据都算法
		// 也就不超过四个
		metaMap := make(map[string]string, 4)
		index = bytes.IndexByte(meta, splitter)
		for index != -1 {
			// 一个键值对
			pair := meta[:index]
			// 我们使用 \r 来切分键值对
			pairIndex := bytes.IndexByte(meta, '\r')
			metaMap[string(pair[:pairIndex])] = string(pair[pairIndex+1:])
			meta = meta[index+1:]
			index = bytes.IndexByte(meta, splitter)
		}
		req.Meta = metaMap
	}

	// 剩下的就是数据了
	req.Data = bs[req.HeadLength:]
	return req
}

type Response struct {
	// 消息长度
	HeadLength uint32
	// 协议版本
	BodyLength uint32
	// 消息 ID
	MessageId uint32
	// 版本，一个字节
	Version uint8
	// 压缩算法
	Compresser uint8
	// 序列化协议
	Serializer uint8

	Error []byte

	Data []byte
}

func (resp *Response) SetHeadLength() {
	resp.HeadLength = uint32(15 + len(resp.Error))
}

// 这里处理 Resp 我直接复制粘贴，是因为我觉得复制粘贴会使可读性更高

func EncodeResp(resp *Response) []byte {
	bs := make([]byte, resp.HeadLength+resp.BodyLength)

	cur := bs
	// 1. 写入 HeadLength，四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.HeadLength)
	cur = cur[4:]
	// 2. 写入 BodyLength 四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.BodyLength)
	cur = cur[4:]

	// 3. 写入 message id, 四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.MessageId)
	cur = cur[4:]

	// 4. 写入 version，因为本身就是一个字节，所以不用进行编码了
	cur[0] = resp.Version
	cur = cur[1:]

	// 5. 写入压缩算法
	cur[0] = resp.Compresser
	cur = cur[1:]

	// 6. 写入序列化协议
	cur[0] = resp.Serializer
	cur = cur[1:]
	// 7. 写入 error
	copy(cur, resp.Error)
	cur = cur[len(resp.Error):]

	// 剩下的数据
	copy(cur, resp.Data)
	return bs
}

// DecodeResp 解析 Response
func DecodeResp(bs []byte) *Response {
	resp := &Response{}
	// 按照 EncodeReq 写下来
	// 1. 读取 HeadLength
	resp.HeadLength = binary.BigEndian.Uint32(bs[:4])
	// 2. 读取 BodyLength
	resp.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	// 3. 读取 message id
	resp.MessageId = binary.BigEndian.Uint32(bs[8:12])
	// 4. 读取 Version
	resp.Version = bs[12]
	// 5. 读取压缩算法
	resp.Compresser = bs[13]
	// 6. 读取序列化协议
	resp.Serializer = bs[14]
	// 7. error 信息
	resp.Error = bs[15:resp.HeadLength]

	// 剩下的就是数据了
	resp.Data = bs[resp.HeadLength:]
	return resp
}
