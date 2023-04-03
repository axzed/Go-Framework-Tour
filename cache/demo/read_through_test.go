package demo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
	"gitee.com/geektime-geekbang/geektime-go/orm"
)

func TestReadThroughCache_Get(t *testing.T) {
	local := NewLocalCache(func(key string, val any) {

	})
	var db *orm.DB
	cache := &ReadThroughCache{
		Cache: local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (any, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// Key = /user/123 ，其中123 是用户 id
				// 这是用户的
				id := strings.Trim(key, "/user/")
				return orm.NewSelector[User](db).Where(orm.C("Id").EQ(id)).Get(ctx)
			} else if strings.HasPrefix(key, "/order/") {
				// 找 Order 数据
			} else if strings.HasPrefix(key, "produce") {
				// 找商品的数据
			}
			// if-else 就没完没了了
			return nil, errors.New("不支持操作")
		},
	}

	cache.Get(context.Background(), "/user/123")

	userCache := &ReadThroughCache{
		Cache: local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (any, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// Key = /user/123 ，其中123 是用户 id
				// 这是用户的
				id := strings.Trim(key, "/user/")
				return orm.NewSelector[User](db).Where(orm.C("Id").EQ(id)).Get(ctx)
			}
			// if-else 就没完没了了
			return nil, errors.New("不支持操作")
		},
	}

	userCache.Get(context.Background(), "/user/123")
	// orderCache

	userCacheV1 := &ReadThroughCacheV1[*User]{
		Cache: local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (*User, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// Key = /user/123 ，其中123 是用户 id
				// 这是用户的
				id := strings.Trim(key, "/user/")
				return orm.NewSelector[User](db).Where(orm.C("Id").EQ(id)).Get(ctx)
			}
			// if-else 就没完没了了
			return nil, errors.New("不支持操作")
		},
	}

	val, err := userCacheV1.Get(context.Background(), "/user/123")
	// val 还是 any, 我干嘛用泛型？？？？？我干嘛要 v1??
	fmt.Println(val)
	fmt.Println(err)


	userCacheV2 := &ReadThroughCacheV2[*User]{
		// 这边要考虑创建一个 CacheV2
		// Cache: local,
		Expiration: time.Minute,
		LoadFunc: func(ctx context.Context, key string) (*User, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// Key = /user/123 ，其中123 是用户 id
				// 这是用户的
				id := strings.Trim(key, "/user/")
				return orm.NewSelector[User](db).Where(orm.C("Id").EQ(id)).Get(ctx)
			}
			// if-else 就没完没了了
			return nil, errors.New("不支持操作")
		},
	}

	user, err := userCacheV2.Get(context.Background(), "/user/123")
	fmt.Println(user.Name)
	fmt.Println(err)


	userCacheV3 := &ReadThroughCacheV3{
		Loader: LoadFunc(func(ctx context.Context, key string) (any, error) {
			if strings.HasPrefix(key, "/user/") {
				// 找用户的数据
				// Key = /user/123 ，其中123 是用户 id
				// 这是用户的
				id := strings.Trim(key, "/user/")
				return orm.NewSelector[User](db).Where(orm.C("Id").EQ(id)).Get(ctx)
			}
			// if-else 就没完没了了
			return nil, errors.New("不支持操作")
		}),
	}
	fmt.Println(userCacheV3)
}
