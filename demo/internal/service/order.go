package service

import (
	"context"
	"sync"

	v1 "demo/api/payment/v1"
)

// OrderService 订单服务接口
type OrderService interface {
	// List 获取订单列表
	List(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error)

	// Detail 获取订单详情
	Detail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error)

	// Cancel 取消订单
	Cancel(ctx context.Context, req *v1.OrderCancelReq) (res *v1.OrderCancelRes, err error)

	// Delete 删除订单
	Delete(ctx context.Context, req *v1.OrderDeleteReq) (res *v1.OrderDeleteRes, err error)

	// UpdateStatus 管理员更新订单状态
	UpdateStatus(ctx context.Context, req *v1.UpdateOrderStatusReq) (res *v1.UpdateOrderStatusRes, err error)

	// UpdateOrderStatus 更新订单状态
	UpdateOrderStatus(ctx context.Context, orderNo string, status int, transactionId string) error
}

var (
	orderServiceInstance OrderService
	orderOnce            sync.Once
)

// Order 获取订单服务实例
func Order() OrderService {
	if orderServiceInstance == nil {
		panic("订单服务未初始化")
	}
	return orderServiceInstance
}

// SetOrder 设置订单服务实例
func SetOrder(service OrderService) {
	orderOnce.Do(func() {
		orderServiceInstance = service
	})
}
