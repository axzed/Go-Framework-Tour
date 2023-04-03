//go:build v4
package web

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
	PathParams map[string]string
}
