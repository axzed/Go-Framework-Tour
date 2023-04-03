//go:build v16
package querylog

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(sql string, args []any)
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args []any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(sql string, args []any) {
			log.Println(sql, args)
		},
	}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			q, err := qc.builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			m.logFunc(q.SQL, q.Args)
			return next(ctx, qc)
		}
	}
}