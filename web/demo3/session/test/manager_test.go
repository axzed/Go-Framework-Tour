package test

import (
	web "gitee.com/geektime-geekbang/geektime-go/web/demo3"
	"gitee.com/geektime-geekbang/geektime-go/web/demo3/session"
	"gitee.com/geektime-geekbang/geektime-go/web/demo3/session/cookie"
	"gitee.com/geektime-geekbang/geektime-go/web/demo3/session/memory"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	s := web.NewHTTPServer()

	m := session.Manager{
		SessCtxKey: "_sess",
		Store:      memory.NewStore(30 * time.Minute),
		Propagator: cookie.NewPropagator("sessid",
			cookie.WithCookieOption(func(c *http.Cookie) {
				c.HttpOnly = true
			})),
	}

	s.Post("/login", func(ctx *web.Context) {
		// 前面就是你登录的时候一大堆的登录校验

		id := uuid.New()
		sess, err := m.InitSession(ctx, id.String())
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		// 然后根据自己的需要设置
		err = sess.Set(ctx.Req.Context(), "mykey", "some value")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
	})
	s.Get("/resource", func(ctx *web.Context) {
		sess, err := m.GetSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		val, err := sess.Get(ctx.Req.Context(), "mykey")
		ctx.RespData = []byte(val)
		ctx.RespStatusCode = 200
	})

	s.Post("/logout", func(ctx *web.Context) {
		_ = m.RemoveSession(ctx)
	})

	s.Use(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 执行校验
			if ctx.Req.URL.Path != "/login" {
				sess, err := m.GetSession(ctx)
				// 不管发生了什么错误，对于用户我们都是返回未授权
				if err != nil {
					ctx.RespStatusCode = http.StatusUnauthorized
					return
				}
				_ = m.Refresh(ctx.Req.Context(), sess.ID())
			}
			next(ctx)
		}
	})

	s.Start(":8081")
}
