package demo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFields(t *testing.T) {

	u1 := &User{
		Name: "大明",
	}
	u2 := &u1

	tests := []struct {
		// 名字
		name string

		// 输入部分
		val any

		// 输出部分
		wantRes map[string]any
		wantErr error
	}{
		{
			name:    "nil",
			val:     nil,
			wantErr: errors.New("不能为 nil"),
		},
		{
			name:    "user",
			val:     User{Name: "Tom"},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
			},
		},
		{
			// 指针
			name: "pointer",
			val:  &User{Name: "Jerry"},
			// 要支持指针
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Jerry",
			},
		},
		{
			// 多重指针
			name: "multiple pointer",
			val:  u2,
			// 要支持指针
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "大明",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := iterateFields(tt.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantRes, res)
		})
	}
}

type User struct {
	Name string
}
