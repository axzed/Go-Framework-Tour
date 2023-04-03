package rpc

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/message"
)

type Proxy interface {
	Invoke(ctx context.Context, req *message.Request) (*message.Response, error)
}

type Service interface {
	ServiceName() string
}
