//go:build v13

package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/v13/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdater_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		u       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name:    "no columns",
			u:       NewUpdater[TestModel](db),
			wantErr: errs.ErrNoUpdatedColumns,
		},
		{
			name: "single column",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age: 18,
			}).Set(C("Age")),
			want: &Query{
				SQL:  "UPDATE `test_model` SET `age`=?;",
				Args: []any{int8(18)},
			},
		},
		{
			name: "assignment",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(C("Age"), Assign("FirstName", "DaMing")),
			want: &Query{
				SQL:  "UPDATE `test_model` SET `age`=?,`first_name`=?;",
				Args: []any{int8(18), "DaMing"},
			},
		},
		{
			name: "where",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(C("Age"), Assign("FirstName", "DaMing")).
				Where(C("Id").EQ(1)),
			want: &Query{
				SQL:  "UPDATE `test_model` SET `age`=?,`first_name`=? WHERE `id` = ?;",
				Args: []any{int8(18), "DaMing", 1},
			},
		},
		{
			name: "incremental",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(Assign("Age", C("Age").Add(1))),
			want: &Query{
				SQL:  "UPDATE `test_model` SET `age`=`age` + ?;",
				Args: []any{1},
			},
		},
		{
			name: "incremental-raw",
			u: NewUpdater[TestModel](db).Update(&TestModel{
				Age:       18,
				FirstName: "Tom",
			}).Set(Assign("Age", Raw("`age`+?", 1))),
			want: &Query{
				SQL:  "UPDATE `test_model` SET `age`=`age`+?;",
				Args: []any{1},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.u.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, q)
		})
	}
}
