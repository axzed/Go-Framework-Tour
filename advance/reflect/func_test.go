package reflect

import (
	"gitee.com/geektime-geekbang/geektime-go/advance/reflect/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFuncs(t *testing.T) {
	testCases := []struct {
		name string

		input any

		wantRes map[string]*FuncInfo
		wantErr error
	}{
		{
			// 普通结构体
			name:  "normal struct",
			input: types.User{},
			wantRes: map[string]*FuncInfo{
				"GetAge": {
					Name:   "GetAge",
					In:     []reflect.Type{reflect.TypeOf(types.User{})},
					Out:    []reflect.Type{reflect.TypeOf(0)},
					Result: []any{0},
				},
			},
		},
		{
			// 指针
			name:  "pointer",
			input: &types.User{},
			wantRes: map[string]*FuncInfo{
				"GetAge": {
					Name:   "GetAge",
					In:     []reflect.Type{reflect.TypeOf(&types.User{})},
					Out:    []reflect.Type{reflect.TypeOf(0)},
					Result: []any{0},
				},
				"ChangeName": {
					Name:   "ChangeName",
					In:     []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf("")},
					Out:    []reflect.Type{},
					Result: []any{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFuncs(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
