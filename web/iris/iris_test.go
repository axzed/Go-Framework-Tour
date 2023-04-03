package iris

import (
	"github.com/kataras/iris/v12"
	"testing"
)

func TestHelloWorld(t *testing.T) {

	// g := gin.Default()
	// ctrl := &UserController{}
	// g.GET("/user/*", ctrl.GetUser)
	// g.POST("/user/*", func(ctx *gin.Context) {
	// 	ctx.String(http.StatusOK, "hello %s", "world")
	// })
	//
	// g.GET("/static", func(context *gin.Context) {
	// 	// 读文件
	// 	// 谐响应
	// })

	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		_, _ = ctx.HTML("Hello <strong>%s</strong>!", "World")
	})

	_ = app.Listen(":8083")
}
