package service

import (
	"context"

	v1 "demo/api/package/v1"
)

// PackageService 套餐服务接口
type PackageService interface {
	// List 获取套餐列表
	List(ctx context.Context, req *v1.PackageListReq) (res *v1.PackageListRes, err error)
	// Detail 获取套餐详情
	Detail(ctx context.Context, req *v1.PackageDetailReq) (res *v1.PackageDetailRes, err error)
	// Create 创建套餐
	Create(ctx context.Context, req *v1.PackageCreateReq) (res *v1.PackageCreateRes, err error)
	// Update 更新套餐
	Update(ctx context.Context, req *v1.PackageUpdateReq) (res *v1.PackageUpdateRes, err error)
	// Delete 删除套餐
	Delete(ctx context.Context, req *v1.PackageDeleteReq) (res *v1.PackageDeleteRes, err error)
	// GetGlobalStatus 获取套餐总开关状态
	GetGlobalStatus(ctx context.Context, req *v1.PackageGlobalStatusReq) (res *v1.PackageGlobalStatusRes, err error)
	// UpdateTopPackageGlobalStatus 更新置顶套餐总开关状态
	UpdateTopPackageGlobalStatus(ctx context.Context, req *v1.TopPackageGlobalStatusUpdateReq) (res *v1.TopPackageGlobalStatusUpdateRes, err error)
	// UpdatePublishPackageGlobalStatus 更新发布套餐总开关状态
	UpdatePublishPackageGlobalStatus(ctx context.Context, req *v1.PublishPackageGlobalStatusUpdateReq) (res *v1.PublishPackageGlobalStatusUpdateRes, err error)
	// WxList 客户端获取套餐列表
	WxList(ctx context.Context, req *v1.WxPackageListReq) (res *v1.WxPackageListRes, err error)
}

var (
	// localPackage 是PackageService接口的单例实现
	localPackage PackageService
)

// Package 返回PackageService接口的单例对象
func Package() PackageService {
	if localPackage == nil {
		panic("implement not found for interface PackageService, forgot register?")
	}
	return localPackage
}

// SetPackage 设置PackageService接口的实现
func SetPackage(s PackageService) {
	localPackage = s
}
