package orm

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo4/model"
)

type QueryContext struct {
	// 用在 UPDATE，DELETE，SELECT，以及 INSERT 语句上的
	Type string
	Builder QueryBuilder

	Model *model.Model
	TableName string
	DBName string
}

type QueryResult struct {
	// SELECT 语句，你的返回值是 T 或者 []T
	// UPDATE, DELETE, INSERT 返回值是 Result
	Result any

	Err error
}


type Handler func(ctx context.Context, qc *QueryContext) *QueryResult
//
// type HandlerV1 func(qc *QueryContext) *QueryResult
// type HandlerV2 func(qc *QueryContext) (*QueryResult, error)

type Middleware func(next Handler) Handler