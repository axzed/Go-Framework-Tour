//go:build v7

package orm

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/v7/internal/valuer"
	"gitee.com/geektime-geekbang/geektime-go/orm/v7/model"
)

type DBOption func(*DB)

type DB struct {
	r          model.Registry
	db         *sql.DB
	valCreator valuer.Creator
}

func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:          model.NewRegistry(),
		db:         db,
		valCreator: valuer.NewUnsafeValue,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

// MustNewDB 创建一个 DB，如果失败则会 panic
// 我个人不太喜欢这种
func MustNewDB(driver string, dsn string, opts ...DBOption) *DB {
	db, err := Open(driver, dsn, opts...)
	if err != nil {
		panic(err)
	}
	return db
}
