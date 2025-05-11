package browse_history

import (
	"context"

	v1 "demo/api/browse_history/v1"
)

// IBrowseHistory 浏览历史记录接口定义
type IBrowseHistory interface {
	// V1 创建v1版本API
	V1() IBrowseHistoryV1
}

// IBrowseHistoryV1 浏览历史记录v1版接口定义
type IBrowseHistoryV1 interface {
	// List 获取浏览历史记录列表
	List(ctx context.Context, req *v1.BrowseHistoryListReq) (res *v1.BrowseHistoryListRes, err error)
	// Clear 清空浏览历史记录
	Clear(ctx context.Context, req *v1.BrowseHistoryClearReq) (res *v1.BrowseHistoryClearRes, err error)
	// Add 添加浏览历史记录
	Add(ctx context.Context, req *v1.BrowseHistoryAddReq) (res *v1.BrowseHistoryAddRes, err error)
	// Count 获取浏览历史记录数量
	Count(ctx context.Context, req *v1.BrowseHistoryCountReq) (res *v1.BrowseHistoryCountRes, err error)
}
