package orm

import (
	"context"
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/internal/valuer"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/model"
)

// tx := db.Begin(ctx, ...)
// tx.Commit()
type Tx struct {
	core
	tx *sql.Tx
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) RollbackIfNotCommit() error {
	err := tx.tx.Rollback()
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}

func (tx *Tx) getCore() core {
	return tx.core
}

func  (tx *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func  (tx *Tx) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

type Session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type core struct {
	r model.Registry
	valCreator valuer.Creator
	dialect Dialect
}