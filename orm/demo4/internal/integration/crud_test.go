//go:build e2e
package integration

import (
	"context"
	orm "gitee.com/geektime-geekbang/geektime-go/orm/demo4"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo4/internal/test"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InsertTestSuite struct {
	suite.Suite
	db *orm.DB

	driver string
	dsn string
}

func (i *InsertTestSuite) SetupSuite() {
	db, err := orm.Open(i.driver, i.dsn)
	if err != nil {
		i.T().Fatal(err)
	}
	i.db = db
	db.Wait()
}

func (i *InsertTestSuite) TestInsert() {
	t := i.T()
	db := i.db

	testCases := []struct{
		name string
		i *orm.Inserter[test.SimpleStruct]

		affected int64
		wantErr error

		wantData *test.SimpleStruct
	}{
		{
			name: "insert single",
			i: orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(15)),
			affected: 1,
			wantData: test.NewSimpleStruct(15),
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
			assert.Equal(t, tc.affected, affected)
			id, err := res.LastInsertId()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			data, err := orm.NewSelector[test.SimpleStruct](db).Where(orm.C("Id").EQ(id)).Get(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.wantData, data)
		})
	}
}

func TestMySQL(t *testing.T) {

	suite.Run(t, &InsertTestSuite{
		driver: "mysql",
		dsn: "root:root@tcp(localhost:13306)/integration_test",
	})
}

//
// func TestSQLite(t *testing.T) {
// 	// 建表语句
// 	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
// 	db.Exec(test.SimpleStruct{})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	suite.Run(t, &InsertTestSuite{
// 		db: db,
// 	})
// }

