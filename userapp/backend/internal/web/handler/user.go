package handler

import (
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/domainobject/entity"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/service"
	"gitee.com/geektime-geekbang/geektime-go/web"
	"gitee.com/geektime-geekbang/geektime-go/web/session"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	userIdKey = "user_id"
)

type UserHandler struct {
	service service.UserService
	sessMgr session.Manager
}

func NewUserHandler(us service.UserService, sessMgr session.Manager) *UserHandler {
	return &UserHandler{
		service: us,
		sessMgr: sessMgr,
	}
}
// vo => view object

var (
	service service.UserService
	sessMgr session.Manager
)

func SetService(s service.UserService) {
	service = s
}

// 强烈不建议
// func init() {
// 	service = service.NewUserService()
// }

func LoginV2(service service.UserService, sessMrg session.Manager, ctx *web.Context){

}

func LoginV1(service service.UserService, sessMrg session.Manager) web.HandleFunc {
	return func(ctx *web.Context) {
		req := loginReq{}
		err := ctx.BindJSON(&req)
		if err != nil {
			zap.L().Error("handler: 解析 JSON 数据格式失败", zap.Error(err))
			_ = ctx.RespJSON(http.StatusBadRequest, Resp{
				Msg: "解析请求失败",
			})
			return
		}
		usr, err := service.Login(ctx.Req.Context(), entity.User{
			Email: req.Email,
			Password: req.Password,
		})

		if errors.Is(err, service.ErrInvalidUserOrPassword) {
			zap.L().Error("登录失败", zap.Error(err))
			_ = ctx.RespJSON(http.StatusBadRequest, Resp{
				Msg: "账号或用户名输入错误",
			})
			return
		}

		if err != nil {
			zap.L().Error("登录失败，系统异常", zap.Error(err))
			_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
				Msg: "系统异常",
			})
			return
		}
		// 准备 session 了
		// session id 我们使用 uuid 就好了
		// 实际中你可以考虑将一些前端信息编码
		sess, err := h.sessMgr.InitSession(ctx, uuid.New().String())
		if err != nil {
			zap.L().Error("登录失败，初始化 session 失败", zap.Error(err))
			_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
				Msg: "系统异常",
			})
			return
		}

		err = sess.Set(ctx.Req.Context(), userIdKey, strconv.FormatUint(usr.Id, 10))
		if err != nil {
			zap.L().Error("登录失败，设置 session 失败", zap.Error(err))
			_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
				Msg: "系统异常",
			})
			return
		}

		err = ctx.RespJSON(http.StatusOK, Resp{
			Msg: "登录成功",
		})
	}
}

