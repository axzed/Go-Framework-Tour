package demo

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}
	msg := "how are you"
	msgLen := len(msg)
	// msgLen how are you
	msgLenBs := make([]byte, 8)
	binary.BigEndian.PutUint64(msgLenBs, uint64(msgLen))
	data := append(msgLenBs, []byte(msg)...)
	_, err = conn.Write(data)
	if err != nil {
		conn.Close()
		return
	}

	respBs := make([]byte, 16)
	_, err = conn.Read(respBs)
	if err != nil {
		conn.Close()
	}
}
