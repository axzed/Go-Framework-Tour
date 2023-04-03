package demo

import (
	"context"
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/message"
	"gitee.com/geektime-geekbang/geektime-go/micro/demo/serialize"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize/json"
	"net"
	"reflect"
	"strconv"
	"time"
)

type Server struct {
	services map[string]reflectionStub
	serializers []serialize.Serializer
}

func NewServer() *Server {
	res := &Server{
		services: map[string]reflectionStub{},
		// 一个 byte 表达的范围就是 -128, 127
		serializers: make([]serialize.Serializer, 32),
	}
	res.RegisterSerializer(json.Serializer{})
	return res
}

func (s *Server) MustRegister(service Service) {
	err := s.Register(service)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Register(service Service) error {
	s.services[service.Name()] = reflectionStub{
		value: reflect.ValueOf(service),
		serializers: s.serializers,
	}
	return nil
}

func (s *Server) RegisterSerializer(serializer serialize.Serializer) {
	s.serializers[serializer.Code()] = serializer
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			// 考虑输出日志，然后返回
			// return
			continue
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				// 这里考虑输出日志
				conn.Close()
				return
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqMsg, err:= ReadMsg(conn)
		if err != nil {
			return err
		}
		req := message.DecodeReq(reqMsg)

		resp := &message.Response{
			Version: req.Version,
			Compresser: req.Compresser,
			Serializer: req.Serializer,
			MessageId: req.MessageId,
		}
		// 可以考虑找到本地的服务，然后发起调用
		service, ok := s.services[req.ServiceName]
		if !ok {
			// 返回客户端一个错误信息
			resp.Error = []byte("找不到服务")
			resp.SetHeadLength()
			_, err = conn.Write(message.EncodeResp(resp))
			if err != nil {
				return err
			}
			continue
		}
		// 把参数传进去
		// context 建起来
		ctx := context.Background()
		var cancel func() = func() {}
		for key, value := range req.Meta {
			if key == "timeout" {
				deadline, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					// 返回客户端一个错误信息
					resp.Error = []byte(err.Error())
					resp.SetHeadLength()
					_, err = conn.Write(message.EncodeResp(resp))
					if err != nil {
						cancel()
						return err
					}
					cancel()
					continue
				}
				ctx, cancel = context.WithDeadline(ctx, time.UnixMilli(deadline))
			} else {
				ctx = context.WithValue(ctx, key, value)
			}
		}
		// 在这里可以检测超时
		data, err := service.invoke(ctx, req)
		cancel()
		if req.Meta["oneway"] != "" {
			continue
		}
		if err != nil {
			// 返回客户端一个错误信息
			resp.Error = []byte(err.Error())
			resp.SetHeadLength()
			_, err = conn.Write(message.EncodeResp(resp))
			if err != nil {
				cancel()
				return err
			}
			continue
		}
		// 在这里可以检测超时
		resp.SetHeadLength()
		resp.BodyLength = uint32(len(data))
		resp.Data = data
		data = message.EncodeResp(resp)
		// 在这里检测超时
		_, err = conn.Write(data)
		if err != nil {
			return err
		}
	}
}

type reflectionStub struct {
	value reflect.Value
	// 借鉴的是 bit array 或者 bitset 的思想
	serializers []serialize.Serializer
}


func (s *reflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	methodName := req.MethodName
	data := req.Data

	serializer := s.serializers[req.Serializer]
	if serializer == nil {
		// 返回客户端一个错误信息
		return nil, errors.New("micro: 不支持的序列化协议")
	}

	method := s.value.MethodByName(methodName)
	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())
	err := serializer.Decode(data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})
	// if req.Meta["oneway"] != "" {
	// 	return nil, nil
	// }
	if len(res) > 1 && !res[1].IsZero() {
		return nil, res[1].Interface().(error)
	}
	return serializer.Encode(res[0].Interface())
}