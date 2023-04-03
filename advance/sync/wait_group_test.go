package sync

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"testing"
)

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	var result int64 = 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(delta int) {
			atomic.AddInt64(&result, int64(delta))
			defer wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println(result)
}

func TestErrgroup(t *testing.T) {
	eg := errgroup.Group{}
	var result int64 = 0
	for i := 0; i < 10; i++ {
		delta := i
		eg.Go(func() error {
			atomic.AddInt64(&result, int64(delta))
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
