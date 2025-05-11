package service

import (
	"context"
	v1 "demo/api/mall/v1"
)

// MallService 商城服务接口
type MallService interface {
	// 商城分类管理接口
	GetCategory(ctx context.Context, req *v1.ShopCategoryGetReq) (res *v1.ShopCategoryGetRes, err error)
	CreateCategory(ctx context.Context, req *v1.ShopCategoryCreateReq) (res *v1.ShopCategoryCreateRes, err error)
	UpdateCategory(ctx context.Context, req *v1.ShopCategoryUpdateReq) (res *v1.ShopCategoryUpdateRes, err error)
	DeleteCategory(ctx context.Context, req *v1.ShopCategoryDeleteReq) (res *v1.ShopCategoryDeleteRes, err error)
	UpdateCategoryStatus(ctx context.Context, req *v1.ShopCategoryStatusUpdateReq) (res *v1.ShopCategoryStatusUpdateRes, err error)
	GetCategoryList(ctx context.Context, req *v1.ShopCategoryListReq) (res *v1.ShopCategoryListRes, err error)

	// 商品分类同步接口
	SyncCategoriesProductCount(ctx context.Context, req *v1.ShopCategorySyncProductCountReq) (res *v1.ShopCategorySyncProductCountRes, err error)

	// 客户端商城分类接口
	WxGetCategoryList(ctx context.Context, req *v1.WxShopCategoryListReq) (res *v1.WxShopCategoryListRes, err error)
}

var localMall MallService

// Mall 获取商城服务
func Mall() MallService {
	if localMall == nil {
		panic("implement not found for interface MallService, forgot register?")
	}
	return localMall
}

// RegisterMall 注册商城服务
func RegisterMall(i MallService) {
	localMall = i
}
