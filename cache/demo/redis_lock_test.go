package demo

import (
	"context"
	"time"
)

func ExampleLockRefresh() {

	// 假如说我们拿到了一个锁
	var l *Lock

	stop := make(chan struct{}, 1)

	bizStop := make(chan struct{}, 1)
	go func() {
		// 间隔时间根据你的锁过期时间来决定
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()
		// 不断续约，直到收到退出信号
		ch := make(chan struct{}, 1)
		retryCnt := 0
		for {
			select {
			case <- ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cancel()
				// error 怎么处理
				// 可能需要对 err 分类处理

				// 超时了
				if err == context.DeadlineExceeded {
					// 可以重试
					// 如果一直重试失败，又怎么办？
					ch <- struct{}{}
					continue
				}
				if err != nil {
					// 不可挽回的错误
					// 你这里要考虑中断业务执行
					bizStop <- struct{}{}
					return
				}

				retryCnt = 0
			case <- ch:
				retryCnt ++
				// 重试信号
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err := l.Refresh(ctx)
				cancel()
				// error 怎么处理
				// 可能需要对 err 分类处理

				// 超时了
				if err == context.DeadlineExceeded {
					// 可以重试
					// 如果一直重试失败，又怎么办？
					if retryCnt > 10 {
						// 考虑中断业务
						return
					} else {
						ch <- struct{}{}
					}
					continue
				}
				if err != nil {
					// 不可挽回的错误
					// 你这里要考虑中断业务执行
					bizStop <- struct{}{}
					return
				}
				retryCnt = 0
			case <- stop:
				return
			}
		}
	}()

	// 这边就是你的业务

	for {
		select {
		case <- bizStop:
			// 要回滚的
			break
		default:
			// 你的业务，你的业务被拆成了好几个步骤，非常多的步骤
		}
	}

	// 业务结束
	// 这边是不是要通知不用再续约了
	stop <- struct{}{}

	// Output:
	// hello world
}