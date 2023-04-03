package sql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlMock(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()

	mock.ExpectBegin()
	// mock 返回的行
	mockRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Tom")
	// 或者 WillReturnError
	mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)

	mockResult := sqlmock.NewResult(12, 1)
	// 或者 WillReturnError
	mock.ExpectExec("UPDATE .*").WillReturnResult(mockResult)
	mock.ExpectCommit()

	tx, err := mockDB.Begin()
	assert.Nil(t, err)
	rows, err := tx.QueryContext(context.Background(), "SELECT * FROM `user`")
	cs, err := rows.Columns()
	assert.Nil(t, err)
	assert.Equal(t, []string{"id", "name"}, cs)
	rows.Next()
	var id int
	var name string
	err = rows.Scan(&id, &name)
	assert.Nil(t, err)
	res, err := tx.ExecContext(context.Background(), "UPDATE `user` SET `age` = 12")
	assert.Nil(t, err)
	affected, err := res.RowsAffected()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), affected)
	err = tx.Commit()
	assert.Nil(t, err)
}
