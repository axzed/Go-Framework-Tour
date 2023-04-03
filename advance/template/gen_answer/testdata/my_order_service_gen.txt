package testdata

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type MyOrderServiceGen struct {
    Endpoint string
    Path string
	Client http.Client
}

func (s *MyOrderServiceGen) Create(ctx context.Context, req *CreateOrderReq) (*CreateOrderResp, error) {
	url := s.Endpoint + s.Path + "/Create"
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
	resp := &CreateOrderResp{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

