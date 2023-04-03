package main

import (
	"gitee.com/geektime-geekbang/geektime-go/web"
	"gitee.com/geektime-geekbang/geektime-go/web/middleware/accesslog"
	"gitee.com/geektime-geekbang/geektime-go/web/middleware/cors"
	"gitee.com/geektime-geekbang/geektime-go/web/middleware/opentelemetry"
	"gitee.com/geektime-geekbang/geektime-go/web/middleware/prometheus"
	"gitee.com/geektime-geekbang/geektime-go/web/middleware/recovery"
	"go.uber.org/zap"
	"net/http"
)

func initSever() *web.HTTPServer {
	server := web.NewHTTPServer()
	// 我把这些 builder 的接收器都改成了结构体，懒得写一个括号

	// 对 VIP 鉴权
	// server.UseAny("/vip", func(next web.HandleFunc) web.HandleFunc {
	//
	// })
	//
	// server.UseAny("/login", func(next web.HandleFunc) web.HandleFunc {
	// 	// 针对登录的限流
	// })

	// 这三个其实不太好确定谁先谁后，你们可以自己琢磨一下自己
	server.UseAny("/*",
		// func(next web.HandleFunc) web.HandleFunc {
		// 	return func(ctx *web.Context) {
		// 		做熔断限流降级
		// 	}
		// },
		// func(next web.HandleFunc) web.HandleFunc {
		// 	// 返回 mock 响应，在开始研发阶段
		// 	return func(ctx *web.Context) {
		// 		// 连上你的 Mock 中心，根据 URL 来找到 mock 数据
		// 		ctx.RespData = []byte("mock 数据")
		// 		ctx.RespStatusCode = 200
		// 	}
		// },
		opentelemetry.MiddlewareBuilder{}.Build(),
		accesslog.NewBuilder().LogFunc(func(accessLog string) {
			zap.L().Info(accessLog)
		}).Build(),
		recovery.MiddlewareBuilder{
			StatusCode: http.StatusInternalServerError,
			ErrMsg: "系统异常",
			LogFunc: func(ctx *web.Context, err any) {
				zap.L().Error("服务 panic", zap.Any("panic", err),
					// 发生 panic 的时候，可能都还没到路由查找那里
					zap.String("route", ctx.MatchedRoute))
			},
		}.Build(),
		prometheus.MiddlewareBuilder{
			Name: "userapp",
			Subsystem: "web",
			// 可以考虑在这里设置 instance id 之类的东西
			// ConstLabels: map[string]string{},
			Help: "userapp 的 web 统计",
		}.Build(),
		cors.MiddlewareBuilder{}.Build())
	return server
}