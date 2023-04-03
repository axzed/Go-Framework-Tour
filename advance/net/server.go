package net

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func Serve(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			handleConn(conn)
		}()
	}
}

func handleConn(conn net.Conn) {
	for {
		// 读数据
		bs := make([]byte, 8)
		_, err := conn.Read(bs)
		if err == io.EOF || err == net.ErrClosed ||
			err == io.ErrUnexpectedEOF {
			// 一般关闭的错误比较懒得管
			// 也可以把关闭错误输出到日志
			_ = conn.Close()
			return
		}
		if err != nil {
			continue
		}
		res := handleMsg(bs)
		_, err = conn.Write(res)
		if err == io.EOF || err == net.ErrClosed ||
			err == io.ErrUnexpectedEOF {
			_ = conn.Close()
			return
		}
	}
}

func handleConnV1(conn net.Conn) {
	for {
		// 读数据
		bs := make([]byte, 8)
		_, err := conn.Read(bs)
		if err != nil {
			// 一般关闭的错误比较懒得管
			// 也可以把关闭错误输出到日志
			_ = conn.Close()
			return
		}
		res := handleMsg(bs)
		_, err = conn.Write(res)
		if err != nil {
			_ = conn.Close()
			return
		}
	}
}

func handleMsg(bs []byte) []byte {
	return []byte("world")
}

type Server struct {
	addr string
}

func (s *Server) StartAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			// 直接在这里处理
			er := s.handleConn(conn)
			if er != nil {
				_ = conn.Close()
				fmt.Printf("connect error: %v", er)
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		// 读数据长度
		bs := make([]byte, lenBytes)
		_, err := conn.Read(bs)
		if err != nil {
			return err
		}

		reqBs := make([]byte, binary.BigEndian.Uint64(bs))
		_, err = conn.Read(reqBs)
		if err != nil {
			return err
		}
		res := string(reqBs) + ", from response"
		// 总长度
		bs = make([]byte, lenBytes, len(res)+lenBytes)
		// 写入消息长度
		binary.BigEndian.PutUint64(bs, uint64(len(res)))
		bs = append(bs, res...)
		_, err = conn.Write(bs)
		if err != nil {
			return err
		}
	}
}
