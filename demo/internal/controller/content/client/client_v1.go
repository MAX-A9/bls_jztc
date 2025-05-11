package client

import (
	"context"

	"demo/api/content/client"
	v1 "demo/api/content/client/v1"
	"demo/internal/service"
)

// controllerV1 客户端内容控制器V1版本
type controllerV1 struct{}

// Favorite 获取收藏控制器实例
func (c *controllerV1) Favorite() client.FavoriteController {
	return NewFavorite()
}

// CategoryList 获取分类列表
func (c *controllerV1) CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	return service.ContentClient().CategoryList(ctx, req)
}

// WxIdleCreate 微信客户端-闲置发布
func (c *controllerV1) WxIdleCreate(ctx context.Context, req *v1.WxIdleCreateReq) (res *v1.WxIdleCreateRes, err error) {
	return service.ContentClient().WxIdleCreate(ctx, req)
}

// WxInfoCreate 微信客户端-信息发布
func (c *controllerV1) WxInfoCreate(ctx context.Context, req *v1.WxInfoCreateReq) (res *v1.WxInfoCreateRes, err error) {
	return service.ContentClient().WxInfoCreate(ctx, req)
}

// RegionContentList 按地区获取内容列表
func (c *controllerV1) RegionContentList(ctx context.Context, req *v1.RegionContentListReq) (res *v1.RegionContentListRes, err error) {
	return service.ContentClient().RegionContentList(ctx, req)
}

// RegionIdleList 按地区获取闲置物品列表
func (c *controllerV1) RegionIdleList(ctx context.Context, req *v1.RegionIdleListReq) (res *v1.RegionIdleListRes, err error) {
	return service.ContentClient().RegionIdleList(ctx, req)
}

// ContentPublicDetail 获取公开内容详情
func (c *controllerV1) ContentPublicDetail(ctx context.Context, req *v1.ContentPublicDetailReq) (res *v1.ContentPublicDetailRes, err error) {
	return service.ContentClient().ContentPublicDetail(ctx, req)
}

// WxClientPackageList 微信客户端-获取套餐列表
func (c *controllerV1) WxClientPackageList(ctx context.Context, req *v1.WxClientPackageListReq) (res *v1.WxClientPackageListRes, err error) {
	return service.ContentClient().WxClientPackageList(ctx, req)
}

// WxMyPublishList 微信客户端-获取我的发布列表
func (c *controllerV1) WxMyPublishList(ctx context.Context, req *v1.WxMyPublishListReq) (res *v1.WxMyPublishListRes, err error) {
	return service.ContentClient().WxMyPublishList(ctx, req)
}

// WxMyPublishCount 微信客户端-获取我的发布数量
func (c *controllerV1) WxMyPublishCount(ctx context.Context, req *v1.WxMyPublishCountReq) (res *v1.WxMyPublishCountRes, err error) {
	return service.ContentClient().WxMyPublishCount(ctx, req)
}
