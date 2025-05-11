package exchange

import (
	"context"
	v1 "demo/api/exchange/v1"
	"demo/internal/service"
)

// Controller 兑换记录控制器
type Controller struct{}

// Get 获取兑换记录详情
func (c *Controller) Get(ctx context.Context, req *v1.ExchangeRecordGetReq) (res *v1.ExchangeRecordGetRes, err error) {
	return service.Exchange().Get(ctx, req)
}

// Create 创建兑换记录
func (c *Controller) Create(ctx context.Context, req *v1.ExchangeRecordCreateReq) (res *v1.ExchangeRecordCreateRes, err error) {
	return service.Exchange().Create(ctx, req)
}

// Update 更新兑换记录
func (c *Controller) Update(ctx context.Context, req *v1.ExchangeRecordUpdateReq) (res *v1.ExchangeRecordUpdateRes, err error) {
	return service.Exchange().Update(ctx, req)
}

// Delete 删除兑换记录
func (c *Controller) Delete(ctx context.Context, req *v1.ExchangeRecordDeleteReq) (res *v1.ExchangeRecordDeleteRes, err error) {
	return service.Exchange().Delete(ctx, req)
}

// UpdateStatus 更新兑换记录状态
func (c *Controller) UpdateStatus(ctx context.Context, req *v1.ExchangeRecordStatusUpdateReq) (res *v1.ExchangeRecordStatusUpdateRes, err error) {
	return service.Exchange().UpdateStatus(ctx, req)
}

// GetList 获取兑换记录列表
func (c *Controller) GetList(ctx context.Context, req *v1.ExchangeRecordListReq) (res *v1.ExchangeRecordListRes, err error) {
	return service.Exchange().GetList(ctx, req)
}

// WxGetPage 微信客户端分页获取兑换记录列表
func (c *Controller) WxGetPage(ctx context.Context, req *v1.WxExchangeRecordPageReq) (res *v1.WxExchangeRecordPageRes, err error) {
	return service.Exchange().WxGetPage(ctx, req)
}

// WxCreate 微信客户端创建兑换记录
func (c *Controller) WxCreate(ctx context.Context, req *v1.WxExchangeRecordCreateReq) (res *v1.WxExchangeRecordCreateRes, err error) {
	return service.Exchange().WxCreate(ctx, req)
}

// WxGetPublicList 微信客户端获取公开兑换记录列表
func (c *Controller) WxGetPublicList(ctx context.Context, req *v1.WxExchangeRecordPublicListReq) (res *v1.WxExchangeRecordPublicListRes, err error) {
	return service.Exchange().WxGetPublicList(ctx, req)
}
