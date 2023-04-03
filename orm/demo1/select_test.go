package demo

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// SELECT * FROM xxx;
// SELECT * FROM xxx WHERE id = ?;
func TestSelector_Build(t *testing.T) {
	tests := []struct {
		name    string
		s QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "from",
			s: NewSelector[TestModel]().From("test_model_tab"),
			want: &Query{
				SQL: "SELECT * FROM test_model_tab;",
			},
		},
		{
			name: "no from",
			s: &Selector[TestModel]{},
			want: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "from but empty",
			s: NewSelector[TestModel]().From(""),
			want: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "with db",
			s: NewSelector[TestModel]().From("`test_db`.`test_model`"),
			want: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			name: "single predicate",
			s: NewSelector[TestModel]().Where(C("id").Eq(12)),
			want: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE `id` = ?;",
				Args: []any{12},
			},
		},
		{
			name: "multi predicate",
			s: NewSelector[TestModel]().
				Where(C("Age").GT(18), C("Age").LT(35)),
			want: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "not predicate",
			s: NewSelector[TestModel]().
				Where(Not(C("Age").GT(18))),
			want: &Query{
				SQL: "SELECT * FROM `TestModel` WHERE  NOT (`Age` > ?);",
				Args: []any{18},
			},
		},
	}
	for _, tt := range tests {
		time.NewTimer(0)
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}