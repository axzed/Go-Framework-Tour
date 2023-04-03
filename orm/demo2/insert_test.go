package orm

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo2/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct{
		name string
		i QueryBuilder
		wantQuery *Query
		wantErr error
	}{
		{
			// 一个都不插入
			name: "no value",
			i: NewInserter[TestModel](db).Values(),
			wantErr: errs.ErrInsertZeroRow,
		},
		{
			// 插入当个值的全部列，其实就是插入一行
			name: "single value",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id: 12,
				FirstName: "Tom",
				Age: 18,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			}),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}},
			},
		},
		{
			// 插入多行
			name: "multi value",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id: 12,
				FirstName: "Tom",
				Age: 18,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			}, &TestModel{
				Id: 13,
				FirstName: "XiaoMing",
				Age: 17,
				LastName: &sql.NullString{Valid: true, String: "Deng"},
			}),
			wantQuery: &Query{
				// INSERT INTO `test_model` (`id`, `age`) VALUES, `id` + 1)
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?);",
				Args: []any{
					int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"},
					int64(13), "XiaoMing", int8(17), &sql.NullString{Valid: true, String: "Deng"},
				},
			},
		},

		// 指定列
		{
			name: "specify columns",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id: 12,
				FirstName: "Tom",
				Age: 18,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			}).Columns("Age", "FirstName"),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`age`,`first_name`) VALUES(?,?);",
				Args: []any{int8(18), "Tom"},
			},
		},

		{
			name: "upsert",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id: 12,
				FirstName: "Tom",
				Age: 18,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			}).OnDuplicateKey().Update(Assign("Age", 19)),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age`=?;",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}, 19},
			},
		},

		{
			name: "upsert multiple",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id: 12,
				FirstName: "Tom",
				Age: 18,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			}).OnDuplicateKey().Update(Assign("Age", 19), C("FirstName")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age`=?,`first_name`=VALUES(`first_name`);",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}, 19},
			},
		},

		{
			name: "upsert use columns",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id: 12,
				FirstName: "Tom",
				Age: 18,
				LastName: &sql.NullString{Valid: true, String: "Jerry"},
			}).OnDuplicateKey().Update(C("Age")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)" +
					" ON DUPLICATE KEY UPDATE `age`=VALUES(`age`);",
				Args: []any{int64(12), "Tom", int8(18), &sql.NullString{Valid: true, String: "Jerry"}},
			},
		},
	}
	for _, tc :=range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.i.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, q)
		})
	}
}
