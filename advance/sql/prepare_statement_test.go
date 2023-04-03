package sql

import (
	"context"
	"github.com/stretchr/testify/assert"
)

func (s *sqlTestSuite) TestPrepareStatement() {
	t := s.T()

	stmt, err := s.db.Prepare("SELECT * FROM `test_model` WHERE `id` = ?")
	if err != nil {
		t.Fatal(err)
	}

	// SELECT * FROM `user` WHERE `id` = 1
	_, err = stmt.QueryContext(context.Background(), 1)
	assert.Nil(t, err)

	// SELECT * FROM `user` WHERE `id` = 1
	_, err = stmt.QueryContext(context.Background(), 2)
	assert.Nil(t, err)

	// 用完就关闭
	err = stmt.Close()
	assert.Nil(t, err)
}
