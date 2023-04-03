package rpc

import (
	"context"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/compress"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/message"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize"
	"gitee.com/geektime-geekbang/geektime-go/micro/rpc/serialize/json"
	"net"
	"reflect"
	"strconv"
	"time"
)

type Server struct {
	services    map[string]*reflectionStub
	serializers []serialize.Serializer
	compressors []compress.Compressor

	listener net.Listener
}

func (s *Server) Start(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = ln
	for {
		conn, err := ln.Accept()
		// 关闭了
		if err == net.ErrClosed {
			return nil
		}
		if err != nil {
			fmt.Printf("accept connection got error: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		bytes, err := ReadMsg(conn)
		if err != nil {
			return
		}
		req := message.DecodeReq(bytes)
		ctx := context.Background()
		deadline, err := strconv.ParseInt(req.Meta["deadline"], 10, 64)
		cancel := func() {}
		if err == nil {
			ctx, cancel = context.WithDeadline(ctx, time.UnixMilli(deadline))
		}
		resp, er := s.Invoke(ctx, req)
		if req.Meta["one-way"] == "true" {
			// 什么也不需要处理。
			// 这样就相当于直接把连接资源释放了，去接收下一个请求了
			cancel()
			continue
		}
		if er != nil {
			resp = message.GetResponse()
			// 服务器本身出错了
			resp.Error = []byte(fmt.Errorf("rpc-server: 服务器异常 %w", er).Error())
			// 计算一下长度
			resp.SetHeadLength()
		}
		encode := message.EncodeResp(resp)
		message.PutResponse(resp)
		_, er = conn.Write(encode)
		if er != nil {
			fmt.Printf("sending response failed: %v", er)
		}
		cancel()
	}
}

func (s *Server) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	service, ok := s.services[req.ServiceName]
	if !ok {
		return nil, fmt.Errorf("server: 未找到服务, 服务名 %s", req.ServiceName)
	}
	return service.invoke(ctx, req)
}

func (s *Server) RegisterSerializer(serializer serialize.Serializer) {
	s.serializers[serializer.Code()] = serializer
}

func (s *Server) RegisterCompressor(c compress.Compressor) {
	s.compressors[c.Code()] = c
}

func (s *Server) RegisterService(service Service) {
	val := reflect.ValueOf(service)
	typ := reflect.TypeOf(service)
	methods := make(map[string]reflect.Value, val.NumMethod())
	for i := 0; i < val.NumMethod(); i++ {
		methodType := typ.Method(i)
		methods[methodType.Name] = val.Method(i)
	}
	s.services[service.ServiceName()] = &reflectionStub{
		s:           service,
		methods:     methods,
		serializers: s.serializers,
		compressors: s.compressors,
	}
}

func NewServer() *Server {
	res := &Server{
		services: make(map[string]*reflectionStub, 4),
		// 一个字节，最多有 256 个实现，直接做成一个简单的 bit array 的东西
		serializers: make([]serialize.Serializer, 256),
		compressors: make([]compress.Compressor, 256),
	}
	// 注册最基本的序列化协议
	res.RegisterSerializer(json.Serializer{})
	res.RegisterCompressor(compress.DoNothingCompressor{})
	return res
}

type reflectionStub struct {
	s           Service
	serializers []serialize.Serializer
	compressors []compress.Compressor
	methods     map[string]reflect.Value
}

func (r *reflectionStub) invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	method, ok := r.methods[req.Method]
	if !ok {
		return nil, fmt.Errorf("server: 未找到目标服务方法 %s", req.Method)
	}
	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())

	c := r.compressors[req.Compresser]
	data, err := c.Uncompress(req.Data)
	if err != nil {
		return nil, err
	}
	s := r.serializers[req.Serializer]
	err = s.Decode(data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})

	respData, err := s.Encode(res[0].Interface())
	if err != nil {
		// 服务器本身的错误
		return nil, err
	}
	respData, err = c.Compress(respData)
	if err != nil {
		return nil, err
	}
	resp := message.GetResponse()
	resp.BodyLength = uint32(len(respData))
	resp.MessageId = req.MessageId
	resp.Compresser = req.Compresser
	// 理论上来说，这里可以换一种序列化协议，但是没必要暴露这种功能给用户
	resp.Serializer = req.Serializer
	resp.Data = respData
	if !res[1].IsZero() {
		resp.Error = []byte(res[1].Interface().(error).Error())
	}
	resp.SetHeadLength()
	return resp, nil
}
