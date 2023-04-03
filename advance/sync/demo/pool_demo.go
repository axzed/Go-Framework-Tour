package demo

import (
	"fmt"
	"sync"
)

type MyCache struct {
	pool sync.Pool
}

func NewMyCache() *MyCache {
	return &MyCache{
		pool: sync.Pool{
			New: func() any {
				fmt.Println("hhh, new")
				return []byte{}
			},
		},
	}
}
