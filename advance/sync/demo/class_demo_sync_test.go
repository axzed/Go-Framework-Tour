package demo

import (
	"fmt"
	"testing"
	"time"
)

func TestDeferRLock(t *testing.T) {
	sm := SafeMap[string, string]{
		values: make(map[string]string, 4),
	}
	sm.LoadOrStoreV1("a", "b")
	fmt.Println("hello")
}

func TestOverride(t *testing.T) {
	sm := SafeMap[string, string]{
		values: make(map[string]string, 4),
	}
	go func() {
		time.Sleep(time.Second)
		sm.LoadOrStoreV2("a", "b")
	}()

	go func() {
		time.Sleep(time.Second)
		sm.LoadOrStoreV2("a", "c")
	}()

	go func() {
		time.Sleep(time.Second)
		sm.LoadOrStoreV2("a", "d")
	}()

	go func() {
		time.Sleep(time.Second)
		sm.LoadOrStoreV2("a", "e")
	}()
	time.Sleep(time.Second)
	fmt.Println("hello")
}
