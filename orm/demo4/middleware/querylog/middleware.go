package querylog

import (
	"context"
	orm "gitee.com/geektime-geekbang/geektime-go/orm/demo4"
	"time"
)

type MiddlewareBuilder struct {
	// 慢查询的阈值，毫秒单位
	threshold int64
	logFunc func(sql string, args...any)
}

func (m *MiddlewareBuilder) SlowQueryThreshold(threshold int64) *MiddlewareBuilder {
	m.threshold = threshold
	return m
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args...any) ) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			start := time.Now()
			q, err := qc.Builder.Build()
			if err != nil {
				// 构造 SQL 失败了
				return &orm.QueryResult{
					Err: err,
				}
			}
			defer func() {
				duration := time.Now().Sub(start)
				// 设置了慢查询阈值，并且触发了
				// 我想知道是哪个数据库
				if m.threshold >0 && duration.Milliseconds() > m.threshold {
					m.logFunc(q.SQL, q.Args...)
				}
			}()
			return next(ctx, qc)
		}
	}
}
