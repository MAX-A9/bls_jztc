package service

import (
	"context"

	v1 "demo/api/content/client/v1"
)

// FavoriteService 收藏服务接口
type FavoriteService interface {
	// Add 添加收藏
	Add(ctx context.Context, req *v1.FavoriteAddReq) (res *v1.FavoriteAddRes, err error)

	// Cancel 取消收藏
	Cancel(ctx context.Context, req *v1.FavoriteCancelReq) (res *v1.FavoriteCancelRes, err error)

	// GetStatus 获取收藏状态
	GetStatus(ctx context.Context, req *v1.FavoriteStatusReq) (res *v1.FavoriteStatusRes, err error)

	// GetList 获取收藏列表
	GetList(ctx context.Context, req *v1.FavoriteListReq) (res *v1.FavoriteListRes, err error)

	// GetCount 获取收藏总数
	GetCount(ctx context.Context, req *v1.FavoriteCountReq) (res *v1.FavoriteCountRes, err error)
}

var localFavorite FavoriteService

// Favorite 获取收藏服务实例
func Favorite() FavoriteService {
	if localFavorite == nil {
		panic("请先调用 SetFavorite 方法设置收藏服务实现")
	}
	return localFavorite
}

// SetFavorite 设置收藏服务实现
func SetFavorite(favorite FavoriteService) {
	localFavorite = favorite
}
