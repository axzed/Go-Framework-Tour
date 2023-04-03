package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/domainobject/entity"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/service"
	usmocks "gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/service/mocks"
	"gitee.com/geektime-geekbang/geektime-go/web"
	"gitee.com/geektime-geekbang/geektime-go/web/session"
	"gitee.com/geektime-geekbang/geektime-go/web/session/cookie"
	"gitee.com/geektime-geekbang/geektime-go/web/session/memory"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	thttp "github.com/stretchr/testify/http"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

// 测试的核心是利用 mock 对象来构建 Handler

// 确保 Handler 自身逻辑是对的
func TestUserHandler_Login(t *testing.T) {
	// ctrl 这个是可以复用的
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct{
		name string
		// 一般就是声明一个 mock 函数字段，然后它会返回所有需要的东西
		mock func() service.UserService

		// 输入
		ctx *web.Context

		// Login 是一个没有返回值的方法，所以我们比较的都是 ctx 里面的 RespData 和 RespStatusCode
		// 还原成 Resp 比较容易写测试用例
		wantResp Resp
		wantCode int
	}{
		// 业务测试的测试用例设计和中间件的测试用例设计基本上差不多。
		// 理论上测试用例应该站在用户的角度来设计，
		// 但是大多数时候你只需要覆盖所有的分支就可以了
		{
			name: "invalid json",
			mock: func() service.UserService {
				return nil
			},
			ctx: func() *web.Context{
				body := bytes.NewBuffer([]byte(`{`))
				req, err := http.NewRequest(http.MethodPost, "/login", body)
				require.NoError(t, err)
				return &web.Context{
					Req: req,
					// 使用测试用的 writer，实际上用不到这个东西
					Resp: &thttp.TestResponseWriter{},
				}
			}(),
			wantResp: Resp{
				Msg: "解析请求失败",
			},
			wantCode: http.StatusBadRequest,
		},

		{
			name: "incorrect password",
			mock: func() service.UserService {
				us := usmocks.NewMockUserService(ctrl)
				us.EXPECT().Login(gomock.Any(), gomock.Any()).Return(entity.User{}, service.ErrInvalidUserOrPassword)
				return us
			},
			ctx: func() *web.Context{
				body := bytes.NewBuffer([]byte(`{}`))
				req, err := http.NewRequest(http.MethodPost, "/login", body)
				require.NoError(t, err)
				return &web.Context{
					Req: req,
					// 使用测试用的 writer，实际上用不到这个东西
					Resp: &thttp.TestResponseWriter{},
				}
			}(),
			wantResp: Resp{
				Msg: "账号或用户名输入错误",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "db error",
			mock: func() service.UserService {
				us := usmocks.NewMockUserService(ctrl)
				us.EXPECT().Login(gomock.Any(), gomock.Any()).Return(entity.User{}, errors.New("mock db error"))
				return us
			},
			ctx: func() *web.Context{
				body := bytes.NewBuffer([]byte(`{}`))
				req, err := http.NewRequest(http.MethodPost, "/login", body)
				require.NoError(t, err)
				return &web.Context{
					Req: req,
					// 使用测试用的 writer，实际上用不到这个东西
					Resp: &thttp.TestResponseWriter{},
				}
			}(),
			wantResp: Resp{
				Msg: "系统异常",
			},
			wantCode: http.StatusInternalServerError,
		},

		// 因为我们这里没有使用 mock 的 session 实现，所以模拟不了 session 操作失败
		// 在使用 mock 的 session 的情况下，可以检测对 session 的调用，来确保我们已经设置了 session
		{
			name: "success",
			mock: func() service.UserService {
				us := usmocks.NewMockUserService(ctrl)
				us.EXPECT().Login(gomock.Any(), gomock.Any()).Return(entity.User{
					Id: 123,
				}, nil)
				return us
			},
			ctx: func() *web.Context{
				body := bytes.NewBuffer([]byte(`{}`))
				req, err := http.NewRequest(http.MethodPost, "/login", body)
				require.NoError(t, err)
				return &web.Context{
					Req: req,
					// 使用测试用的 writer，实际上用不到这个东西
					Resp: &thttp.TestResponseWriter{},
				}
			}(),
			wantResp: Resp{
				Msg: "登录成功",
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			us := tc.mock()
			sessMgr := session.Manager{
				Store: memory.NewStore(time.Minute * 15),
				Propagator: cookie.NewPropagator("sessid"),
				SessCtxKey: "_sess",
			}
			handler := NewUserHandler(us, sessMgr)
			handler.Login(tc.ctx)
			var res Resp
			err := json.Unmarshal(tc.ctx.RespData, &res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantResp, res)
			assert.Equal(t, tc.wantCode, tc.ctx.RespStatusCode)
		})
	}
}