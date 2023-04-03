//go:build v13

package orm

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/v13/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 一个都不插入
			name:    "no value",
			q:       NewInserter[TestModel](db).Values(),
			wantErr: errs.ErrInsertZeroRow,
		},
		{
			name: "single values",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
		{
			name: "multiple values",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				},
				&TestModel{
					Id:        2,
					FirstName: "Da",
					Age:       19,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?);",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true},
					int64(2), "Da", int8(19), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
		{
			// 指定列
			name: "specify columns",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).Columns("FirstName", "LastName"),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`first_name`,`last_name`) VALUES(?,?);",
				Args: []any{"Deng", &sql.NullString{String: "Ming", Valid: true}},
			},
		},
		{
			// 指定列
			name: "invalid columns",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).Columns("FirstName", "Invalid"),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},

		{
			// upsert
			name: "upsert",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).OnDuplicateKey().Update(Assign("FirstName", "Da")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=?;",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true}, "Da"},
			},
		},
		{
			// upsert invalid column
			name: "upsert invalid column",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).OnDuplicateKey().Update(Assign("Invalid", "Da")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			// 使用原本插入的值
			name: "upsert use insert value",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				},
				&TestModel{
					Id:        2,
					FirstName: "Da",
					Age:       19,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).OnDuplicateKey().Update(C("FirstName"), C("LastName")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=VALUES(`first_name`),`last_name`=VALUES(`last_name`);",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true},
					int64(2), "Da", int8(19), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestUpsert_SQLite3_Build(t *testing.T) {
	db := memoryDB(t, DBWithDialect(SQLite3))
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// upsert
			name: "upsert",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Id").
				Update(Assign("FirstName", "Da")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?) " +
					"ON CONFLICT(`id`) DO UPDATE SET `first_name`=?;",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true}, "Da"},
			},
		},
		{
			// upsert invalid column
			name: "upsert invalid column",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Id").
				Update(Assign("Invalid", "Da")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			// conflict invalid column
			name: "conflict invalid column",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Invalid").
				Update(Assign("FirstName", "Da")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			// 使用原本插入的值
			name: "upsert use insert value",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "Deng",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				},
				&TestModel{
					Id:        2,
					FirstName: "Da",
					Age:       19,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Id").
				Update(C("FirstName"), C("LastName")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?) " +
					"ON CONFLICT(`id`) DO UPDATE SET `first_name`=excluded.`first_name`,`last_name`=excluded.`last_name`;",
				Args: []any{int64(1), "Deng", int8(18), &sql.NullString{String: "Ming", Valid: true},
					int64(2), "Da", int8(19), &sql.NullString{String: "Ming", Valid: true}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
