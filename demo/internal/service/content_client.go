// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package service

import (
	"context"

	v1 "demo/api/content/client/v1"
)

// ContentClientService 客户端内容服务接口
type ContentClientService interface {
	// 获取分类列表
	CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error)

	// 微信客户端-闲置发布
	WxIdleCreate(ctx context.Context, req *v1.WxIdleCreateReq) (res *v1.WxIdleCreateRes, err error)

	// 微信客户端-信息发布
	WxInfoCreate(ctx context.Context, req *v1.WxInfoCreateReq) (res *v1.WxInfoCreateRes, err error)

	// 按地区获取内容列表
	RegionContentList(ctx context.Context, req *v1.RegionContentListReq) (res *v1.RegionContentListRes, err error)

	// 按地区获取闲置物品列表
	RegionIdleList(ctx context.Context, req *v1.RegionIdleListReq) (res *v1.RegionIdleListRes, err error)

	// 获取内容详情
	ContentPublicDetail(ctx context.Context, req *v1.ContentPublicDetailReq) (res *v1.ContentPublicDetailRes, err error)

	// 微信客户端-获取套餐列表
	WxClientPackageList(ctx context.Context, req *v1.WxClientPackageListReq) (res *v1.WxClientPackageListRes, err error)

	// 微信客户端-获取我的发布列表
	WxMyPublishList(ctx context.Context, req *v1.WxMyPublishListReq) (res *v1.WxMyPublishListRes, err error)

	// 微信客户端-获取我的发布数量
	WxMyPublishCount(ctx context.Context, req *v1.WxMyPublishCountReq) (res *v1.WxMyPublishCountRes, err error)
}

var localContentClient ContentClientService

// ContentClient 获取客户端内容服务实例
func ContentClient() ContentClientService {
	if localContentClient == nil {
		panic("请先调用service.SetContentClient设置客户端内容服务")
	}
	return localContentClient
}

// SetContentClient 设置客户端内容服务实例
func SetContentClient(s ContentClientService) {
	localContentClient = s
}
