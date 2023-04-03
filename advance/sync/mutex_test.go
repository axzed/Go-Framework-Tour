package sync

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayList_DeleteAt(t *testing.T) {
	testCases := []struct{
		name string
		index int
		input []int
		wantVals []int
	} {
		{
			// 删除第一个
			name: "first",
			index: 0,
			input: []int{1, 2, 3},
			wantVals: []int{2, 3},
		},
		{
			// 删除最后一个
			name: "last",
			index: 2,
			input: []int{1, 2, 3},
			wantVals: []int{1, 2},
		},
		{
			// 删除中间一个
			name:"middle",
			index: 2,
			input: []int{1, 2, 3, 4, 5},
			wantVals: []int{1, 2, 4, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewArrayList[int](12)
			a.vals = tc.input
			_ = a.DeleteAt(tc.index)
			assert.Equal(t, tc.wantVals, a.vals)
		})
	}
}
