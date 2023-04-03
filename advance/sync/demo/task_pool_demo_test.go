package demo

import (
	"fmt"
	"testing"
	"time"
)

func TestTaskPool_Do(t1 *testing.T) {
	tp := NewTaskPool(2)
	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task1")
	})

	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task2")
	})

	tp.Do(func() {
		MyTask(1, "13")
	})
}

func TestTaskPoolWithCache_Do(t1 *testing.T) {
	tp := NewTaskPoolWithCache(2, 10)
	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task1")
	})

	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task2")
	})

	id := 1
	name := "Tom"
	tp.Do(func() {
		MyTask(id, name)
	})

	time.Sleep(2 * time.Second)
}

func MyTask(a int, b string) {
	//
}
