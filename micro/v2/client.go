package rpc

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/message"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize/json"
	"github.com/gotomicro/ekit/bean/option"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"sync/atomic"
	"time"
)

var messageId uint32 = 0

type Client struct {
	connPool   pool.Pool
	serializer serialize.Serializer
}

func ClientWithSerializer(s serialize.Serializer) option.Option[Client] {
	return func(client *Client) {
		client.serializer = s
	}
}

func NewClient(address string, opts ...option.Option[Client]) (*Client, error) {
	poolConfig := &pool.Config{
		InitialCap:  5,
		MaxIdle:     20,
		MaxCap:      30,
		Factory:     func() (interface{}, error) { return net.Dial("tcp", address) },
		Close:       func(v interface{}) error { return v.(net.Conn).Close() },
		IdleTimeout: time.Minute,
	}
	connPool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		return nil, err
	}

	res := &Client{
		connPool:   connPool,
		serializer: json.Serializer{},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	conn, err := c.connPool.Get()
	if err != nil {
		return nil, fmt.Errorf("client: 获得获取一个可用连接 %w", err)
	}
	// put back
	defer c.connPool.Put(conn)
	cn := conn.(net.Conn)
	bs := message.EncodeReq(req)
	_, err = cn.(net.Conn).Write(bs)
	if err != nil {
		return nil, err
	}

	bs, err = ReadMsg(cn.(net.Conn))
	if err != nil {
		return nil, fmt.Errorf("client: 无法读取响应 %w", err)
	}
	resp := message.DecodeResp(bs)
	return resp, nil
}

func (c *Client) InitService(val Service) error {
	return setFuncField(c.serializer, val, c)
}

// 这个单独的拆出来，就是为了测试，我们可以考虑传入一个 mock proxy
func setFuncField(s serialize.Serializer, val Service, proxy Proxy) error {
	v := reflect.ValueOf(val)
	ele := v.Elem()
	t := ele.Type()
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		fieldValue := ele.Field(i)
		if fieldValue.CanSet() {
			fn := func(args []reflect.Value) (results []reflect.Value) {
				in := args[1].Interface()
				inData, err := s.Encode(in)
				out := reflect.Zero(field.Type.Out(0))
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}

				req := &message.Request{
					BodyLength: uint32(len(inData)),
					MessageId:  atomic.AddUint32(&messageId, 1),
					// 目前还没有支持压缩，需要你们作业支持
					// Compresser:
					Serializer:  s.Code(),
					ServiceName: val.ServiceName(),
					Method:      field.Name,
					Data:        inData,
				}
				req.SetHeadLength()

				resp, err := proxy.Invoke(args[0].Interface().(context.Context), req)
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}
				var retErr error
				if len(resp.Error) > 0 {
					retErr = errors.New(string(resp.Error))
				}
				if len(resp.Data) > 0 {
					out = reflect.New(field.Type.Out(0))
					err = s.Decode(resp.Data, out.Interface())
					if err != nil {
						return []reflect.Value{out, reflect.ValueOf(err)}
					}
					out = out.Elem()
				}

				var errVal reflect.Value
				if retErr == nil {
					errVal = reflect.Zero(reflect.TypeOf(new(error)).Elem())
				} else {
					errVal = reflect.ValueOf(retErr)
				}
				return []reflect.Value{out, errVal}
			}
			fieldValue.Set(reflect.MakeFunc(field.Type, fn))
		}
	}
	return nil
}
