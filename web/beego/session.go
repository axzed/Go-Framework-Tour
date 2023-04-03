package beego

import "github.com/beego/beego/v2/server/web"

// to test this
// you should browse POST http://localhost:8080/session to set session firstly
func main() {

	// enable session
	// or you can put "session=true" into your config file
	// or you can set SessionGCMaxLifetime
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 10

	// create contr
	ctrl := &MainController{}

	// POST http://localhost:8080/session => ctrl.PutSession()
	web.Router("/session", ctrl, "post:PutSession")

	// GET http://localhost:8080/session => ctrl.ReadSession()
	web.Router("/session", ctrl, "get:ReadSession")

	// GET http://localhost:8080/session => ctrl.ReadSession()
	web.Router("/session", ctrl, "delete:DeleteSession")

	web.Run()
}

type MainController struct {
	web.Controller
}

func (ctrl *MainController) PutSession() {
	// put something into session
	ctrl.SetSession("name", "web session")

	// web-example/views/hello_world.html
	ctrl.TplName = "hello_world.html"
	ctrl.Data["name"] = "PutSession"
	_ = ctrl.Render()
}

func (ctrl *MainController) ReadSession() {
	// web-example/views/hello_world.html
	ctrl.TplName = "hello_world.html"
	ctrl.Data["name"] = ctrl.GetSession("name")
	// don't forget this
	_ = ctrl.Render()
}

func (ctrl *MainController) DeleteSession() {
	// delete session all
	ctrl.DestroySession()
	// web-example/views/hello_world.html
	ctrl.TplName = "hello_world.html"
	_ = ctrl.Render()
}
