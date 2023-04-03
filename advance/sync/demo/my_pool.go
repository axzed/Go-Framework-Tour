package demo

import (
	"sync"
	"unsafe"
)

type MyPool struct {
	p      sync.Pool
	maxCnt int32
	cnt    int32
}

func (p *MyPool) Get() any {
	return p.p.Get()
}

func (p *MyPool) Put(val any) {
	// 大对象不放回去
	if unsafe.Sizeof(val) > 1024 {
		return
	}

	p.p.Put(val)
}
