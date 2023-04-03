package sql

import (
	"context"
	"database/sql"
	"time"
)

func (s *sqlTestSuite) TestTx() {
	t := s.T()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	// 这种依旧没有在事务里面
	// s.db.Exec()
	res, err := tx.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (2, 'Tom', 20, 'Jerry')")
	if err != nil {
		t.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}
	// 回滚则是 tx.Rollback
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
