//go:build v8
package web

type Middleware func(next HandleFunc) HandleFunc