func Login(ctx *web.Context) {
	req := loginReq{}
	err := ctx.BindJSON(&req)
	if err != nil {
		zap.L().Error("handler: 解析 JSON 数据格式失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusBadRequest, Resp{
			Msg: "解析请求失败",
		})
		return
	}
	usr, err := service.Login(ctx.Req.Context(), entity.User{
		Email: req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		zap.L().Error("登录失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusBadRequest, Resp{
			Msg: "账号或用户名输入错误",
		})
		return
	}

	if err != nil {
		zap.L().Error("登录失败，系统异常", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}
	// 准备 session 了
	// session id 我们使用 uuid 就好了
	// 实际中你可以考虑将一些前端信息编码
	sess, err := h.sessMgr.InitSession(ctx, uuid.New().String())
	if err != nil {
		zap.L().Error("登录失败，初始化 session 失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}

	err = sess.Set(ctx.Req.Context(), userIdKey, strconv.FormatUint(usr.Id, 10))
	if err != nil {
		zap.L().Error("登录失败，设置 session 失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}

	err = ctx.RespJSON(http.StatusOK, Resp{
		Msg: "登录成功",
	})
}

func (h *UserHandler) Login(ctx *web.Context) {
	req := loginReq{}
	err := ctx.BindJSON(&req)
	if err != nil {
		zap.L().Error("handler: 解析 JSON 数据格式失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusBadRequest, Resp{
			Msg: "解析请求失败",
		})
		return
	}
	usr, err := h.service.Login(ctx.Req.Context(), entity.User{
		Email: req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		zap.L().Error("登录失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusBadRequest, Resp{
			Msg: "账号或用户名输入错误",
		})
		return
	}

	if err != nil {
		zap.L().Error("登录失败，系统异常", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}
	// 准备 session 了
	// session id 我们使用 uuid 就好了
	// 实际中你可以考虑将一些前端信息编码
	sess, err := h.sessMgr.InitSession(ctx, uuid.New().String())
	if err != nil {
		zap.L().Error("登录失败，初始化 session 失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}

	err = sess.Set(ctx.Req.Context(), userIdKey, strconv.FormatUint(usr.Id, 10))
	if err != nil {
		zap.L().Error("登录失败，设置 session 失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}

	err = ctx.RespJSON(http.StatusOK, Resp{
		Msg: "登录成功",
	})
}

func (h *UserHandler) Update(ctx *web.Context) {
	u := User{}
	err := ctx.BindJSON(&u)
	if err != nil{
		zap.L().Error("web: 解析 JSON 数据错误", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
	}

	uid, err := h.getId(ctx)
	if err != nil {
		zap.L().Error("handler: 无法获得 user id", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}

	err = h.service.EditProfile(ctx.Req.Context(), entity.User{
		// 一般是前端传了什么，这边就往下传什么
		Id: uid,
		Name: u.Name,
		Email: u.Email,
	})
	if err != nil {
		zap.L().Error("handler: 无法更新用户详情", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}

	// 可以考虑忽略，不过不嫌麻烦还是要和其它方法一样处理一下
	_ = ctx.RespJSON(http.StatusOK, Resp{
		Msg: "ok",
	})
}

func (h *UserHandler) Profile(ctx *web.Context) {
	uid, err := h.getId(ctx)
	if err != nil {
		zap.L().Error("handler: 无法获得 user id", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, Resp{
			Msg: "系统异常",
		})
		return
	}
	usr, err := h.service.FindById(ctx.Req.Context(), uid)
	if err != nil {
		// 这里已经没有必要区别用户存不存在了，因为 id 本身来自我们的 session
		zap.L().Error("web: 查找用户失败", zap.Error(err))
		_ = ctx.RespString(http.StatusInternalServerError, "system error")
		return
	}
	err = ctx.RespJSON(http.StatusOK, Resp{
		Data: User{
			Email: usr.Email,
			Name: usr.Name,
			Avatar: usr.Avatar,
		},
	})
	if err != nil {
		zap.L().Error("返回响应失败", zap.Error(err))
	}
}

func (h *UserHandler) SignUp(ctx *web.Context) {
	u := &signUpReq{}
	err := ctx.BindJSON(u)
	if err != nil {
		zap.L().Error("web: 解析 JSON 数据格式失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusBadRequest, Resp{
			Msg: "解析请求失败",
		})
		return
	}

	_, err = h.service.CreateUser(ctx.Req.Context(), entity.User{
		Email: u.Email,
		Password: u.Password,
	})
	if errors.Is(err, service.ErrInvalidNewUser) {
		zap.L().Error("创建用户失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusBadRequest, Resp{
			Msg: "用户输入错误",
		})
		return
	}
	if errors.Is(err, service.ErrDuplicateEmail) {
		zap.L().Error("创建用户失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusBadRequest, Resp{
			Msg: "邮箱已被注册",
		})
		return
	}
	if err != nil {
		zap.L().Error("创建用户失败", zap.Error(err))
		_ = ctx.RespJSON(http.StatusInternalServerError, &Resp{
			Msg: "创建用户失败",
		})
		return
	}
	// 如果你知道 RespOk 不可能返回 error，那么你可以考虑在这里忽略掉这个错误
	// 或者在 web 框架里面添加一个新的方法，MustRespOk，它会在 error 的时候 panic
	err = ctx.RespOk("创建成功")
	if err != nil {
		zap.L().Error("返回响应失败", zap.Error(err))
	}
}

func (h *UserHandler) getId(ctx *web.Context) (uint64, error){
	sess, err := h.sessMgr.GetSession(ctx)
	if err != nil {
		return 0, err
	}
	uidStr, err := sess.Get(ctx.Req.Context(), userIdKey)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(uidStr, 10, 64)
}
