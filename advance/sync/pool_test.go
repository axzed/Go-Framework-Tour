package sync

import (
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	p := sync.Pool{
		New: func() interface{} {
			// 创建函数，sync.Pool 会回调
			return nil
		},
	}

	obj := p.Get()
	// 在这里使用取出来的对象
	// 用完再还回去
	p.Put(obj)
}
