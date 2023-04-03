package orm

import (
	"context"
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo4/internal/valuer"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo4/model"
	"go.uber.org/multierr"
)

type DBOption func(*DB)

// DB 是sql.DB 的装饰器
type DB struct {
	core
	db *sql.DB
	dialect Dialect
}

func (db *DB) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
	}, nil
}
func (db *DB) DoTx(ctx context.Context, opts *sql.TxOptions,
	task func(ctx context.Context, tx *Tx) error ) (err error) {
	tx, err := db.Begin(ctx, opts)
	if err != nil {
		return err
	}
	panicked := true
	defer func() {
		if panicked || err != nil {
			er := tx.Rollback()
			err = multierr.Combine(err, er)
		} else {
			err = multierr.Combine(err, tx.Commit())
		}
	}()

	err =  task(ctx, tx)
	panicked = false
	return
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
		core: core{
			r:          model.NewRegistry(),
			valCreator: valuer.NewUnsafeValue,
			dialect: &mysqlDialect{},
		},
		db:         db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBName(name string) DBOption {
	return func(db *DB) {
		db.dbName = name
	}
}

func DBWithMiddlewares(ms...Middleware) DBOption {
	return func(db *DB) {
		db.ms = ms
	}
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

func (db *DB) getCore() core {
	return db.core
}

func  (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func  (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// Wait 用于测试等待容器启动成功
func (db *DB) Wait() {
	for db.db.PingContext(context.Background()) != nil {

	}
}