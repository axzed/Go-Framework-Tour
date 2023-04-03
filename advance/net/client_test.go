package net

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Send(t *testing.T) {
	testCases := []struct {
		req  string
		resp string
	}{
		{
			req:  "hello",
			resp: "hello, from response",
		},
		{
			req:  "aaa bbb cc \n",
			resp: "aaa bbb cc \n, from response",
		},
	}

	c := &Client{
		addr: "localhost:8080",
	}
	for _, tc := range testCases {
		t.Run(tc.req, func(t *testing.T) {
			resp, err := c.Send(tc.req)
			assert.Nil(t, err)
			assert.Equal(t, tc.resp, resp)
		})
	}
}

func TestConnect(t *testing.T) {
	err := Connect("localhost:8080")
	assert.Nil(t, err)
}
