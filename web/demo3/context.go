package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Context struct {
	// 用户不能使用 Req.Context()
	Req  *http.Request
	Resp http.ResponseWriter
	// Ctx  context.Context

	PathParams map[string]string

	// 缓存住你的响应
	RespStatusCode int
	// RespData []byte
	RespData     []byte
	MatchedRoute string

	tplEngine TemplateEngine

	UserValues map[string]any
}

func (ctx *Context) BindJSON(val any) error {
	if ctx.Req.Body == nil {
		return errors.New("web: body 为空")
	}
	decoder := json.NewDecoder(ctx.Req.Body)
	return decoder.Decode(val)
}

func (ctx *Context) FormValue(key string) (string, error) {
	err := ctx.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return ctx.Req.FormValue(key), nil
}

func (ctx *Context) FormValueOrDefault(key string, def string) string {
	val, err := ctx.FormValue(key)
	if err != nil || val == "" {
		return def
	}
	return val
}
func (ctx *Context) QueryValue(key string) (string, error) {
	params := ctx.Req.URL.Query()
	vals, ok := params[key]
	if !ok || len(vals) == 0 {
		return "", nil
	}

	return vals[0], nil
}

func (ctx *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx.RespData = bs
	ctx.RespStatusCode = code
	return err
}

func (ctx *Context) Render(tplName string, val any) error {
	data, err := ctx.tplEngine.Render(ctx.Req.Context(), tplName, val)
	if err != nil {
		return err
	}
	ctx.RespData = data
	ctx.RespStatusCode = 200
	return nil
}

func (ctx *Context) QueryValueV1(key string) StringValue {
	params := ctx.Req.URL.Query()
	vals, ok := params[key]
	if !ok || len(vals) == 0 {
		return StringValue{err: errors.New("key not found")}
	}

	return StringValue{val: vals[0]}
}

// ctx.QueryValueV2[int64]("age")
// func (ctx *Context) QueryValueV2[T any](key string) T {
//
// }

type StringValue struct {
	val string
	err error
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// func (s StringValue[T]) To() (T, error) {
//
// }
