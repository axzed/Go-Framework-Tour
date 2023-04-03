package reflect

import (
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/advance/reflect/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TDD => test driven development
func TestIterateFields(t *testing.T) {
	up := &types.User{}
	up2 := &up
	testCases := []struct {
		name       string
		input      any
		wantFields map[string]any
		wantErr    error
	}{
		{
			// 普通结构体
			name: "normal struct",
			input: types.User{
				Name: "Tom",
				// age:  18,
			},
			wantFields: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			// 指针
			name: "pointer",
			input: &types.User{
				Name: "Tom",
			},
			wantFields: map[string]any{
				"Name": "Tom",
				"age":  0,
			},
		},
		{
			// 多重指针
			name:  "multiple pointer",
			input: up2,
			wantFields: map[string]any{
				"Name": "",
				"age":  0,
			},
		},
		{
			// 非法输入
			name:    "slice",
			input:   []string{},
			wantErr: errors.New("非法类型"),
		},
		{
			// 非法指针输入
			name:    "pointer to map",
			input:   &(map[string]string{}),
			wantErr: errors.New("非法类型"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := iterateFields(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantFields, res)
		})
	}
}

func TestSetField(t *testing.T) {
	testCases := []struct {
		name string

		field  string
		entity any
		newVal any

		wantErr error
	}{
		{
			name:    "struct",
			entity:  types.User{},
			field:   "Name",
			wantErr: errors.New("非法类型"),
		},
		{
			name:    "private field",
			entity:  &types.User{},
			field:   "age",
			wantErr: errors.New("不可修改字段"),
		},
		{
			name:    "invalid field",
			entity:  &types.User{},
			field:   "invalid_field",
			wantErr: errors.New("字段不存在"),
		},
		{
			name: "pass",
			entity: &types.User{
				Name: "",
			},
			field:  "Name",
			newVal: "Tom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
