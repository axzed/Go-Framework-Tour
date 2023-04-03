package web

type Middleware func(next HandleFunc) HandleFunc
