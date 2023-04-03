//go:build v6

package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModelWithTableName(t *testing.T) {
	testCases := []struct {
		name          string
		val           any
		opt           ModelOpt
		wantTableName string
		wantErr       error
	}{
		{
			// 我们没有对空字符串进行校验
			name:          "empty string",
			val:           &TestModel{},
			opt:           ModelWithTableName(""),
			wantTableName: "",
		},
		{
			name:          "table name",
			val:           &TestModel{},
			opt:           ModelWithTableName("test_model_t"),
			wantTableName: "test_model_t",
		},
	}

	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantTableName, m.tableName)
		})
	}
}

func TestModelWithColumnName(t *testing.T) {
	testCases := []struct {
		name        string
		val         any
		opt         ModelOpt
		field       string
		wantColName string
		wantErr     error
	}{
		{
			name:        "new name",
			val:         &TestModel{},
			opt:         ModelWithColumnName("FirstName", "first_name_new"),
			field:       "FirstName",
			wantColName: "first_name_new",
		},
		{
			name:        "empty new name",
			val:         &TestModel{},
			opt:         ModelWithColumnName("FirstName", ""),
			field:       "FirstName",
			wantColName: "",
		},
		{
			// 不存在的字段
			name:    "invalid field name",
			val:     &TestModel{},
			opt:     ModelWithColumnName("FirstNameXXX", "first_name"),
			field:   "FirstNameXXX",
			wantErr: errs.NewErrUnknownField("FirstNameXXX"),
		},
	}

	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd := m.fieldMap[tc.field]
			assert.Equal(t, tc.wantColName, fd.colName)
		})
	}
}
