//go:build v14
package test

import (
	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonColumn_Scan(t *testing.T) {
	type User struct {
		Name string
	}
	testCases := []struct{
		name string
		input any
		wantVal User
		wantErr string
	} {
		{
			name: "empty string",
			input: ``,
		},
		{
			name: "no fields",
			input: `{}`,
			wantVal: User{},
		},
		{
			name: "string",
			input: `{"name":"Tom"}`,
			wantVal: User{Name: "Tom"},
		},
		{
			name: "nil bytes",
			input: []byte(nil),
		},
		{
			name: "empty bytes",
			input: []byte(""),
		},
		{
			name: "bytes",
			input: []byte(`{"name":"Tom"}`),
			wantVal: User{Name: "Tom"},
		},
		{
			name: "nil",
		},
		{
			name: "empty bytes ptr",
			input: ekit.ToPtr[[]byte]([]byte("")),
		},
		{
			name: "bytes ptr",
			input: ekit.ToPtr[[]byte]([]byte(`{"name":"Tom"}`)),
			wantVal: User{Name: "Tom"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := &JsonColumn{}
			err := js.Scan(tc.input)
			if tc.wantErr != "" {
				assert.EqualError(t, err, tc.wantErr)
				return
			} else {
				assert.Nil(t, err)
			}
			_, err = js.Value()
			assert.Nil(t, err)
			assert.EqualValues(t, tc.wantVal, js.Val)
		})
	}
}