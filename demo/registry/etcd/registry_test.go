package etcd

import (
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
)

func TestRegistry_Subscribe(t *testing.T) {
	testCases := []struct{
		name string
		mock func() clientv3.Watcher
		wantErr error
	} {
		{
			// mock: func() clientv3.Watcher {
			// 	watcher := mocks.NewMockWatcher()
			// },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &Registry{
				client: &clientv3.Client{
					Watcher: tc.mock(),
				},
			}
			ch, err := r.Subscribe("service-name")
			assert.Equal(t, tc.wantErr, err)
			event := <- ch
			// 你在这里进一步断言你预期中 event，还要测试 close 的例子
			log.Println(event)
		})
	}
}
