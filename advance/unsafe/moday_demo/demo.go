package moday_demo

type Handler interface {
	Handle()
}

type HandlerBasedOnTree struct {
}

type privateHandler struct {
}

// func NewPrivateHandler() Handler {
//
// }

// 个人偏好
// func NewHandlerBasedOnTreeV1() Handler {
//
// }

// Go 推荐的，返回具体类型
// func NewHandlerBasedOnTreeV2() *HandlerBasedOnTree {
//
// }

//
// func NewHandlerBasedOnTreeV3() HandlerBasedOnTree {
//
// }

type MyService struct {
	// 使用接口，可以注入 mock 的实现
	handler Handler
}

type MyHandler struct {
}

func (m *MyHandler) Handle() {
	// 假如说这里调用一个下游
	// 可以是 http
	// 可以是 RPC
}

func (ms *MyService) Serve() {
	ms.handler.Handle()
}
