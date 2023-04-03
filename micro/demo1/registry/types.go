package registry

import (
	"context"
)

//go:generate mockgen -package=mocks -destination=mocks/registry.mock.go -source=types.go Registry
type Registry interface {
	Register(ctx context.Context, ins ServiceInstance) error
	// Unregister(ctx context.Context, serviceName string) error
	Unregister(ctx context.Context, ins ServiceInstance) error
	ListService(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	// 可以考虑利用 ctx 来 close 掉返回的 channel
	// Subscribe(ctx context.Context, serviceName string) (<- chan Event, error)

	Subscribe(serviceName string) (<- chan Event, error)
	// 可有可无，不定义的话，具体的实现也可以额外的添加
	Close() error
}

// ServiceInstance 代表的是一个实例
type ServiceInstance struct {
	ServiceName string
	Address string
}

// type Event interface {
// 	Type() string
// 	Parameter()
// }

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeAdd
	EventTypeDelete
	EventTypeUpdate
	// EventTypeErr
)

type Event struct {
	Type     EventType
	Instance ServiceInstance
}


// 利用 ctx 来退出监听循环
// func Subscribe(ctx context.Context, serviceName string) (<- chan Event, error) {
// 	ch := make(chan Event)
// 	go func() {
// 		for {
// 			select {
// 			case <- ctx.Done():
// 				return
// 			case xxx:
// 				ch <- event
// 			}
// 		}
// 	}()
// 	return ch, nil
// }