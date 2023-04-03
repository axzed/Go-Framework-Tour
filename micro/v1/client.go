package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/silenceper/pool"
	"net"
)

type Client struct {
	connPool pool.Pool
	invoker  Proxy
}

type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}

func NewClient(address string) (*Client, error) {
	// Create a connection pool: Initialize the number of connections to 5, the maximum idle connection is 20, and the maximum concurrent connection is 30
	poolConfig := &pool.Config{
		InitialCap:  5,
		MaxIdle:     20,
		MaxCap:      30,
		Factory:     func() (interface{}, error) { return net.Dial("tcp", address) },
		Close:       func(v interface{}) error { return v.(net.Conn).Close() },
		IdleTimeout: time.Minute,
	}
	connPool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		return nil, err
	}
	res := &Client{
		connPool: connPool,
	}
	return res, nil
}

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	conn, err := c.connPool.Get()
	if err != nil {
		return nil, fmt.Errorf("client: 获得获取一个可用连接 %w", err)
	}
	// put back
	defer c.connPool.Put(conn)
	cn := conn.(net.Conn)
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("client: 无法序列化请求, %w", err)
	}

	encode := EncodeMsg(bs)
	_, err = cn.(net.Conn).Write(encode)
	if err != nil {
		return nil, err
	}

	bs, err = ReadMsg(cn.(net.Conn))
	if err != nil {
		return nil, fmt.Errorf("client: 无法读取响应 %w", err)
	}

	resp := &Response{}
	err = json.Unmarshal(bs, resp)
	return resp, nil
}
