package testdata

import "context"

type (
	// OrderService 订单相关操作
	// 重新改一下名字
	// @HttpClient
	// @ServiceName MyOrderService
	OrderService interface {
		Create(ctx context.Context, req *CreateOrderReq) (*CreateOrderResp, error)
	}
)

type CreateOrderReq struct {

}

type CreateOrderResp struct {

}