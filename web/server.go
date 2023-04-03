package web

import (
	"net/http"
	"strconv"
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
	addRoute(method string, path string, handler HandleFunc, ms...Middleware)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type ServerOption func(server *HTTPServer)

type HTTPServer struct {
	router
	tplEngine TemplateEngine
	log Logger
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

// func (s *HTTPServer) Use(mdls ...Middleware) {
// 	if s.mdls == nil {
// 		s.mdls = mdls
// 		return
// 	}
// 	s.mdls = append(s.mdls, mdls...)
// }

// Use 会执行路由匹配，只有匹配上了的 mdls 才会生效
// 这个只需要稍微改造一下路由树就可以实现
func (s *HTTPServer) Use(method, path string, mdls ...Middleware) {
	s.addRoute(method, path, nil, mdls...)
}

// UseAny 这个名字确实不咋的
func (s *HTTPServer) UseAny(path string, mdls ...Middleware) {
	s.addRoute(http.MethodGet, path, nil, mdls...)
	s.addRoute(http.MethodPost, path, nil, mdls...)
	s.addRoute(http.MethodOptions, path, nil, mdls...)
	s.addRoute(http.MethodConnect, path, nil, mdls...)
	s.addRoute(http.MethodDelete, path, nil, mdls...)
	s.addRoute(http.MethodHead, path, nil, mdls...)
	s.addRoute(http.MethodPatch, path, nil, mdls...)
	s.addRoute(http.MethodPut, path, nil, mdls...)
	s.addRoute(http.MethodTrace, path, nil, mdls...)
}
// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:       request,
		Resp:      writer,
		tplEngine: s.tplEngine,
	}

	// ctx pool.Get()
	// defer func(){
	//     ctx.Reset()
	//     pool.Put(ctx)
	// }
	s.serve(ctx)
}

// Start 启动服务器，编程接口
// 要求用户自己去配置文件读端口
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
	if mi.n != nil {
		ctx.PathParams = mi.pathParams
		ctx.MatchedRoute = mi.n.route
	}
	// 最后一个应该是执行用户代码
	var root HandleFunc = func(ctx *Context) {
		if !ok || mi.n == nil || mi.n.handler == nil {
			ctx.RespStatusCode = 404
			return
		}
		mi.n.handler(ctx)
	}
	// 从后往前组装
	for i := len(mi.mdls) - 1; i >= 0; i-- {
		root = mi.mdls[i](root)
	}
	// 第一个应该是回写响应的
	// 因为它在调用next之后才回写响应，
	// 所以实际上 flashResp 是最后一个步骤
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (s *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	ctx.Resp.Header().Set("Content-Length", strconv.Itoa(len(ctx.RespData)))
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		// s.log.Fatalln("回写响应失败", err)
		defaultLogger.Fatalln("回写响应失败", err)
	}
}

var defaultLogger Logger

func SetDefaultLogger(log Logger) {
	defaultLogger=log
}
type Logger interface {
	Fatalln(msg string, args...any)
}


type HTTPServerV1 struct {
	router
	tplEngine TemplateEngine
	log Logger
}

func NewHTTPServerV1(cfgFile string) *HTTPServerV1 {
	// 这里去读取配置文件
	// 初始化实例
	panic("implement me")
}