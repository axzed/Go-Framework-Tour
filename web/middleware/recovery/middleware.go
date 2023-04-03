package recovery

import (
	"gitee.com/geektime-geekbang/geektime-go/web"
)

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *web.Context, err any)
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = []byte(m.ErrMsg)
					// 万一 LogFunc 也panic，那我们也无能为力了
					m.LogFunc(ctx, err)
				}
			}()

			// 这里就是before route, before execute
			next(ctx)
			// 这里就是after route, after execute
		}
	}
}
