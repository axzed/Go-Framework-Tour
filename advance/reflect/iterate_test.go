package reflect

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterate(t *testing.T) {
	testCases := []struct {
		name string

		input any
		
		wantRes []any
		wantErr error
	}{
		{
			name:    "slice",
			input:   []int{1, 2, 3},
			wantRes: []any{1, 2, 3},
		},
		{
			name:    "array",
			input:   [5]int{1, 2, 3, 4, 5},
			wantRes: []any{1, 2, 3, 4, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Iterate(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestIterateMap(t *testing.T) {
	testCases := []struct {
		name       string
		input      any
		wantKeys   []any
		wantValues []any
		wantErr    error
	}{
		{
			name:    "nil",
			input:   nil,
			wantErr: errors.New("非法类型"),
		},
		{
			name: "happy case",
			input: map[string]string{
				"a_k": "a_v",
			},
			wantKeys:   []any{"a_k"},
			wantValues: []any{"a_v"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keys, vals, err := IterateMapV1(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantKeys, keys)
			assert.Equal(t, tc.wantValues, vals)

			keys, vals, err = IterateMapV2(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantKeys, keys)
			assert.Equal(t, tc.wantValues, vals)
		})
	}
}

type UserService struct {
	GetByIdV1 func()
}

func (u *UserService) GetByIdV2() {
	fmt.Println("aa")
}
