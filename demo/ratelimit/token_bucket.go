package ratelimit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

type TokenBucketLimiter struct {
	tokens chan struct{}
	closed chan struct{}
}

// NewTokenBucketLimiter buffer 最多积攒多少个令牌
// interval 就是间隔多久产生一个令牌
func NewTokenBucketLimiter(buffer int, interval time.Duration) *TokenBucketLimiter {
	res :=  &TokenBucketLimiter{
		tokens: make(chan struct{}, buffer),
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		select {
		case <- ticker.C:
			res.tokens <- struct{}{}
		case <- res.closed:
			close(res.tokens)
			return
		}
		// for range ticker.C {
		//
		// 	res.tokens <- struct{}{}
		//
		// 	// 这个地方你可能放满
		// 	// select {
		// 	// case res.tokens <- struct{}{}:
		// 	// default:
		// 	//
		// 	// }
		// }
	}()
	return res
}

func (t TokenBucketLimiter) BuildUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// select {
		// case <- ctx.Done():
		// 	// 缺陷是 channel 是 FIFO 的
		// 	// 意味着等待最久的，会拿到令牌
		// 	// 这意味着，你大概率在业务处理的时候会超时。要小心不同超时时间设置
		// 	return ctx.Err()
		// case _, ok := <- t.tokens:
		// 	if ok {
		// 		return invoker(ctx, method, req, reply, cc, opts...)
		// 	}
		// }

		// 怎么样处理？
		select {
		case _, ok := <- t.tokens:
			if ok {
				return handler(ctx, req)
			}
		default:
			// 拿不到令牌就直接拒绝
		}

		// 熔断限流降级之间区别在这里了
		// 1. 返回默认值 get_user -> GetUserResp
		// 2. 打个标记位，后面执行快路径，或者兜底路径
		return nil, errors.New("你被限流了")

	}
}

// func (t TokenBucketLimiter) BuildUnary() grpc.UnaryClientInterceptor {
//
// }

// 老生常谈的多次 Close 的问题
func (t TokenBucketLimiter) Close() error {
	t.closed <- struct{}{}
	return nil
}
