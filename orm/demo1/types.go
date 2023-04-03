package demo

import (
	"context"
	"database/sql"
)

// SELECT 语句
type Queier[T any] interface {
	// user := xxx.Get(ctx)

	// 不再需要写成
	// var user User
	// Get(ctx, &user)
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// UPDATE, DELETE, INSERT
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type QueryBuilder interface {
	Build() (*Query, error)
}

// db.Exec
// db.QueryRow
// db.Query
type Query struct {
	SQL string
	Args []any
}
