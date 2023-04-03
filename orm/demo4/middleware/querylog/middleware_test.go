package querylog

import (
	"context"
	"fmt"
	orm "gitee.com/geektime-geekbang/geektime-go/orm/demo4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := &MiddlewareBuilder{}
	builder.LogFunc(func(sql string, args ...any) {
		fmt.Println(sql)
	}).SlowQueryThreshold(100) // 100 ms 就是慢查询
	db, err := orm.Open("sqlite3",
		"file:test.db?cache=shared&mode=memory", orm.DBWithMiddlewares(
			builder.Build(), func(next orm.Handler) orm.Handler {
				return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
					time.Sleep(time.Millisecond)
					return next(ctx, qc)
				}
			}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = orm.NewSelector[TestModel](db).Get(context.Background())
	assert.NotNil(t, err)
}

type TestModel struct {
}
