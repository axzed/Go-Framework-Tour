package demo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

// 用多少个字节来表达长度
// 用二进制去表达你的请求和响应，最多需要多少个字节
const lenBytes = 8

// 拿到完整的数据
func ReadMsg(conn net.Conn) (bs []byte, err error) {
	msgLenBytes := make([]byte, lenBytes)
	length, err := conn.Read(msgLenBytes)
	defer func() {
		if msg := recover(); msg != nil {
			err = errors.New(fmt.Sprintf("%v", msg))
		}
	}()
	if err != nil {
		return nil, err
	}

	if length != lenBytes {
		return nil, errors.New("could not read the length data")
	}

	headLength :=binary.BigEndian.Uint32(msgLenBytes[:4])
	bodyLength := binary.BigEndian.Uint32(msgLenBytes[4:8])
	// 整个请求读出来了
	bs = make([]byte, headLength + bodyLength)
	_, err = io.ReadFull(conn, bs[lenBytes:])
	// 你没读够，你根本不知道接下来怎么处理这个连接里面的数据
	// if n != headLength + bodyLength - lenBytes {
	// 	返回错误并且关掉连接
	// 	return nil, errors.New("")
	// }
	copy(bs, msgLenBytes)
	return bs, err
}

func EncodeMsg(data []byte) []byte {
	resp := make([]byte, len(data) +lenBytes)
	l := len(data)
	// 大顶端编码，把长度编码成二进制，然后放到了 resp 的前八个字节
	binary.BigEndian.PutUint64(resp, uint64(l))
	copy(resp[lenBytes:], data)
	return resp
}
