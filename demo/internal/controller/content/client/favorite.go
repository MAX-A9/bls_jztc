package client

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"demo/api/content/client"
	v1 "demo/api/content/client/v1"
	"demo/internal/service"
	"demo/utility/auth"
)

// FavoriteController 收藏控制器
type FavoriteController struct{}

// NewFavorite 创建收藏控制器实例
func NewFavorite() client.FavoriteController {
	return &FavoriteController{}
}

// Add 添加收藏
func (c *FavoriteController) Add(ctx context.Context, req *v1.FavoriteAddReq) (res *v1.FavoriteAddRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	// 调用服务层处理业务逻辑
	return service.Favorite().Add(ctx, req)
}

// Cancel 取消收藏
func (c *FavoriteController) Cancel(ctx context.Context, req *v1.FavoriteCancelReq) (res *v1.FavoriteCancelRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	// 调用服务层处理业务逻辑
	return service.Favorite().Cancel(ctx, req)
}

// GetStatus 获取收藏状态
func (c *FavoriteController) GetStatus(ctx context.Context, req *v1.FavoriteStatusReq) (res *v1.FavoriteStatusRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	// 调用服务层处理业务逻辑
	return service.Favorite().GetStatus(ctx, req)
}

// GetList 获取收藏列表
func (c *FavoriteController) GetList(ctx context.Context, req *v1.FavoriteListReq) (res *v1.FavoriteListRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	// 调用服务层处理业务逻辑
	return service.Favorite().GetList(ctx, req)
}

// GetCount 获取收藏总数
func (c *FavoriteController) GetCount(ctx context.Context, req *v1.FavoriteCountReq) (res *v1.FavoriteCountRes, err error) {
	// 验证客户身份
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未授权访问")
	}

	// 调用服务层处理业务逻辑
	return service.Favorite().GetCount(ctx, req)
}
