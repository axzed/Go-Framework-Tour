package gzip

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompressor(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
	}{
		{
			name:  "hello world",
			input: []byte("hello, world"),
		},
	}
	c := Compressor{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := c.Compress(tc.input)
			require.NoError(t, err)
			data, err = c.Uncompress(data)
			require.NoError(t, err)
			assert.Equal(t, tc.input, data)
		})
	}
}
