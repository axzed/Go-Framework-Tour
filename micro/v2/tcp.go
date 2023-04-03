package rpc

import (
	"encoding/binary"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/message"
	"io"
	"net"
)

func ReadMsg(conn net.Conn) (bs []byte, err error) {
	headLenBytes := make([]byte, message.HeadLengthBytes)
	_, err = conn.Read(headLenBytes)
	if err != nil {
		return nil, err
	}

	headLen := binary.BigEndian.Uint32(headLenBytes)

	bodyLenBytes := make([]byte, message.BodyLengthBytes)
	_, err = conn.Read(bodyLenBytes)
	if err != nil {
		return nil, err
	}
	bodyLen := binary.BigEndian.Uint32(bodyLenBytes)

	bs = make([]byte, headLen+bodyLen)
	// 放回去后边统一解析，这里也是可以优化的点
	copy(bs, headLenBytes)
	copy(bs[4:], bodyLenBytes)
	_, err = io.ReadFull(conn, bs[8:])
	return bs, err
}
