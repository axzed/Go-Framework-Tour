package micro

import (
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/registry"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"time"
)

type ClientOption func(client *Client)

type Client struct {
	rb       resolver.Builder
	insecure bool
	balancer balancer.Builder
}

func NewClient(opts ...ClientOption) *Client {
	client := &Client{}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *Client) Dial(ctx context.Context, service string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithResolvers(c.rb)}
	address := fmt.Sprintf("registry:///%s", service)
	if c.insecure {
		opts = append(opts, grpc.WithInsecure())
	}
	if c.balancer != nil {
		opts = append(opts, grpc.WithDefaultServiceConfig(
			fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`,
				c.balancer.Name())))
	}
	return grpc.DialContext(ctx, address, opts...)
}

func ClientWithRegistry(r registry.Registry, timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.rb = NewResolverBuilder(r, timeout)
	}
}

func ClientWithInsecure() ClientOption {
	return func(client *Client) {
		client.insecure = true
	}
}

func ClientWithPickerBuilder(name string, b base.PickerBuilder) ClientOption {
	return func(client *Client) {
		builder := base.NewBalancerBuilder(name, b, base.Config{HealthCheck: true})
		balancer.Register(builder)
		client.balancer = builder
	}
}

// 伪代码
// func (c *Client) DialPsu(ctx context.Context, service string) (*grpc.ClientConn, error) {
// 	resolver := c.rb
//
// 	grpc.DialContext(ctx,
// 		"registry:///user-service",
// 		grpc.WithResolvers(resolver))
// }
