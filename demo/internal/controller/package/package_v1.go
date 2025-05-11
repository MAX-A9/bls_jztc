package package_controller

import (
	"context"

	v1 "demo/api/package/v1"
	"demo/internal/service"
)

// Controller 套餐控制器
type Controller struct{}

// V1 创建V1版本控制器
func (c *Controller) V1() *controllerV1 {
	return &controllerV1{}
}

// 套餐V1版本控制器
type controllerV1 struct{}

// List 获取套餐列表
func (c *controllerV1) List(ctx context.Context, req *v1.PackageListReq) (res *v1.PackageListRes, err error) {
	return service.Package().List(ctx, req)
}

// Detail 获取套餐详情
func (c *controllerV1) Detail(ctx context.Context, req *v1.PackageDetailReq) (res *v1.PackageDetailRes, err error) {
	return service.Package().Detail(ctx, req)
}

// Create 创建套餐
func (c *controllerV1) Create(ctx context.Context, req *v1.PackageCreateReq) (res *v1.PackageCreateRes, err error) {
	return service.Package().Create(ctx, req)
}

// Update 更新套餐
func (c *controllerV1) Update(ctx context.Context, req *v1.PackageUpdateReq) (res *v1.PackageUpdateRes, err error) {
	return service.Package().Update(ctx, req)
}

// Delete 删除套餐
func (c *controllerV1) Delete(ctx context.Context, req *v1.PackageDeleteReq) (res *v1.PackageDeleteRes, err error) {
	return service.Package().Delete(ctx, req)
}

// GetGlobalStatus 获取套餐总开关状态
func (c *controllerV1) GetGlobalStatus(ctx context.Context, req *v1.PackageGlobalStatusReq) (res *v1.PackageGlobalStatusRes, err error) {
	return service.Package().GetGlobalStatus(ctx, req)
}

// UpdateTopPackageGlobalStatus 更新置顶套餐总开关状态
func (c *controllerV1) UpdateTopPackageGlobalStatus(ctx context.Context, req *v1.TopPackageGlobalStatusUpdateReq) (res *v1.TopPackageGlobalStatusUpdateRes, err error) {
	return service.Package().UpdateTopPackageGlobalStatus(ctx, req)
}

// UpdatePublishPackageGlobalStatus 更新发布套餐总开关状态
func (c *controllerV1) UpdatePublishPackageGlobalStatus(ctx context.Context, req *v1.PublishPackageGlobalStatusUpdateReq) (res *v1.PublishPackageGlobalStatusUpdateRes, err error) {
	return service.Package().UpdatePublishPackageGlobalStatus(ctx, req)
}

// WxList 客户端获取套餐列表
func (c *controllerV1) WxList(ctx context.Context, req *v1.WxPackageListReq) (res *v1.WxPackageListRes, err error) {
	return service.Package().WxList(ctx, req)
}
