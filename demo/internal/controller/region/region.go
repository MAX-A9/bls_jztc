package region

import (
	"context"

	v1 "demo/api/region/v1"
	"demo/internal/service"
)

// Controller 地区管理控制器
type Controller struct{}

// New 创建地区管理控制器
func New() *Controller {
	return &Controller{}
}

// V1 获取V1版本API控制器
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{}
}

// ControllerV1 地区管理V1控制器
type ControllerV1 struct{}

// List 获取地区列表
func (c *ControllerV1) List(ctx context.Context, req *v1.RegionListReq) (res *v1.RegionListRes, err error) {
	return service.Region().List(ctx, req)
}

// Detail 获取地区详情
func (c *ControllerV1) Detail(ctx context.Context, req *v1.RegionDetailReq) (res *v1.RegionDetailRes, err error) {
	return service.Region().Detail(ctx, req)
}

// Create 创建地区
func (c *ControllerV1) Create(ctx context.Context, req *v1.RegionCreateReq) (res *v1.RegionCreateRes, err error) {
	return service.Region().Create(ctx, req)
}

// Update 更新地区
func (c *ControllerV1) Update(ctx context.Context, req *v1.RegionUpdateReq) (res *v1.RegionUpdateRes, err error) {
	return service.Region().Update(ctx, req)
}

// Delete 删除地区
func (c *ControllerV1) Delete(ctx context.Context, req *v1.RegionDeleteReq) (res *v1.RegionDeleteRes, err error) {
	return service.Region().Delete(ctx, req)
}

// WxClientList 客户端获取地区列表
func (c *ControllerV1) WxClientList(ctx context.Context, req *v1.WxClientRegionListReq) (res *v1.WxClientRegionListRes, err error) {
	return service.Region().WxClientList(ctx, req)
}
