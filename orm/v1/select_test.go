//go:build v1
package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	testCases := []struct{
		name     string
		q        QueryBuilder
		wantQuery *Query
		wantErr  error
	} {
		{
			// From 都不调用
			name:    "no from",
			q:       NewSelector[TestModel](),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			// 调用 FROM
			name:    "with from",
			q:       NewSelector[TestModel]().From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			// 调用 FROM，但是传入空字符串
			name:    "empty from",
			q:       NewSelector[TestModel]().From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			// 调用 FROM，同时出入看了 DB
			name:    "with db",
			q:       NewSelector[TestModel]().From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
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