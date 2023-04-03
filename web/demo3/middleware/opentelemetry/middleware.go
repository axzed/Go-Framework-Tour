package opentelemetry

import (
	web "gitee.com/geektime-geekbang/geektime-go/web/demo3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			spanCtx, span := m.Tracer.Start(ctx.Req.Context(), "Unknown")
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			ctx.Req = ctx.Req.WithContext(spanCtx)
			defer func() {
				if ctx.MatchedRoute != "" {
					// 不是 404
					span.SetName(ctx.MatchedRoute)
				}
				span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))

			}()

			// ctx.Ctx = spanCtx
			next(ctx)

		}
	}
}
