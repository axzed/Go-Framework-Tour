package orm

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo2/internal/valuer"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo2/model"
)

type DBOption func(*DB)

// DB 是sql.DB 的装饰器
type DB struct {
	db *sql.DB
	r  model.Registry
	dialect Dialect

	valCreator valuer.Creator
}

// 如果用户指定了 registry，就用用户指定的，否则用默认的

// db := Open()

// r1 := NewRegistry()
// db1 := Open(r1)
// db2 := Open(r1)

func Open(driver string, dsn string, opts...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)

	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

// OpenDB
// 我可以利用 OpenDB 来传入一个 mock 的DB
// sqlmock.Open 的 DB
func OpenDB(db *sql.DB, opts...DBOption) (*DB, error) {
	res := &DB{
		r:          model.NewRegistry(),
		db:         db,
		valCreator: valuer.NewUnsafeValue,
		dialect: &mysqlDialect{},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}


// func MustNewDB(opts...DBOption) *DB{
// 	res, err := Open(opts...)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return res
// }

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

//
// func (db *DB) NewSelector[T any]() *Selector[T] {
// 	return &Selector[T]{
// 		db: db,
// 	}
// }