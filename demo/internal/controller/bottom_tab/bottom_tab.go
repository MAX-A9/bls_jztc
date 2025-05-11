package bottom_tab

import (
	"context"
	v1 "demo/api/content/v1"
	"demo/internal/service"
)

// Controller 底部导航栏控制器
type Controller struct{}

// V1 创建V1版本控制器
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{}
}

// ControllerV1 V1版本底部导航栏控制器
type ControllerV1 struct{}

// GetBottomTabList 获取底部导航栏列表
func (c *ControllerV1) GetBottomTabList(ctx context.Context, req *v1.BottomTabListReq) (res *v1.BottomTabListRes, err error) {
	return service.BottomTab().GetBottomTabList(ctx, req)
}

// CreateBottomTab 创建底部导航项
func (c *ControllerV1) CreateBottomTab(ctx context.Context, req *v1.BottomTabCreateReq) (res *v1.BottomTabCreateRes, err error) {
	return service.BottomTab().CreateBottomTab(ctx, req)
}

// UpdateBottomTab 更新底部导航项
func (c *ControllerV1) UpdateBottomTab(ctx context.Context, req *v1.BottomTabUpdateReq) (res *v1.BottomTabUpdateRes, err error) {
	return service.BottomTab().UpdateBottomTab(ctx, req)
}

// DeleteBottomTab 删除底部导航项
func (c *ControllerV1) DeleteBottomTab(ctx context.Context, req *v1.BottomTabDeleteReq) (res *v1.BottomTabDeleteRes, err error) {
	return service.BottomTab().DeleteBottomTab(ctx, req)
}

// UpdateBottomTabStatus 更新底部导航项状态
func (c *ControllerV1) UpdateBottomTabStatus(ctx context.Context, req *v1.BottomTabStatusUpdateReq) (res *v1.BottomTabStatusUpdateRes, err error) {
	return service.BottomTab().UpdateBottomTabStatus(ctx, req)
}

// WxGetBottomTabList 微信客户端获取底部导航栏列表
func (c *ControllerV1) WxGetBottomTabList(ctx context.Context, req *v1.WxBottomTabListReq) (res *v1.WxBottomTabListRes, err error) {
	return service.BottomTab().WxGetBottomTabList(ctx, req)
}

// New 创建底部导航栏控制器实例
func New() *Controller {
	return &Controller{}
}
