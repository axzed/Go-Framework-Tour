package web

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type ServerOption func(server *HTTPServer)

type HTTPServer struct {
	router
	ms        []Middleware
	tplEngine TemplateEngine
}

func NewHTTPServer(opts ...ServerOption) *HTTPServer {
	s := &HTTPServer{
		router: newRouter(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func ServerWithTemplateEngine(engine TemplateEngine) ServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = engine
	}
}
func (s *HTTPServer) Use(ms ...Middleware) {
	s.ms = ms
}

var Debug bool

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:       request,
		Resp:      writer,
		tplEngine: s.tplEngine,
		// Ctx:  request.Context(),
	}

	root := s.serve
	for i := len(s.ms) - 1; i >= 0; i-- {
		m := s.ms[i]
		root = m(root)
	}

	var flushMdl Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.writeResp(ctx)
		}
	}
	root = flushMdl(root)
	root(ctx)
}

func (s *HTTPServer) writeResp(ctx *Context) {
	ctx.Resp.WriteHeader(ctx.RespStatusCode)
	ctx.Resp.Write(ctx.RespData)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)

	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.RespStatusCode = 404
		ctx.RespData = []byte("Not Found")
		return
	}
	ctx.PathParams = mi.pathParams
	ctx.MatchedRoute = mi.n.route
	mi.n.handler(ctx)
}
