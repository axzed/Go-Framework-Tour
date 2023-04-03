package testdata

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type UserServiceGen struct {
    Endpoint string
    Path string
	Client http.Client
}

func (s *UserServiceGen) Get(ctx context.Context, req *GetUserReq) (*GetUserResp, error) {
	url := s.Endpoint + s.Path + "/Get"
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	body.Write(bs)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := s.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(httpResp.Body)
	resp := &GetUserResp{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

func (s *UserServiceGen) Update(ctx context.Context, req *UpdateUserReq) (*UpdateUserResp, error) {
	url := s.Endpoint + s.Path + "/user/update"
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	body.Write(bs)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := s.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(httpResp.Body)
	resp := &UpdateUserResp{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

