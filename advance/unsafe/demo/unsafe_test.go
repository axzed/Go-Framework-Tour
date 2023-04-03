package demo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsafeAccessor_Field(t *testing.T) {
	testCases := []struct {
		name    string
		entity  interface{}
		field   string
		wantVal int
		wantErr error
	}{
		{
			name:    "invalid field",
			entity:  &User{Age: 18},
			field:   "AgeXXXXXX",
			wantErr: errors.New("不存在字段"),
		},
		{
			name:    "normal case",
			entity:  &User{Age: 18},
			field:   "Age",
			wantVal: 18,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor, err := NewUnsafeAccessor(tc.entity)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			val, err := accessor.Field(tc.field)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestUnsafeAccessor_SetField(t *testing.T) {
	testCases := []struct {
		name    string
		entity  *User
		field   string
		newVal  int
		wantErr error
	}{
		{
			name:   "normal case",
			entity: &User{},
			field:  "Age",
			newVal: 18,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor, err := NewUnsafeAccessor(tc.entity)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			err = accessor.SetField(tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.newVal, tc.entity.Age)
		})
	}
}

type User struct {
	Age int
}
