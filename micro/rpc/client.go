package rpc

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/compress"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/message"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize/json"
	"github.com/gotomicro/ekit/bean/option"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
)

var messageId uint32 = 0

type Client struct {
	connPool   pool.Pool
	serializer serialize.Serializer
	compressor compress.Compressor
}

func ClientWithSerializer(s serialize.Serializer) option.Option[Client] {
	return func(client *Client) {
		client.serializer = s
	}
}

func ClientWithCompressor(c compress.Compressor) option.Option[Client] {
	return func(client *Client) {
		client.compressor = c
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
		// 避免 nil 检测
		compressor: compress.DoNothingCompressor{},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var (
		resp *message.Response
		err  error
	)
	ch := make(chan struct{})

	go func() {
		bs := message.EncodeReq(req)
		resp, err = c.doInvoke(ctx, bs)
		ch <- struct{}{}
		close(ch)
	}()
	select {
	case <-ch:
		return resp, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) doInvoke(ctx context.Context, bs []byte) (*message.Response, error) {
	conn, err := c.connPool.Get()
	if err != nil {
		return nil, fmt.Errorf("client: 获得获取一个可用连接 %w", err)
	}
	// put back
	defer c.connPool.Put(conn)

	cn := conn.(net.Conn)

	_, err = cn.(net.Conn).Write(bs)
	if err != nil {
		return nil, err
	}

	if isOneway(ctx) {
		// 返回一个 error，防止有用户真的去接收结果
		return nil, errors.New("client: 这是 oneway 调用")
	}

	bs, err = ReadMsg(cn.(net.Conn))
	if err != nil {
		return nil, fmt.Errorf("client: 无法读取响应 %w", err)
	}
	resp := message.DecodeResp(bs)
	return resp, nil
}

func (c *Client) InitService(val Service) error {
	return setFuncField(c.serializer, c.compressor, val, c)
}

// 这个单独的拆出来，就是为了测试，我们可以考虑传入一个 mock proxy
func setFuncField(s serialize.Serializer, c compress.Compressor, val Service, proxy Proxy) error {
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
				out := reflect.Zero(field.Type.Out(0))
				inData, err := s.Encode(in)
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}

				inData, err = c.Compress(inData)
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}
				ctx := args[0].Interface().(context.Context)
				// 暂时先写死，后面我们考虑通用的链路元数据传递再重构
				meta := make(map[string]string, 2)
				oneway := isOneway(ctx)
				if oneway {
					meta["one-way"] = "true"
				}
				if deadline, ok := ctx.Deadline(); ok {
					// 传输字符串，需要更加多的空间
					meta["deadline"] = strconv.FormatInt(deadline.UnixMilli(), 10)
				}
				req := message.GetRequest()
				defer message.PutRequest(req)
				req.Meta = meta
				req.BodyLength = uint32(len(inData))
				req.MessageId = atomic.AddUint32(&messageId, 1)
				// 目前还没有支持压缩，需要你们作业支持
				req.Compresser = c.Code()
				req.Serializer = s.Code()
				req.ServiceName = val.ServiceName()
				req.Method = field.Name
				req.Data = inData
				req.SetHeadLength()
				resp, err := proxy.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}
				defer message.PutResponse(resp)

				var retErr error
				if len(resp.Error) > 0 {
					retErr = errors.New(string(resp.Error))
				}
				if len(resp.Data) > 0 {
					out = reflect.New(field.Type.Out(0).Elem())
					var data []byte
					data, err = c.Uncompress(resp.Data)
					if err != nil {
						return []reflect.Value{out, reflect.ValueOf(err)}
					}
					err = s.Decode(data, out.Interface())
					if err != nil {
						return []reflect.Value{out, reflect.ValueOf(err)}
					}
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
