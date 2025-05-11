package mall

import (
	"context"
	v1 "demo/api/mall/v1"
	"demo/internal/service"
)

// Controller 商城控制器
type Controller struct{}

// GetCategory 获取商城分类详情
func (c *Controller) GetCategory(ctx context.Context, req *v1.ShopCategoryGetReq) (res *v1.ShopCategoryGetRes, err error) {
	return service.Mall().GetCategory(ctx, req)
}

// CreateCategory 创建商城分类
func (c *Controller) CreateCategory(ctx context.Context, req *v1.ShopCategoryCreateReq) (res *v1.ShopCategoryCreateRes, err error) {
	return service.Mall().CreateCategory(ctx, req)
}

// UpdateCategory 更新商城分类
func (c *Controller) UpdateCategory(ctx context.Context, req *v1.ShopCategoryUpdateReq) (res *v1.ShopCategoryUpdateRes, err error) {
	return service.Mall().UpdateCategory(ctx, req)
}

// DeleteCategory 删除商城分类
func (c *Controller) DeleteCategory(ctx context.Context, req *v1.ShopCategoryDeleteReq) (res *v1.ShopCategoryDeleteRes, err error) {
	return service.Mall().DeleteCategory(ctx, req)
}

// UpdateCategoryStatus 更新商城分类状态
func (c *Controller) UpdateCategoryStatus(ctx context.Context, req *v1.ShopCategoryStatusUpdateReq) (res *v1.ShopCategoryStatusUpdateRes, err error) {
	return service.Mall().UpdateCategoryStatus(ctx, req)
}

// GetCategoryList 获取商城分类列表
func (c *Controller) GetCategoryList(ctx context.Context, req *v1.ShopCategoryListReq) (res *v1.ShopCategoryListRes, err error) {
	return service.Mall().GetCategoryList(ctx, req)
}

// WxGetCategoryList 微信客户端获取商城分类列表
func (c *Controller) WxGetCategoryList(ctx context.Context, req *v1.WxShopCategoryListReq) (res *v1.WxShopCategoryListRes, err error) {
	return service.Mall().WxGetCategoryList(ctx, req)
}

// SyncCategoriesProductCount 同步商品分类数量
func (c *Controller) SyncCategoriesProductCount(ctx context.Context, req *v1.ShopCategorySyncProductCountReq) (res *v1.ShopCategorySyncProductCountRes, err error) {
	return service.Mall().SyncCategoriesProductCount(ctx, req)
}
