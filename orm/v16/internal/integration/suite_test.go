//go:build e2e
package integration

import (
	"gitee.com/geektime-geekbang/geektime-go/orm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	driver string
	dsn string

	db *orm.DB
}

func (i *Suite) SetupSuite() {
	db, err := orm.Open(i.driver, i.dsn)
	require.NoError(i.T(), err)
	err = db.Wait()
	require.NoError(i.T(), err)
	i.db = db
}