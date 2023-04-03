package demo

import (
	"context"
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/message"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/serialize"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/serialize/json"
	"github.com/silenceper/pool"
	"net"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
)

var messageId uint32 = 0

type Client struct {
	connPool  pool.Pool
	serialzer serialize.Serializer
}

func NewClient(addr string) (*Client, error){
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 10,
		MaxCap: 100,
		MaxIdle: 50,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", addr)
		},
		IdleTimeout: time.Minute,
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		serialzer: json.Serializer{},
		connPool:  p,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	// 你可以在这里检测
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var (
		resp *message.Response
		err error
	)

	ch := make(chan struct{})
	go func() {
		resp, err = c.doInvoke(ctx, req)
		close(ch)
	}()
	select {
	case <- ctx.Done():
		return nil, ctx.Err()
	case <- ch:
		return resp, err
	}
}

func (c *Client) doInvoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	// 拿一个连接
	obj, err := c.connPool.Get()
	// 这个 error 是框架 error，而不是用户返回的 error
	if err != nil {
		return nil, err
	}
	conn := obj.(net.Conn)
	// 发请求

	// 发送请求过去
	data := message.EncodeReq(req)
	i, err := conn.Write(data)
	if err != nil {
		return  nil, err
	}

	// 可以检测超时

	if i != len(data) {
		return nil, errors.New("micro: 未写入全部数据")
	}
	// 读响应
	// 我怎么知道该读多长数据？相应地，服务端读请求，该读多长？
	// 先读长度字段

	// 读取全部的响应
	// 装响应的 bytes
	// if isOneway {
	// 	return
	// }
	respMsg, err := ReadMsg(conn)
	// 还可以在这里检测超时
	if err != nil {
		return nil, err
	}
	return message.DecodeResp(respMsg), nil
}

// 客户端
// 代码演示第一部分
// 1. 首先反射拿到 Request，核心是服务名字，方法名字和参数
// 2. 将 Request 进行编码，要注意序列化并且加上长度字段
// 3. 使用连接池，或者一个连接，把请求发过去

// 代码演示第四部分
// 4. 从连接里面读取响应，解析成结构体

// 服务端
// 代码演示第二部分
// 1. 启动一个服务器，监听一个端口
// 2. 读取长度字段，再根据长度，读完整个消息
// 3. 解析成 Request
// 4. 查找服务，会对应的方法
// 5. 构造方法对应的输入
// 6. 反射执行调用

// 代码演示第三部分
// 7. 编码响应
// 8. 写回响应



func (c *Client) InitService(service Service) error {
	// 你可以做校验，确保它必须是一个指向结构体的指针
	val := reflect.ValueOf(service).Elem()
	typ := reflect.TypeOf(service).Elem()
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			// 可以报错，也可以跳掉
			continue
		}
		// if fieldType.Type.Kind() != reflect.Func {
		// 	continue
		// }
		// 替换新的实现
		// 替换为一个新的方法实现
		fn := reflect.MakeFunc(fieldType.Type,
			func(args []reflect.Value) (results []reflect.Value) {
				// 实际上你在这里需要对 args 和 results 进行校验
				// 第一个返回值，真的返回值，指向 GetIdResp
				outType := fieldType.Type.Out(0)
				ctx := args[0].Interface().(context.Context)
				// if !ok {
				// 	return errors.Ne
				// }
				arg := args[1].Interface()

				bs, err := c.serialzer.Encode(arg)
				if err != nil {
					results = append(results, reflect.Zero(outType))
					// 这个是 error
					results = append(results, reflect.ValueOf(err))
					return
				}
				msgId := atomic.AddUint32(&messageId, 1)

				meta := make(map[string]string,2)
				// 能不能遍历 ctx 里面所有的 key？
				// 然后传递给服务端？答案是不能
				if isOneway(ctx) {
					meta = map[string]string{
						"oneway": "true",
					}
				}
				deadline, ok := ctx.Deadline()
				// 我确实有超时控制
				if ok {
					// 传秒
					meta["timeout"] = strconv.FormatInt(deadline.UnixMilli(), 10)
				}
				// 你要在这里把调用信息拼凑起来
				// 服务名，方法名，参数值，参数类型不需要
				req := &message.Request{
					// 要计算头部长度和响应体长度

					BodyLength: uint32(len(bs)),
					// 这里要构建完整
					Version: 0,
					Compresser: 0,
					Serializer: c.serialzer.Code(),
					MessageId: msgId,
					ServiceName: service.Name(),
					// 客户端和服务端可能叫不一样的名字
					// ServiceName: typ.PkgPath() + typ.Name(),
					// 服务名从哪里来？
					// 对应的是字段名
					MethodName: fieldType.Name,
					Data: bs,
					Meta: meta,
				}
				req.CalHeadLength()
				resp, err := c.Invoke(ctx, req)
				// if isOneway {
				// 	return
				// }

				if err != nil {
					results = append(results, reflect.Zero(outType))
					// 这个是 error
					results = append(results, reflect.ValueOf(err))
					return
				}
				// 第一个返回值，真的返回值，指向 GetIdResp
				first := reflect.New(outType).Interface()
				// 有没有可能请求是 json 序列化的，但是响应是 protobuf 序列化的
				// 如果你要支持这种，你在这里就要有一个查找 serializer 的过程
				err = c.serialzer.Decode(resp.Data, first)
				// if len(resp.Data) > 0 {
				//
				// }
				results = append(results, reflect.ValueOf(first).Elem())

				// 这个是 error
				if err != nil {
					results = append(results, reflect.ValueOf(err))
				} else {
					results = append(results,  reflect.Zero(reflect.TypeOf(new(error)).Elem()))
				}

				return
		})
		fieldValue.Set(fn)
	}
	return nil
}



type Service interface {
	Name() string
}