package demo

import (
	"context"
	"time"
)

func ExampleCacheV2() {
	var userCache CacheV2[User]
	// userCache.Set(context.Background(), "my-order-01", Order{}, time.Minute)
	userCache.Set(context.Background(), "my-user-01", User{}, time.Minute)
	var orderCache CacheV2[Order]
	orderCache.Set(context.Background(), "my-order-01", Order{}, time.Minute)
}

type User struct {
	Name string
}

type Order struct {

}