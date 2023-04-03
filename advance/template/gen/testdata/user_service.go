package testdata

import "context"

type (
	// UserService 定义了和 User 有关的操作
	// @HttpClient
	UserService interface {
		// Get 查询。虽然理论上我们应该使用 http GET 方法，但是为了简化处理，我们都用 POST
		Get(ctx context.Context, req *GetUserReq) (*GetUserResp, error)
		// Update 更新用户
		// @Path /user/update
		Update(ctx context.Context, req *UpdateUserReq) (*UpdateUserResp, error)
	}

	// UserServiceIgnore 会被忽略掉，因为它没有 HttpClient 注解
	UserServiceIgnore interface {
		Get(ctx context.Context, req *GetUserReq) (*GetUserResp, error)
		Update(ctx context.Context, req *UpdateUserReq) (*UpdateUserResp, error)
	}
)

type UpdateUserReq struct {

}

type UpdateUserResp struct {

}

type GetUserReq struct {

}

type GetUserResp struct {

}