package demo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_addRoute(t *testing.T) {
	tests := []struct {

		// 输入
		method string
		path   string
	}{
		// 静态匹配
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "//home",
		},
		{
			method: http.MethodGet,
			path:   "//home1///",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/detail/profile",
		},
		{
			method: http.MethodGet,
			path:   "/order/cancel",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:order_sn",
		},

		{
			method: http.MethodPost,
			path:   "/order/cancel",
		},
	}

	var handleFunc HandleFunc = func(context *Context) {

	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: handleFunc,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: handleFunc,
						children: map[string]*node{
							"detail": &node{
								path: "detail",
								children: map[string]*node{
									"profile": &node{
										path:    "profile",
										handler: handleFunc,
									},
								},
							},
						},
					},
					"order": &node{
						path: "order",
						starChild: &node{
							path:    "*",
							handler: handleFunc,
						},
						children: map[string]*node{
							"cancel": &node{
								path:    "cancel",
								handler: handleFunc,
							},
							"detail": &node{
								path: "detail",
								paramChild: &node{
									path:    ":order_sn",
									handler: handleFunc,
								},
							},
						},
					},
					"home": &node{
						path:    "home",
						handler: handleFunc,
					},
					"home1": &node{
						path:    "home1",
						handler: handleFunc,
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"cancel": &node{
								path:    "cancel",
								handler: handleFunc,
							},
						},
					},
				},
			},
		},
	}

	res := &router{
		trees: map[string]*node{},
	}
	for _, tc := range tests {
		res.addRoute(tc.method, tc.path, handleFunc)
	}
	errStr, ok := wantRouter.equal(res)
	assert.True(t, ok, errStr)

	findCases := []struct {
		name   string
		method string
		path   string

		found      bool
		wantPath   string
		hasHandler bool
	}{
		{
			name:       "/",
			method:     http.MethodGet,
			path:       "/",
			found:      true,
			wantPath:   "/",
			hasHandler: true,
		},
		{
			name:       "/user",
			method:     http.MethodGet,
			path:       "/user",
			found:      true,
			wantPath:   "user",
			hasHandler: true,
		},
		{
			name:       "/user/detail",
			method:     http.MethodGet,
			path:       "/user/detail",
			found:      true,
			wantPath:   "detail",
			hasHandler: false,
		},
		{
			name:       "/order/*",
			method:     http.MethodGet,
			path:       "/order/abc",
			found:      true,
			wantPath:   "*",
			hasHandler: true,
		},
		{
			name:       "/order/*",
			method:     http.MethodGet,
			path:       "/order/detail/:order_sn",
			found:      true,
			wantPath:   ":order_sn",
			hasHandler: true,
		},
	}
	for _, tc := range findCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, ok := res.findRoute(http.MethodGet, tc.path)
			assert.Equal(t, tc.found, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantPath, mi.n.path)
			assert.Equal(t, tc.hasHandler, mi.n.handler != nil)
		})
	}

}

func (r router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}
