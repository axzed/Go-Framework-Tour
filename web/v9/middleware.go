//go:build v9
package web

type Middleware func(next HandleFunc) HandleFunc
