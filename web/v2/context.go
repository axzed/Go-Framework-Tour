//go:build v2
package web

import "net/http"

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
}
