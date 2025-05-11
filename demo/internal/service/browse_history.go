package service

import (
	"context"

	v1 "demo/api/browse_history/v1"
)

// 浏览历史记录服务接口
type BrowseHistoryService interface {
	// List 获取浏览历史记录列表
	List(ctx context.Context, req *v1.BrowseHistoryListReq) (res *v1.BrowseHistoryListRes, err error)
	// Clear 清空浏览历史记录
	Clear(ctx context.Context, req *v1.BrowseHistoryClearReq) (res *v1.BrowseHistoryClearRes, err error)
	// Add 添加浏览历史记录
	Add(ctx context.Context, req *v1.BrowseHistoryAddReq) (res *v1.BrowseHistoryAddRes, err error)
	// Count 获取浏览历史记录数量
	Count(ctx context.Context, req *v1.BrowseHistoryCountReq) (res *v1.BrowseHistoryCountRes, err error)
}

var (
	localBrowseHistory BrowseHistoryService
)

// BrowseHistory 获取浏览历史记录服务
func BrowseHistory() BrowseHistoryService {
	if localBrowseHistory == nil {
		panic("请先调用 SetBrowseHistory 方法初始化浏览历史记录服务")
	}
	return localBrowseHistory
}

// SetBrowseHistory 设置浏览历史记录服务
func SetBrowseHistory(s BrowseHistoryService) {
	localBrowseHistory = s
}
