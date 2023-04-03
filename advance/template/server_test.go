package template

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type User struct {
	Id   int
	Name string
}

type GetUserReq struct {
	Id int
}

type DeleteByIdReq struct {
	Id int
}

type DeleteByIdResp struct {
	Ok bool
}

type UpdateUserResp struct {
	Ok bool
}

type UserService interface {
	Create(ctx context.Context, req *User) (*User, error)
	Update(ctx context.Context, req *User) (*UpdateUserResp, error)
	GetById(ctx context.Context, req *GetUserReq) (*User, error)
	DeleteById(ctx context.Context, req *DeleteByIdReq) (*DeleteByIdResp, error)
}

type userService struct {
	Endpoint string
	Path     string
	Client   *http.Client
}

func (u *userService) Create(ctx context.Context, req *User) (*User, error) {
	url := u.Endpoint + u.Path + "/Create"
	bs, err := json.Marshal(req)
	body := &bytes.Buffer{}
	body.Write(bs)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := u.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(httpResp.Body)
	resp := &User{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

func (u *userService) Update(ctx context.Context, req *User) (*UpdateUserResp, error) {
	url := u.Endpoint + u.Path + "/Update"
	bs, err := json.Marshal(req)
	body := &bytes.Buffer{}
	body.Write(bs)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := u.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(httpResp.Body)
	resp := &UpdateUserResp{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

func (u *userService) GetById(ctx context.Context, req *GetUserReq) (*User, error) {
	url := u.Endpoint + u.Path + "/GetById"
	bs, err := json.Marshal(req)
	body := &bytes.Buffer{}
	body.Write(bs)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := u.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(httpResp.Body)
	resp := &User{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

func (u *userService) DeleteById(ctx context.Context, req *DeleteByIdReq) (*DeleteByIdResp, error) {
	url := u.Endpoint + u.Path + "/DeleteById"
	bs, err := json.Marshal(req)
	body := &bytes.Buffer{}
	body.Write(bs)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := u.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(httpResp.Body)
	resp := &DeleteByIdResp{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}

func TestUserService(t *testing.T) {
	go startTestServer()
	// sleep 一下，等待服务器启动
	time.Sleep(3 * time.Second)

	us := &UserServiceGen{
		Endpoint: "http://localhost:8080",
		Path: "/user",
		Client: http.Client{},
	}
	u, err := us.Create(context.Background(), &User{Name: "Tom"})
	assert.Nil(t, err)
	assert.Equal(t, &User{Id: 12, Name: "Tom"}, u)

	uRes, err := us.Update(context.Background(), &User{Name: "Jerry"})
	assert.Nil(t, err)
	assert.True(t, uRes.Ok)

	u, err = us.GetById(context.Background(), &GetUserReq{Id: 12})
	assert.Nil(t, err)
	assert.Equal(t, &User{Id: 12, Name: "Tom"}, u)

	dRes, err:= us.DeleteById(context.Background(), &DeleteByIdReq{Id: 12})
	assert.Nil(t, err)
	assert.True(t, dRes.Ok)
}

func startTestServer() {
	http.HandleFunc("/user/Create", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		u := &User{}
		fmt.Printf("Create user: %s \n", string(bs))
		err = json.Unmarshal(bs, u)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		u.Id = 12
		u.Name = "Tom"
		bs, err = json.Marshal(u)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write(bs)
	})

	http.HandleFunc("/user/Update", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("Update user: %s \n", string(bs))
		_, _ = writer.Write([]byte(`{"ok": true}`))
	})

	http.HandleFunc("/user/GetById", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		u := &User{}
		err = json.Unmarshal(bs, u)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		u.Name = "Tom"
		bs, err = json.Marshal(u)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		_, _ = writer.Write(bs)
	})

	http.HandleFunc("/user/DeleteById", func(writer http.ResponseWriter, request *http.Request) {
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("Delete user: %s \n", string(bs))
		_, _ = writer.Write([]byte(`{"ok": true}`))
	})
	_ = http.ListenAndServe(":8080", nil)
}


// 从 gen_http_test 里面复制过来的
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



