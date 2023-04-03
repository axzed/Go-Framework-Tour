package net

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Option func(p *SimplePool)

type SimplePool struct {
	idleChan    chan conn
	waitChan chan *conReq

	factory     func() (net.Conn, error)
	idleTimeout time.Duration

	maxCnt int32
	// 连接数
	cnt int32

	l sync.Mutex
}

func NewSimplePool(factory func()(net.Conn, error), opts...Option) *SimplePool {
	res := &SimplePool {
		idleChan: make(chan conn, 16),
		waitChan: make(chan *conReq, 128),
		factory: factory,
		maxCnt: 128,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (p *SimplePool) Get() (net.Conn, error) {
	for {
		select {
		case c := <-p.idleChan:
			// 超时，直接关闭.
			// 有没有觉得奇怪，就是明明我们就是需要一个连接，但是我们还关闭了
			if c.lastActive.Add(p.idleTimeout).Before(time.Now()) {
				atomic.AddInt32(&p.cnt, -1)
				_ = c.c.Close()
				continue
			}
			return c.c, nil
		default:
			cnt := atomic.AddInt32(&p.cnt, 1)
			if cnt <= p.maxCnt {
				return p.factory()
			}
			atomic.AddInt32(&p.cnt, -1)
			req := &conReq{
				con: make(chan conn, 1),
			}
			// 可能阻塞在这两句，对应不同的情况。
			// 所以实际上 waitChan 根本不需要设计很大的容量
			// 另外，这里需不需要加锁？
			p.waitChan <- req
			c := <- req.con
			return c.c, nil
		}
	}
}

func (p *SimplePool) Put(c net.Conn) {
	// 为什么我只在这个部分加锁，其余部分都不加？
	p.l.Lock()
	if len(p.waitChan) > 0 {
		req := <- p.waitChan
		p.l.Unlock()
		req.con <- conn{c: c, lastActive: time.Now()}
		return
	}

	p.l.Unlock()

	select {
	case p.idleChan <- conn{c: c, lastActive: time.Now()}:
	default:
		defer func() {
			atomic.AddInt32(&p.maxCnt, -1)
		}()
		_ = c.Close()
	}
}

// WithMaxIdleCnt 自定义最大空闲连接数量
func WithMaxIdleCnt(maxIdleCnt int32) Option {
	return func(p *SimplePool) {
		p.idleChan = make(chan conn, maxIdleCnt)
	}
}

// WithMaxCnt 自定义最大连接数量
func WithMaxCnt(maxCnt int32) Option {
	return func(p *SimplePool) {
		p.maxCnt = maxCnt
	}
}

type conn struct {
	c          net.Conn
	lastActive time.Time
}

type conReq struct {
	con chan conn
}