package orm

import (
	"context"
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/internal/errs"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// From 都不调用
			name: "no from",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM
			name: "with from",
			q:    NewSelector[TestModel](db).From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name: "empty from",
			q:    NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name: "with db",
			q:    NewSelector[TestModel](db).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			// 单一简单条件
			name: "single and simple predicate",
			q:    NewSelector[TestModel](db).From("`test_model_t`").
				Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
					// TestModel -> test_model
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			q:    NewSelector[TestModel](db).
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			q:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
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

func TestSelector_Get(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:"single row",
			query: "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows{
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"))
				return rows
			}(),
			wantVal: &TestModel{
				Id: 123,
				FirstName: "Ming",
				Age: 18,
				LastName: &sql.NullString{Valid: true, String: "Deng"},
			},
		},

		{
			// SELECT 出来的行数小于你结构体的行数
			name:"less columns",
			query: "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows{
				rows := sqlmock.NewRows([]string{"id", "first_name"})
				rows.AddRow([]byte("123"), []byte("Ming"))
				return rows
			}(),
			wantVal: &TestModel{
				Id: 123,
				FirstName: "Ming",
			},
		},

		{
			name:"invalid columns",
			query: "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows{
				rows := sqlmock.NewRows([]string{"id", "first_name", "gender"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("male"))
				return rows
			}(),
			wantErr: errs.NewErrUnknownColumn("gender"),
		},

		{
			name:"more columns",
			query: "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows{
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name",  "first_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"), []byte("明明"))
				return rows
			}(),
			wantErr: errs.ErrTooManyReturnedColumns,
		},
	}

	for _, tc := range testCases {
		if tc.mockErr != nil {
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		} else {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		}
	}


	db, err := OpenDB(mockDB)
	require.NoError(t, err)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVal, res)
		})
	}
}

func memoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 指定列
			name: "specify columns",
			q: NewSelector[TestModel](db).Select(C("Id"), C("Age")),
			wantQuery: &Query{
				SQL: "SELECT `id`,`age` FROM `test_model`;",
			},
		},
		{
			// 指定聚合函数
			// AVG, COUNT, SUM, MIN, MAX(xxx)
			name: "specify aggregate",
			q: NewSelector[TestModel](db).Select(Min("Id"), Avg("Age")),
			wantQuery: &Query{
				SQL: "SELECT MIN(`id`),AVG(`age`) FROM `test_model`;",
			},
		},
		{
			// count distinct
			name: "specify aggregate",
			q: NewSelector[TestModel](db).Select(Raw("DISTINCT `first_name`")),
			wantQuery: &Query{
				SQL: "SELECT DISTINCT `first_name` FROM `test_model`;",
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

// func TestABC(t *testing.T) {
// 	db := memoryDB(t)
// 	tx, err := db.Begin(context.Background(), &sql.TxOptions{})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	NewSelector[TestModel](db)

	// dao := &UserDAO{
	// 	sess: db.db,
	// }
	//
	// txDao := &UserDAO{
	// 	sess: tx.tx,
	// }
	//
	// dao.Config.Name = "abc"
// }

type UserDAO struct {
	sess interface{
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
//
// 	Config struct{
// 		Name string
// 	}
}
//
//  dao.GetByID(context.WithValue("sess", db))
//  dao.GetByID(context.WithValue("sess", tx))
// func (dao *UserDAO) GetByID(ctx context.Context, id uint64) (*User, error){
	// sess := ctx.Value("sess")
	// if sess == nil {
	// 	// 开事务
	// 	// 报错
	// } else {
	// 	// 执行语句
	// }
	// // 处理结果集
// }

// func TestA(t *testing.T) {
// 	ch := Chan{
// 		// bizChan: ///
// 		testChan:
// 	}
//
// 	data <- ch.testChan
// 	// 验证 data
// }

type Chan struct {
	bizChan chan struct{}
	testChan chan struct{}
}

func (c Chan) Put(data struct{}) {
	c.bizChan <- data
	c.testChan <- data
}