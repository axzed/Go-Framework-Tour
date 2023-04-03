package gin

import "github.com/gin-gonic/gin"

type UserController struct {
}

func (c *UserController) GetUser(ctx *gin.Context) {
	panic("一些业务错误")
	ctx.String(200, "hello, world")
}
