package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGen(t *testing.T) {
	testCases := []struct {
		name    string
		def     ServiceDefinition
		wantErr error
		wantGen string
	}{
		{
			name: "user service",
			def:  ServiceDefinition{
				Package: "gen_http",
				Name: "UserService",
				Methods: []ServiceMethod{
					{
						Name: "Create",
						Path: "/Create",
						ReqTypeName: "User",
						RespTypeName: "User",
					},
					{
						Name: "Update",
						Path: "/Update",
						ReqTypeName: "User",
						RespTypeName: "UpdateUserResp",
					},
					{
						Name: "GetById",
						Path: "/GetById",
						ReqTypeName: "GetUserReq",
						RespTypeName: "User",
					},
					{
						Name: "DeleteById",
						Path: "/DeleteById",
						ReqTypeName: "DeleteByIdReq",
						RespTypeName: "DeleteByIdResp",
					},
				},
			},
			wantGen: `package gen_http

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

func (s *UserServiceGen) Create(ctx context.Context, req *User) (*User, error) {
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
	resp := &User{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

func (s *UserServiceGen) Update(ctx context.Context, req *User) (*UpdateUserResp, error) {
	url := s.Endpoint + s.Path + "/Update"
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

func (s *UserServiceGen) GetById(ctx context.Context, req *GetUserReq) (*User, error) {
	url := s.Endpoint + s.Path + "/GetById"
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
	resp := &User{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

func (s *UserServiceGen) DeleteById(ctx context.Context, req *DeleteByIdReq) (*DeleteByIdResp, error) {
	url := s.Endpoint + s.Path + "/DeleteById"
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
	resp := &DeleteByIdResp{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs := &bytes.Buffer{}
			err := Gen(bs, tc.def)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantGen, bs.String())
		})
	}
}
