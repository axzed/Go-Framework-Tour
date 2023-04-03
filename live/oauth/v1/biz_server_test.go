package v1

import (
	"gitee.com/geektime-geekbang/geektime-go/web"
	"gitee.com/geektime-geekbang/geektime-go/web/session"
	"gitee.com/geektime-geekbang/geektime-go/web/session/cookie"
	"gitee.com/geektime-geekbang/geektime-go/web/session/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"html/template"
	"net/http"
	"testing"
	"time"
)

func TestBizServer(t *testing.T) {
	tpl, err := template.ParseGlob("template/*.gohtml")
	require.NoError(t, err)
	engine := &web.GoTemplateEngine{
		T: tpl,
	}
	server := web.NewHTTPServer(web.ServerWithTemplateEngine(engine))
	sessMgr := session.Manager{
		Store:      memory.NewStore(time.Minute * 15),
		Propagator: cookie.NewPropagator("sso_sess"),
		SessCtxKey: "sso_sess",
	}

	server.UseAny("/", func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 登录界面或者登录请求
			if ctx.Req.URL.Path == "/login" {
				next(ctx)
				return
			}
			_, err := sessMgr.GetSession(ctx)
			if err != nil {
				ctx.Redirect("login")
				return
			}
			next(ctx)
		}
	})

	server.Get("/login", func(ctx *web.Context) {
		_ = ctx.Render("login.gohtml", nil)
	})

	server.Post("/login", func(ctx *web.Context) {
		email, err := ctx.FormValue("email").String()
		if err != nil {
			_ = ctx.RespString(http.StatusBadRequest, "邮箱输入错误")
			return
		}
		password, err := ctx.FormValue("password").String()
		if err != nil {
			_ = ctx.RespString(http.StatusBadRequest, "密码错误")
			return
		}
		if password == "123" && email == "123@demo.com" {
			ssid := uuid.New().String()
			_, err = sessMgr.InitSession(ctx, ssid)
			if err != nil {
				_ = ctx.RespServerError("登录失败")
				return
			}
			ctx.Redirect("/profile")
			return
		}
		_ = ctx.RespString(http.StatusBadRequest, "用户名密码不对")
	})

	// 假如说我们登录成功之后我们就访问对应的资源
	// 这个就是模拟登录后的请求
	server.Get("/profile", func(ctx *web.Context) {
		_ = ctx.RespOk("hello, world")
	})
	server.Start(":8082")
}
