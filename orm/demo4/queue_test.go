package orm

// func TestPriorityQueue_Enqueue(t *testing.T) {
// 	testCases := []struct{
// 		name string
//
//
// 		// 构造的输入
// 		q *PriorityQueue[int]
// 		input int
//
//
// 		// 期望的返回
// 		wantErr error
//
// 		// 内部细节
// 		wantCap int
// 		wantData []int
//
// 	}{
//
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			err := tc.q.Enqueue(tc.input)
// 			require.Equal(t, tc.wantErr, err)
//
// 			// 为了验证下标为 0 的位置，不放用户数据
// 			require.Equal(t, tc.wantCap, tc.q.capacity)
//
// 			// 为了验证你的代码保持住了堆结构，即入队出队之后，你要维持住大顶堆或者小顶堆的结构
// 			require.Equal(t, tc.wantData, tc.q.data)
// 		})
// 	}
// }
