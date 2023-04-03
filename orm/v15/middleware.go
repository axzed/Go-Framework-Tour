//go:build v15
package orm

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm/model"
)

type QueryContext struct {
	// Type 声明查询类型。即 SELECT, UPDATE, DELETE 和 INSERT
	Type string

	// builder 使用的时候，大多数情况下你需要转换到具体的类型
	// 才能篡改查询
	Builder QueryBuilder
	Model *model.Model
}

type QueryResult struct {
	// Result 在不同的查询里面，类型是不同的
	// Selector.Get 里面，这会是单个结果
	// Selector.GetMulti，这会是一个切片
	// 其它情况下，它会是 Result 类型
	Result any
	Err error
}

type Middleware func(next HandleFunc) HandleFunc

type HandleFunc func(ctx context.Context, qc *QueryContext) *QueryResult


