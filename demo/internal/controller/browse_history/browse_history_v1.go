package browse_history

import (
	"context"

	v1 "demo/api/browse_history/v1"
	"demo/internal/service"
)

// List 获取浏览历史记录列表
func (c *ControllerV1) List(ctx context.Context, req *v1.BrowseHistoryListReq) (res *v1.BrowseHistoryListRes, err error) {
	return service.BrowseHistory().List(ctx, req)
}

// Clear 清空浏览历史记录
func (c *ControllerV1) Clear(ctx context.Context, req *v1.BrowseHistoryClearReq) (res *v1.BrowseHistoryClearRes, err error) {
	return service.BrowseHistory().Clear(ctx, req)
}

// Add 添加浏览历史记录
func (c *ControllerV1) Add(ctx context.Context, req *v1.BrowseHistoryAddReq) (res *v1.BrowseHistoryAddRes, err error) {
	return service.BrowseHistory().Add(ctx, req)
}

// Count 获取浏览历史记录数量
func (c *ControllerV1) Count(ctx context.Context, req *v1.BrowseHistoryCountReq) (res *v1.BrowseHistoryCountRes, err error) {
	return service.BrowseHistory().Count(ctx, req)
}
