//go:build v14

package valuer

import (
	"database/sql/driver"
	"gitee.com/geektime-geekbang/geektime-go/orm/v14/internal/errs"
	"gitee.com/geektime-geekbang/geektime-go/orm/v14/internal/test"
	"gitee.com/geektime-geekbang/geektime-go/orm/v14/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_unsafeValue_SetColumn(t *testing.T) {
	testCases := []struct {
		name    string
		cs      map[string][]byte
		val     *test.SimpleStruct
		wantVal *test.SimpleStruct
		wantErr error
	}{
		{
			name: "normal value",
			cs: map[string][]byte{
				"id":               []byte("1"),
				"bool":             []byte("true"),
				"bool_ptr":         []byte("false"),
				"int":              []byte("12"),
				"int_ptr":          []byte("13"),
				"int8":             []byte("8"),
				"int8_ptr":         []byte("-8"),
				"int16":            []byte("16"),
				"int16_ptr":        []byte("-16"),
				"int32":            []byte("32"),
				"int32_ptr":        []byte("-32"),
				"int64":            []byte("64"),
				"int64_ptr":        []byte("-64"),
				"uint":             []byte("14"),
				"uint_ptr":         []byte("15"),
				"uint8":            []byte("8"),
				"uint8_ptr":        []byte("18"),
				"uint16":           []byte("16"),
				"uint16_ptr":       []byte("116"),
				"uint32":           []byte("32"),
				"uint32_ptr":       []byte("132"),
				"uint64":           []byte("64"),
				"uint64_ptr":       []byte("164"),
				"float32":          []byte("3.2"),
				"float32_ptr":      []byte("-3.2"),
				"float64":          []byte("6.4"),
				"float64_ptr":      []byte("-6.4"),
				"byte":             []byte("8"),
				"byte_ptr":         []byte("18"),
				"byte_array":       []byte("hello"),
				"string":           []byte("world"),
				"null_string_ptr":  []byte("null string"),
				"null_int16_ptr":   []byte("16"),
				"null_int32_ptr":   []byte("32"),
				"null_int64_ptr":   []byte("64"),
				"null_bool_ptr":    []byte("true"),
				"null_float64_ptr": []byte("6.4"),
				"json_column":      []byte(`{"name": "Tom"}`),
			},
			val:     &test.SimpleStruct{},
			wantVal: test.NewSimpleStruct(1),
		},
		{
			name: "invalid field",
			cs: map[string][]byte{
				"invalid_column": nil,
			},
			wantErr: errs.NewErrUnknownColumn("invalid_column"),
		},
	}
	r := model.NewRegistry()
	meta, err := r.Get(&test.SimpleStruct{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = db.Close() }()
			val := NewUnsafeValue(tc.val, meta)
			cols := make([]string, 0, len(tc.cs))
			colVals := make([]driver.Value, 0, len(tc.cs))
			for k, v := range tc.cs {
				cols = append(cols, k)
				colVals = append(colVals, v)
			}
			mock.ExpectQuery("SELECT *").
				WillReturnRows(sqlmock.NewRows(cols).
					AddRow(colVals...))
			rows, _ := db.Query("SELECT *")
			rows.Next()
			err = val.SetColumns(rows)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			if tc.wantErr != nil {
				t.Fatalf("期望得到错误，但是并没有得到 %v", tc.wantErr)
			}
			assert.Equal(t, tc.wantVal, tc.val)
		})
	}

}
