package demo

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	// 开始监听端口
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}
	for {
		// 这边开始接收连接
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			handle(conn)
		}()
	}
}

func handle(conn net.Conn) {
	for {
		lenBs := make([]byte, 8)
		_, err := conn.Read(lenBs)
		if err != nil {
			conn.Close()
			return
		}
		msgLen := binary.BigEndian.Uint64(lenBs)
		reqBs := make([]byte, msgLen)
		_, err = conn.Read(reqBs)

		if err != nil {
			conn.Close()
			return
		}

		_, err = conn.Write([]byte("hello, world"))
		if err != nil {
			conn.Close()
			return
		}

	}
}
