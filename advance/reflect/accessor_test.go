package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReflectAccessor_Field(t *testing.T) {
	testCases := []struct {
		name string

		// 这个是输入
		entity interface{}
		field  string

		// 这个是期望输出
		wantVal int
		wantErr error
	}{
		{
			name:    "normal case",
			entity:  &User{Age: 18},
			field:   "Age",
			wantVal: 18,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accessor, err := NewReflectAccessor(tc.entity)
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

func TestReflectAccessor_SetField(t *testing.T) {
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
			accessor, err := NewReflectAccessor(tc.entity)
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
