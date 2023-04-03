package demo

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/message"
)

// type Middleware func(next Proxy) Proxy

// type HandleFunc func(ctx context.Context, req *Request) (*Response, error)

type Proxy interface {
	Invoke(ctx context.Context, req *message.Request) (*message.Response, error)
}

// type Filter func(ctx context.Context, req *Request) (*Response, error)
//
// type FiltersProxy struct {
// 	Proxy
// 	filters []Filter
// }
//
// func (f FiltersProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
// 	for _, flt := range f.filters {
// 		resp, err := flt(ctx, req)
// 		if resp != nil {
// 			return resp, err
// 		}
// 	}
// 	res := f.Proxy.Invoke()
// }

