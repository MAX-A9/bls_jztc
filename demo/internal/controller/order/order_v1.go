package order

import (
	"context"

	v1 "demo/api/payment/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ControllerV1 订单控制器V1版本
type ControllerV1 struct{}

// List 获取订单列表
func (c *ControllerV1) List(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error) {
	// 检查权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 验证用户角色是否为管理员
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看订单列表", nil))
	}

	return service.Order().List(ctx, req)
}

// Detail 获取订单详情
func (c *ControllerV1) Detail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error) {
	// 检查权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 验证用户角色是否为管理员
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看订单详情", nil))
	}

	return service.Order().Detail(ctx, req)
}

// Cancel 取消订单
func (c *ControllerV1) Cancel(ctx context.Context, req *v1.OrderCancelReq) (res *v1.OrderCancelRes, err error) {
	// 检查权限
	_, _, _, err = auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	return service.Order().Cancel(ctx, req)
}

// Delete 删除订单
func (c *ControllerV1) Delete(ctx context.Context, req *v1.OrderDeleteReq) (res *v1.OrderDeleteRes, err error) {
	// 检查权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 验证用户角色是否为管理员
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限删除订单", nil))
	}

	return service.Order().Delete(ctx, req)
}

// UpdateStatus 更新订单状态
func (c *ControllerV1) UpdateStatus(ctx context.Context, req *v1.UpdateOrderStatusReq) (res *v1.UpdateOrderStatusRes, err error) {
	// 检查权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(401, "未登录或登录已过期", nil))
	}

	// 验证用户角色是否为管理员
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限更新订单状态", nil))
	}

	return service.Order().UpdateStatus(ctx, req)
}

// NewV1 创建订单控制器V1版本实例
func NewV1() *ControllerV1 {
	return &ControllerV1{}
}
