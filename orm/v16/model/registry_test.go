//go:build v16
package model

import (
	"database/sql"
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestModelWithTableName(t *testing.T) {
	testCases := []struct{
		name string
		val           any
		opt           Option
		wantTableName string
		wantErr error
	} {
		{
			// 我们没有对空字符串进行校验
			name:          "empty string",
			val:           &TestModel{},
			opt:           WithTableName(""),
			wantTableName: "",
		},
		{
			name:          "table name",
			val:           &TestModel{},
			opt:           WithTableName("test_model_t"),
			wantTableName: "test_model_t",
		},
	}

	r := NewRegistry().(*registry)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantTableName, m.TableName)
		})
	}
}

func TestWithColumnName(t *testing.T) {
	testCases := []struct{
		name string
		val   any
		opt   Option
		field string
		wantColName string
		wantErr error
	} {
		{
			name:        "new name",
			val:         &TestModel{},
			opt:         WithColumnName("FirstName", "first_name_new"),
			field:       "FirstName",
			wantColName: "first_name_new",
		},
		{
			name:        "empty new name",
			val:         &TestModel{},
			opt:         WithColumnName("FirstName", ""),
			field:       "FirstName",
			wantColName: "",
		},
		{
			// 不存在的字段
			name:    "invalid Field name",
			val:     &TestModel{},
			opt:     WithColumnName("FirstNameXXX", "first_name"),
			field:   "FirstNameXXX",
			wantErr: errs.NewErrUnknownField("FirstNameXXX"),
		},
	}

	r := NewRegistry().(*registry)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd := m.FieldMap[tc.field]
			assert.Equal(t, tc.wantColName, fd.ColName)
		})
	}
}

