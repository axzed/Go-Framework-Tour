package web

import (
	"fmt"
	"sync"
	"testing"
)

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	Debug = true
	s := NewHTTPServer()
	// s.Use(repeat_body.Middleware(), accesslog.MiddlewareBuilder{}.Build())
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		// age, err := ctx.QueryValueV1("age").ToInt64()
		ctx.Resp.Write([]byte("hello, user"))

	})

	// s.Post("/upload", (&FileUploader{}).Upload)

	s.Post("/user", func(ctx *Context) {
		u := &User{}
		err := ctx.BindJSON(u)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(u)
	})

	// / -----------------------------------------
	s.Start(":8081")
}

type User struct {
	Name string `json:"name"`
}

// type SafeHTTPServer struct {
// 	Server
// 	l sync.RWMutex
// }
//
// func (s *SafeHTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
// 	s.l.RLock()
// 	defer s.l.Unlock()
// 	s.Server.ServeHTTP()
// }

type SafeContext struct {
	c Context
	l sync.RWMutex
}

func (ctx *SafeContext) RespJSON(code int, val any) error {
	ctx.l.Lock()
	defer ctx.l.Unlock()
	return ctx.c.RespJSON(code, val)
}
