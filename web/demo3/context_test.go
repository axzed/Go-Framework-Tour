package web

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestContext_BindJSON(t *testing.T) {
	testCases := []struct {
		name string

		req   *http.Request
		input any

		wantErr error
		wantVal any
	}{
		{
			name: "happy case",
			req: func() *http.Request {
				mockData := &bytes.Buffer{}
				mockData.WriteString(`{"name": "Tom"}`)
				req, err := http.NewRequest(http.MethodPost, "/user", mockData)
				if err != nil {
					t.Fatal(err)
				}
				return req
			}(),
			input:   &User{},
			wantVal: &User{Name: "Tom"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &Context{
				Req: tc.req,
			}
			err := ctx.BindJSON(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, tc.input)
		})
	}
}
