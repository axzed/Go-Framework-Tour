package handler

// 如果是复杂的项目，那么 vo 可以是一个单独的项目，或者 web 下的 handler 也可以继续垂直拆分

type loginReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
type signUpReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPwd string `json:"confirm_pwd"`
}

type User struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
}

type Resp struct {
	Msg string `json:"msg"`
	Data any `json:"data"`
}
