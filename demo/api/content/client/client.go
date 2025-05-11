// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package client

import (
	"context"

	v1 "demo/api/content/client/v1"
)

// FavoriteController 收藏控制器接口
type FavoriteController interface {
	// Add 添加收藏
	Add(ctx context.Context, req *v1.FavoriteAddReq) (res *v1.FavoriteAddRes, err error)
	// Cancel 取消收藏
	Cancel(ctx context.Context, req *v1.FavoriteCancelReq) (res *v1.FavoriteCancelRes, err error)
	// GetStatus 获取收藏状态
	GetStatus(ctx context.Context, req *v1.FavoriteStatusReq) (res *v1.FavoriteStatusRes, err error)
	// GetList 获取收藏列表
	GetList(ctx context.Context, req *v1.FavoriteListReq) (res *v1.FavoriteListRes, err error)
	// GetCount 获取收藏数量
	GetCount(ctx context.Context, req *v1.FavoriteCountReq) (res *v1.FavoriteCountRes, err error)
}

// IClientV1 客户端内容控制器V1接口
type IClientV1 interface {
	// Favorite 获取收藏控制器实例
	Favorite() FavoriteController
	// CategoryList 获取分类列表
	CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error)
	// WxIdleCreate 微信客户端-闲置发布
	WxIdleCreate(ctx context.Context, req *v1.WxIdleCreateReq) (res *v1.WxIdleCreateRes, err error)
	// WxInfoCreate 微信客户端-信息发布
	WxInfoCreate(ctx context.Context, req *v1.WxInfoCreateReq) (res *v1.WxInfoCreateRes, err error)
	// RegionContentList 按地区获取内容列表
	RegionContentList(ctx context.Context, req *v1.RegionContentListReq) (res *v1.RegionContentListRes, err error)
	// RegionIdleList 按地区获取闲置物品列表
	RegionIdleList(ctx context.Context, req *v1.RegionIdleListReq) (res *v1.RegionIdleListRes, err error)
	// ContentPublicDetail 获取公开内容详情
	ContentPublicDetail(ctx context.Context, req *v1.ContentPublicDetailReq) (res *v1.ContentPublicDetailRes, err error)
	// WxClientPackageList 微信客户端-获取套餐列表
	WxClientPackageList(ctx context.Context, req *v1.WxClientPackageListReq) (res *v1.WxClientPackageListRes, err error)
	// WxMyPublishList 微信客户端-获取我的发布列表
	WxMyPublishList(ctx context.Context, req *v1.WxMyPublishListReq) (res *v1.WxMyPublishListRes, err error)
	// WxMyPublishCount 微信客户端-获取我的发布数量
	WxMyPublishCount(ctx context.Context, req *v1.WxMyPublishCountReq) (res *v1.WxMyPublishCountRes, err error)
} 