//go:build v6
package web

type Middleware func(next HandleFunc) HandleFunc
