//go:build v7
package web

type Middleware func(next HandleFunc) HandleFunc
