package demo

import "strings"

// 代表路由
type router struct {
	// trees 代表的是森林， HTTP method => 树的根节点
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		// 根节点
		root = &node{path: "/"}
		r.trees[method] = root
	}
	if path == "/" {
		root.handler = handleFunc
		return
	}

	// 把前后的 / 都去掉
	path = strings.Trim(path, "/")

	cur := root
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		cur = cur.childOrCreate(seg)
	}
	cur.handler = handleFunc
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{n: root}, true
	}

	path = strings.Trim(path, "/")

	cur := root
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if cur.children == nil {
			if cur.paramChild != nil {
				mi := &matchInfo{
					n: cur.paramChild,
					pathParams: map[string]string{
						cur.paramChild.path[1:]: seg,
					},
				}
				return mi, true
			}
			return &matchInfo{n: cur.starChild}, cur.starChild != nil
		}
		child, ok := cur.children[seg]
		if !ok {
			if cur.paramChild != nil {
				mi := &matchInfo{
					n: cur.paramChild,
					pathParams: map[string]string{
						cur.paramChild.path[1:]: seg,
					},
				}
				return mi, true
			}
			return &matchInfo{n: cur.starChild}, cur.starChild != nil
		}
		cur = child
	}
	return &matchInfo{n: cur}, true
}

// childOrCreate 查找子节点，如果子节点不存在就创建一个
// 并且将子节点放回去了 children 中
func (n *node) childOrCreate(path string) *node {
	// /a/*/c
	if path == "*" {
		if n.starChild == nil {
			n.starChild = &node{
				path: path,
			}
		}
		return n.starChild
	}
	if path[0] == ':' {
		if n.paramChild == nil {
			n.paramChild = &node{
				path: path,
			}
		}
		return n.paramChild
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

type node struct {
	// /a/b/c 中的 b 这一段
	path    string
	handler HandleFunc

	// path => 到子节点的映射
	// 静态路由
	children map[string]*node
	// children []*node

	// 通配符匹配
	starChild *node

	paramChild *node
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
