package beego

import "github.com/beego/beego/v2/server/web"

type UserController struct {
	web.Controller
}

func (c *UserController) GetUser() {
	c.Ctx.WriteString("你好，我是大明")
}

func (c *UserController) CreateUser() {
	u := &User{}
	err := c.Ctx.BindJSON(u)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	_ = c.Ctx.JSONResp(u)
}

type User struct {
	Name string
}