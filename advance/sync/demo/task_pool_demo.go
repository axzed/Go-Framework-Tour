package demo

type TaskPool struct {
	ch chan struct{}
}

func NewTaskPool(limit int) *TaskPool {
	t := &TaskPool{
		ch: make(chan struct{}, limit),
	}
	// 提前准备好了令牌
	for i := 0; i < limit; i++ {
		t.ch <- struct{}{}
	}
	return t
}

func (t *TaskPool) Do(f func()) {
	token := <-t.ch
	// 异步执行
	go func() {
		f()
		t.ch <- token
	}()

	// 同步执行
	// f()
	// t.ch <- token
}

type TaskPoolWithCache struct {
	cache chan func()
}

func NewTaskPoolWithCache(limit int, cacheSize int) *TaskPoolWithCache {
	t := &TaskPoolWithCache{
		cache: make(chan func(), cacheSize),
	}
	// 直接把 goroutine 开好
	for i := 0; i < limit; i++ {
		go func() {
			for {
				// 在 goroutine 里面不断尝试从 cache 里面拿到任务
				select {
				case task, ok := <-t.cache:
					if !ok {
						return
					}
					task()
				}
			}
		}()
	}
	return t
}

func (t *TaskPoolWithCache) Do(f func()) {
	t.cache <- f
}

// 显式控制生命周期
// func (t *TaskPoolWithCache) Start() {
// 	for i := 0; i < t.limit; i++ {
// 		go func() {
// 			for {
// 				// 在 goroutine 里面不断尝试从 cache 里面拿到任务
// 				select {
// 				case task, ok := <-t.cache:
// 					if !ok {
// 						return
// 					}
// 					task()
// 				}
// 			}
// 		}()
// 	}
// }
