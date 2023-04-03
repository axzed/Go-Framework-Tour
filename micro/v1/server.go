package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	services map[string]*reflectionStub
}

func (s *Server) Start(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accept connection got error: %v", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		bytes, err := ReadMsg(conn)
		if err != nil {
			return
		}
		// go func() {
		u := &Request{}
		err = json.Unmarshal(bytes, u)
		resp, er := s.Invoke(context.Background(), u)
		if resp == nil {
			resp = &Response{}
		}
		if er != nil && len(resp.Error) == 0 {
			resp.Error = er.Error()
		}
		encode, er := s.encode(resp)
		if er != nil {
			fmt.Printf("encode resp failed: %v", er)
			return
		}
		_, er = conn.Write(encode)
		if er != nil {
			fmt.Printf("sending response failed: %v", er)
		}
	}
}

func (s *Server) encode(m interface{}) ([]byte, error) {
	respData, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal response failed")
		return nil, err
	}
	return EncodeMsg(respData), nil
}

func (s *Server) Invoke(ctx context.Context, req *Request) (*Response, error) {
	resp := &Response{}
	service, ok := s.services[req.ServiceName]
	if !ok {
		return resp, fmt.Errorf("server: 未找到服务, 服务名 %s", req.ServiceName)
	}
	respData, err := service.invoke(ctx, req.Method, req.Data)
	if err != nil {
		return resp, err
	}
	resp.Data = respData
	return resp, nil
}

func (s *Server) RegisterService(service Service) {
	s.services[service.ServiceName()] = &reflectionStub{
		s:     service,
		value: reflect.ValueOf(service),
	}
}

func NewServer() *Server {
	res := &Server{
		services: make(map[string]*reflectionStub, 4),
	}
	return res
}

type reflectionStub struct {
	s     Service
	value reflect.Value
}

func (s *reflectionStub) invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	method := s.value.MethodByName(methodName)
	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())
	err := json.Unmarshal(data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})
	if len(res) > 1 && !res[1].IsZero() {
		return nil, res[1].Interface().(error)
	}
	return json.Marshal(res[0].Interface())
}
