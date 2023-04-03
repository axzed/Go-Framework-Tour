package main

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/advance/ctx/graceful_shutdown/service"
	"log"
	"net/http"
	"time"
)

// 注意要从命令行启动，否则不同的 IDE 可能会吞掉关闭信号
func main() {
	s1 := service.NewServer("business", "localhost:8080")
	s1.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello"))
	}))
	s2 := service.NewServer("admin", "localhost:8081")
	app := service.NewApp([]*service.Server{s1, s2}, service.WithShutdownCallbacks(StoreCacheToDBCallback))
	app.StartAndServe()
}

func StoreCacheToDBCallback(ctx context.Context) {
	done := make(chan struct{}, 1)
	go func() {
		// 你的业务逻辑，比如说这里我们模拟的是将本地缓存刷新到数据库里面
		// 这里我们简单的睡一段时间来模拟
		log.Printf("刷新缓存中……")
		time.Sleep(1 * time.Second)
	}()
	select {
	case <-ctx.Done():
		log.Printf("刷新缓存超时")
	case <-done:
		log.Printf("缓存被刷新到了 DB")
	}
}
