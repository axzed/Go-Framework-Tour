//go:build v16
package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/valuer"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Join(t *testing.T) {
	db := memoryDB(t)

	type Order struct {
		Id int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId int

		UsingCol1 string
		UsingCol2 string
	}

	type Item struct {
		Id int
	}

	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 虽然泛型是 Order，但是我们传入 OrderDetail
			name: "specify table",
			q: NewSelector[Order](db).From(TableOf(&OrderDetail{})),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order_detail`;",
			},
		},
		{
			name: "join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{})
				return NewSelector[Order](db).
					From(t1.Join(t2).On(t1.C("Id").EQ(t2.C("OrderId"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` ON `t1`.`id` = `order_id`);",
			},
		},
		{
			name: "multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelector[Order](db).
					From(t1.Join(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						Join(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((`order` AS `t1` JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},
		{
			name: "left multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelector[Order](db).
					From(t1.LeftJoin(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
						LeftJoin(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((`order` AS `t1` LEFT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) LEFT JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},
		{
			name: "right multiple join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{}).As("t2")
				t3 := TableOf(&Item{}).As("t3")
				return NewSelector[Order](db).
					From(t1.RightJoin(t2).
						On(t1.C("Id").EQ(t2.C("OrderId"))).
					RightJoin(t3).On(t2.C("ItemId").EQ(t3.C("Id"))))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((`order` AS `t1` RIGHT JOIN `order_detail` AS `t2` ON `t1`.`id` = `t2`.`order_id`) RIGHT JOIN `item` AS `t3` ON `t2`.`item_id` = `t3`.`id`);",
			},
		},

		{
			name: "join multiple using",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{}).As("t1")
				t2 := TableOf(&OrderDetail{})
				return NewSelector[Order](db).
					From(t1.Join(t2).Using("UsingCol1", "UsingCol2"))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (`order` AS `t1` JOIN `order_detail` USING (`using_col1`,`using_col2`));",
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

func TestSelector_Subquery(t *testing.T) {
	db := memoryDB(t)
	type Order struct {
		Id int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId int
	}

	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "from",
			q: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).AsSubquery("sub")
				return NewSelector[Order](db).From(sub)
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM (SELECT * FROM `order_detail`) AS `sub`;",
			},
		},
		{
			name: "in",
			q: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
				return NewSelector[Order](db).Where(C("Id").InQuery(sub))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE `id` IN (SELECT `order_id` FROM `order_detail`);",
			},
		},
		{
			name: "exist",
			q: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
				return NewSelector[Order](db).Where(Exist(sub))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE  EXIST (SELECT `order_id` FROM `order_detail`);",
			},
		},
		{
			name: "not exist",
			q: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
				return NewSelector[Order](db).Where(Not(Exist(sub)))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE  NOT ( EXIST (SELECT `order_id` FROM `order_detail`));",
			},
		},
		{
			name: "all",
			q: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
				return NewSelector[Order](db).Where(C("Id").GT(All(sub)))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE `id` > ALL (SELECT `order_id` FROM `order_detail`);",
			},
		},
		{
			name: "some and any",
			q: func() QueryBuilder {
				sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
				return NewSelector[Order](db).Where(C("Id").GT(Some(sub)), C("Id").LT(Any(sub)))
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `order` WHERE (`id` > SOME (SELECT `order_id` FROM `order_detail`)) AND (`id` < ANY (SELECT `order_id` FROM `order_detail`;SELECT `order_id` FROM `order_detail`));",
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
// Join 和 Subquery 混合使用
func TestSelector_SubqueryAndJoin(t *testing.T) {
	db := memoryDB(t)

	type Order struct {
		Id int
		UsingCol1 string
		UsingCol2 string
	}

	type OrderDetail struct {
		OrderId int
		ItemId int

		UsingCol1 string
		UsingCol2 string
	}

	type Item struct {
		Id int
	}

	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 虽然泛型是 Order，但是我们传入 OrderDetail
			name: "table and join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubquery("sub")
				return NewSelector[Order](db).Select(sub.C("ItemId")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId")))).Where()
			}(),
			wantQuery: &Query{
				SQL: "SELECT `sub`.`item_id` FROM (`order` JOIN (SELECT * FROM `order_detail`) AS `sub` ON `id` = `sub`.`order_id`);",
			},
		},
		{
			name: "table and left join",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubquery("sub")
				return NewSelector[Order](db).From(sub.Join(t1).On(t1.C("Id").EQ(sub.C("OrderId")))).Where()
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((SELECT * FROM `order_detail`) AS `sub` JOIN `order` ON `id` = `sub`.`order_id`);",
			},
		},
		{
			name: "join and join",
			q: func() QueryBuilder {
				sub1 := NewSelector[OrderDetail](db).AsSubquery("sub1")
				sub2 := NewSelector[OrderDetail](db).AsSubquery("sub2")
				return NewSelector[Order](db).From(sub1.RightJoin(sub2).Using("Id")).Where()
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((SELECT * FROM `order_detail`) AS `sub1` RIGHT JOIN (SELECT * FROM `order_detail`) AS `sub2` USING (`id`));",
			},
		},
		{
			name: "join sub sub",
			q: func() QueryBuilder {
				sub1 := NewSelector[OrderDetail](db).AsSubquery("sub1")
				sub2 := NewSelector[OrderDetail](db).From(sub1).AsSubquery("sub2")
				t1 := TableOf(&Order{}).As("o1")
				return NewSelector[Order](db).From(sub2.Join(t1).Using("Id")).Where()
			}(),
			wantQuery: &Query{
				SQL: "SELECT * FROM ((SELECT * FROM (SELECT * FROM `order_detail`) AS `sub1`) AS `sub2` JOIN `order` AS `o1` USING (`id`));",
			},
		},
		{
			name: "invalid field",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubquery("sub")
				return NewSelector[Order](db).Select(sub.C("Invalid")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId")))).Where()
			}(),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "invalid field in predicates",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubquery("sub")
				return NewSelector[Order](db).Select(sub.C("ItemId")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("Invalid")))).Where()
			}(),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "invalid field in aggregate function",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).AsSubquery("sub")
				return NewSelector[Order](db).Select(Max("Invalid")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId")))).Where()
			}(),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "not selected",
			q: func() QueryBuilder {
				t1 := TableOf(&Order{})
				sub := NewSelector[OrderDetail](db).Select(C("OrderId")).AsSubquery("sub")
				return NewSelector[Order](db).Select(sub.C("ItemId")).From(t1.Join(sub).On(t1.C("Id").EQ(sub.C("OrderId")))).Where()
			}(),
			wantErr: errs.NewErrUnknownField("ItemId"),
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

