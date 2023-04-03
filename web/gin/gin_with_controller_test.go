package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestUserController_GetUser(t *testing.T) {
	g := gin.Default()
	ctrl := &UserController{}
	g.GET("/user/*", ctrl.GetUser)
	g.POST("/user/*", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello %s", "world")
	})

	g.GET("/static", func(context *gin.Context) {
		// 读文件
		// 谐响应
	})
	_ = g.Run(":8082")
}
