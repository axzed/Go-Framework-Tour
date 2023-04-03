//go:build e2e
package integration

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm"
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InsertTestSuite struct {
	Suite
}

func (i *InsertTestSuite) TearDownTest() {
	res := orm.RawQuery[any](i.db, "TRUNCATE TABLE `simple_struct`").
		Exec(context.Background())
	require.NoError(i.T(), res.Err())
}

func TestInsertMySQL8(t *testing.T) {
	suite.Run(t, &InsertTestSuite{
		Suite: Suite {
			driver: "mysql",
			dsn: "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}


// func testInsert(t *testing.T, driver string, dsn string) {
// 	db, err := orm.Open(driver, dsn)
// 	defer func() {
// 		orm.RawQuery[any](db, "TRUNCATE TABLE `simple_struct`").Exec(context.Background())
// 	}()
// 	require.NoError(t, err)
// 	testCases := []struct{
// 		name string
// 		i *orm.Inserter[test.SimpleStruct]
// 		wantData *test.SimpleStruct
// 		rowsAffected int64
// 		wantErr error
// 	} {
// 		{
// 			name: "id only",
// 			i: orm.NewInserter[test.SimpleStruct](db).Values(&test.SimpleStruct{Id: 1}),
// 			rowsAffected: 1,
// 			wantData: &test.SimpleStruct{ Id: 1},
// 		},
// 		{
// 			name: "all field",
// 			i: orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(2)),
// 			rowsAffected: 1,
// 			wantData: test.NewSimpleStruct(2),
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			res := tc.i.Exec(context.Background())
// 			affected, err := res.RowsAffected()
// 			assert.Equal(t, tc.wantErr, err)
// 			if err != nil {
// 				return
// 			}
// 			assert.Equal(t, tc.rowsAffected, affected)
// 			data, err := orm.NewSelector[test.SimpleStruct](db).
// 				Where(orm.C("Id").EQ(tc.wantData.Id)).
// 				Get(context.Background())
// 			require.NoError(t, err)
// 			assert.Equal(t, tc.wantData, data)
// 		})
// 	}
// }
//
// func TestInsert(t *testing.T) {
// 	testInsert(t, "mysql", "root:root@tcp(localhost:13306)/integration_test")
// }

func (i *InsertTestSuite) TestInsert() {
	t := i.T()
	db, err := orm.Open("mysql", "root:root@tcp(localhost:13306)/integration_test")
	defer func() {
		orm.RawQuery[any](db, "TRUNCATE TABLE `simple_struct`").Exec(context.Background())
	}()
	require.NoError(t, err)
	testCases := []struct{
		name string
		i *orm.Inserter[test.SimpleStruct]
		wantData *test.SimpleStruct
		rowsAffected int64
		wantErr error
	} {
		{
			name: "id only",
			i: orm.NewInserter[test.SimpleStruct](db).Values(&test.SimpleStruct{Id: 1}),
			rowsAffected: 1,
			wantData: &test.SimpleStruct{ Id: 1},
		},
		{
			name: "all field",
			i: orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(2)),
			rowsAffected: 1,
			wantData: test.NewSimpleStruct(2),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			affected, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.rowsAffected, affected)
			data, err := orm.NewSelector[test.SimpleStruct](db).
				Where(orm.C("Id").EQ(tc.wantData.Id)).
				Get(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.wantData, data)
		})
	}
}