func TestSelector_OffsetLimit(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "offset only",
			q:    NewSelector[TestModel](db).Offset(10),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` OFFSET ?;",
				Args: []any{10},
			},
		},
		{
			name: "limit only",
			q:    NewSelector[TestModel](db).Limit(10),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` LIMIT ?;",
				Args: []any{10},
			},
		},
		{
			name: "limit offset",
			q:    NewSelector[TestModel](db).Limit(20).Offset(10),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` LIMIT ? OFFSET ?;",
				Args: []any{20, 10},
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

func TestSelector_Having(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 调用了，但是啥也没传
			name: "none",
			q:    NewSelector[TestModel](db).GroupBy(C("Age")).Having(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`;",
			},
		},
		{
			// 单个条件
			name: "single",
			q:    NewSelector[TestModel](db).GroupBy(C("Age")).
				Having(C("FirstName").EQ("Deng")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age` HAVING `first_name` = ?;",
				Args: []any{"Deng"},
			},
		},
		{
			// 多个条件
			name: "multiple",
			q:    NewSelector[TestModel](db).GroupBy(C("Age")).
				Having(C("FirstName").EQ("Deng"), C("LastName").EQ("Ming")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age` HAVING (`first_name` = ?) AND (`last_name` = ?);",
				Args: []any{"Deng", "Ming"},
			},
		},
		{
			// 聚合函数
			name: "avg",
			q:    NewSelector[TestModel](db).GroupBy(C("Age")).
				Having(Avg("Age").EQ(18)),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age` HAVING AVG(`age`) = ?;",
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

func TestSelector_GroupBy(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 调用了，但是啥也没传
			name: "none",
			q:    NewSelector[TestModel](db).GroupBy(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			// 单个
			name: "single",
			q:    NewSelector[TestModel](db).GroupBy(C("Age")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`;",
			},
		},
		{
			// 多个
			name: "multiple",
			q:    NewSelector[TestModel](db).GroupBy(C("Age"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` GROUP BY `age`,`first_name`;",
			},
		},
		{
			// 不存在
			name: "invalid column",
			q:    NewSelector[TestModel](db).GroupBy(C("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
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

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// 没有指定
			name: "all",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name: "invalid column",
			q:    NewSelector[TestModel](db).Select(Avg("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "partial columns",
			q:    NewSelector[TestModel](db).Select(C("Id"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT `id`,`first_name` FROM `test_model`;",
			},
		},
		{
			name: "avg",
			q:    NewSelector[TestModel](db).Select(Avg("Age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "raw expression",
			q:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},
		// 别名
		{
			name: "alias",
			q:    NewSelector[TestModel](db).
				Select(C("Id").As("my_id"),
					Avg("Age").As("avg_age")),
			wantQuery: &Query{
				SQL: "SELECT `id` AS `my_id`,AVG(`age`) AS `avg_age` FROM `test_model`;",
			},
		},
		// WHERE 忽略别名
		{
			name: "where ignore alias",
			q:    NewSelector[TestModel](db).
				Where(C("Id").As("my_id").LT(100)),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE `id` < ?;",
				Args: []any{100},
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
			// 单一简单条件
			name: "single and simple predicate",
			q:    NewSelector[TestModel](db).Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			q: NewSelector[TestModel](db).
				Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
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
		{
			// 非法列
			name: "invalid column",
			q:    NewSelector[TestModel](db).Where(Not(C("Invalid").GT(18))),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			// 使用 RawExpr
			name: "raw expression",
			q:    NewSelector[TestModel](db).
				Where(Raw("`age` < ?", 18).AsPredicate()),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age` < ?;",
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
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := OpenDB(mockDB)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			// 查询返回错误
			name:    "query error",
			mockErr: errors.New("invalid query"),
			wantErr: errors.New("invalid query"),
			query:   "SELECT .*",
		},
		{
			name:     "no row",
			wantErr:  ErrNoRows,
			query:    "SELECT .*",
			mockRows: sqlmock.NewRows([]string{"id"}),
		},
		{
			name:    "too many column",
			wantErr: errs.ErrTooManyReturnedColumns,
			query:   "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name", "extra_column"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"), []byte("nothing"))
				return res
			}(),
		},
		{
			name:  "get data",
			query: "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"))
				return res
			}(),
			wantVal: &TestModel{
				Id:        1,
				FirstName: "Da",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			},
		},
	}

	for _, tc := range testCases {
		exp := mock.ExpectQuery(tc.query)
		if tc.mockErr != nil {
			exp.WillReturnError(tc.mockErr)
		} else {
			exp.WillReturnRows(tc.mockRows)
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}

// 在 orm 目录下执行
// go test -bench=BenchmarkQuerier_Get -benchmem -benchtime=10000x
// 我的输出结果
// goos: linux
// goarch: amd64
// pkg: gitee.com/geektime-geekbang/geektime-go/orm
// cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
// BenchmarkQuerier_Get/unsafe-12             10000            453677 ns/op            3246 B/op        108 allocs/op
// BenchmarkQuerier_Get/reflect-12            10000           1173199 ns/op            3427 B/op        117 allocs/op
// PASS
// ok      gitee.com/geektime-geekbang/geektime-go/orm     16.324s
func BenchmarkQuerier_Get(b *testing.B) {
	db, err := Open("sqlite3", fmt.Sprintf("file:benchmark_get.db?cache=shared&mode=memory"))
	if err != nil {
		b.Fatal(err)
	}
	_, err = db.db.Exec(TestModel{}.CreateSQL())
	if err != nil {
		b.Fatal(err)
	}

	res, err := db.db.Exec("INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`)" +
		"VALUES (?,?,?,?)", 12, "Deng", 18, "Ming")

	if err != nil {
		b.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		b.Fatal(err)
	}
	if affected == 0 {
		b.Fatal()
	}

	b.Run("unsafe", func(b *testing.B) {
		db.valCreator = valuer.NewUnsafeValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("reflect", func(b *testing.B) {
		db.valCreator = valuer.NewReflectValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}