func TestRegistry_get(t *testing.T) {
	var tm TestModel
	testCases := []struct{
		name string
		val any
		wantModel *Model
		wantErr error
	} {
		{
			name:    "test Model",
			val:     TestModel{},
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			// 指针
			name: "pointer",
			val: &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				Fields: []*Field {tm.IdField(), tm.FirstNameField(), tm.AgeField(), tm.LastNameField()},
				FieldMap: map[string]*Field{
					"Id": tm.IdField(),
					"FirstName" : tm.FirstNameField(),
					"Age": tm.AgeField(),
					"LastName": tm.LastNameField(),
				},
				ColumnMap: map[string]*Field{
					"id": tm.IdField(),
					"first_name" :tm.FirstNameField(),
					"age": tm.AgeField(),
					"last_name": tm.LastNameField(),
				},
			},
		},
		{
			// 多级指针
			name: "multiple pointer",
			// 因为 Go 编译器的原因，所以我们写成这样
			val: func() any {
				val := &TestModel{}
				return &val
			}(),
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name: "map",
			val: map[string]string{},
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name: "slice",
			val: []int{},
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},
		{
			name: "basic type",
			val: 0,
			wantErr: errors.New("orm: 只支持一级指针作为输入，例如 *User"),
		},

		// 标签相关测试用例
		{
			name: "column tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `orm:"column=id"`
				}
				return &ColumnTag{}
			}(),
			wantModel: &Model{
				TableName: "column_tag",
				Fields: []*Field{{
					ColName: "id",
					Type: reflect.TypeOf(uint64(0)),
					GoName: "ID",
				}},
				FieldMap: map[string]*Field{
					"ID": {
						ColName: "id",
						Type: reflect.TypeOf(uint64(0)),
						GoName: "ID",
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Type: reflect.TypeOf(uint64(0)),
						GoName: "ID",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是传入一个空字符串，那么会用默认的名字
			name: "empty column",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type EmptyColumn struct {
					FirstName string `orm:"column="`
				}
				return &EmptyColumn{}
			}(),
			wantModel: &Model{
				TableName: "empty_column",
				Fields: []*Field{{
					ColName: "first_name",
					Type: reflect.TypeOf(""),
					GoName: "FirstName",
				}},
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						Type: reflect.TypeOf(""),
						GoName: "FirstName",
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						Type: reflect.TypeOf(""),
						GoName: "FirstName",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是没有赋值
			name: "invalid tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type InvalidTag struct {
					FirstName string `orm:"column"`
				}
				return &InvalidTag{}
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
		{
			// 如果用户设置了一些奇奇怪怪的内容，这部分内容我们会忽略掉
			name: "ignore tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type IgnoreTag struct {
					FirstName string `orm:"abc=abc"`
				}
				return &IgnoreTag{}
			}(),
			wantModel: &Model{
				TableName: "ignore_tag",
				Fields: []*Field{{
					ColName: "first_name",
					Type: reflect.TypeOf(""),
					GoName: "FirstName",
				}},
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						Type: reflect.TypeOf(""),
						GoName: "FirstName",
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						Type: reflect.TypeOf(""),
						GoName: "FirstName",
					},
				},
			},
		},

		// 利用接口自定义模型信息
		{
			name: "table name",
			val: &CustomTableName{},
			wantModel: &Model{
				TableName: "custom_table_name_t",
				Fields: []*Field{{
					ColName: "name",
					GoName: "Name",
					Type: reflect.TypeOf(""),
				}},
				FieldMap: map[string]*Field{
					"Name": {
						ColName: "name",
						GoName: "Name",
						Type: reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"name": {
						ColName: "name",
						GoName: "Name",
						Type: reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name: "table name ptr",
			val: &CustomTableNamePtr{},
			wantModel: &Model{
				TableName: "custom_table_name_ptr_t",
				Fields: []*Field{{
					ColName: "name",
					GoName: "Name",
					Type: reflect.TypeOf(""),
				}},
				FieldMap: map[string]*Field{
					"Name": {
						ColName: "name",
						GoName: "Name",
						Type: reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"name": {
						ColName: "name",
						GoName: "Name",
						Type: reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name: "empty table name",
			val: &EmptyTableName{},
			wantModel: &Model{
				TableName: "empty_table_name",
				Fields: []*Field{{
					ColName: "name",
					GoName: "Name",
					Type: reflect.TypeOf(""),
				}},
				FieldMap: map[string]*Field{
					"Name": {
						ColName: "name",
						GoName: "Name",
						Type: reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"name": {
						ColName: "name",
						GoName: "Name",
						Type: reflect.TypeOf(""),
					},
				},
			},
		},
	}

	r := &registry{}
	for _, tc :=range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func Test_underscoreName(t *testing.T) {
	testCases := []struct {
		name    string
		srcStr  string
		wantStr string
	}{
		// 我们这些用例就是为了确保
		// 在忘记 underscoreName 的行为特性之后
		// 可以从这里找回来
		// 比如说过了一段时间之后
		// 忘记了 ID 不能转化为 id
		// 那么这个测试能帮我们确定 ID 只能转化为 i_d
		{
			name:    "upper cases",
			srcStr:  "ID",
			wantStr: "i_d",
		},
		{
			name:    "use number",
			srcStr:  "Table1Name",
			wantStr: "table1_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := underscoreName(tc.srcStr)
			assert.Equal(t, tc.wantStr, res)
		})
	}
}

type CustomTableName struct {
	Name string
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

type CustomTableNamePtr struct {
	Name string
}

func (c *CustomTableNamePtr) TableName() string {
	return "custom_table_name_ptr_t"
}

type EmptyTableName struct {
	Name string
}

func (c *EmptyTableName) TableName() string {
	return ""
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func (TestModel) IdField() *Field {
	return &Field{
		ColName: "id",
		Type: reflect.TypeOf(int64(0)),
		GoName: "Id",
		Offset: 0,
		Index: 0,
	}
}

func (TestModel) FirstNameField() *Field {
	return &Field{
		ColName: "first_name",
		Type: reflect.TypeOf(""),
		GoName: "FirstName",
		Offset: 8,
		Index: 1,
	}
}

func (TestModel) AgeField() *Field {
	return &Field {
		ColName: "age",
		Type: reflect.TypeOf(int8(0)),
		GoName: "Age",
		Offset: 24,
		Index: 2,
	}
}

func (TestModel) LastNameField() *Field {
	return &Field{
		ColName: "last_name",
		Type: reflect.TypeOf(&sql.NullString{}),
		GoName: "LastName",
		Offset: 32,
		Index: 3,
	}
}