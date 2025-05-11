package service

import (
	"context"
	v1 "demo/api/exchange/v1"
)

// ExchangeService 兑换记录服务接口
type ExchangeService interface {
	// 后台管理接口
	Get(ctx context.Context, req *v1.ExchangeRecordGetReq) (res *v1.ExchangeRecordGetRes, err error)
	Create(ctx context.Context, req *v1.ExchangeRecordCreateReq) (res *v1.ExchangeRecordCreateRes, err error)
	Update(ctx context.Context, req *v1.ExchangeRecordUpdateReq) (res *v1.ExchangeRecordUpdateRes, err error)
	Delete(ctx context.Context, req *v1.ExchangeRecordDeleteReq) (res *v1.ExchangeRecordDeleteRes, err error)
	UpdateStatus(ctx context.Context, req *v1.ExchangeRecordStatusUpdateReq) (res *v1.ExchangeRecordStatusUpdateRes, err error)
	GetList(ctx context.Context, req *v1.ExchangeRecordListReq) (res *v1.ExchangeRecordListRes, err error)

	// 客户端接口
	WxGetPage(ctx context.Context, req *v1.WxExchangeRecordPageReq) (res *v1.WxExchangeRecordPageRes, err error)
	WxCreate(ctx context.Context, req *v1.WxExchangeRecordCreateReq) (res *v1.WxExchangeRecordCreateRes, err error)

	// 公开接口
	WxGetPublicList(ctx context.Context, req *v1.WxExchangeRecordPublicListReq) (res *v1.WxExchangeRecordPublicListRes, err error)
}

var localExchange ExchangeService

// Exchange 获取兑换记录服务
func Exchange() ExchangeService {
	if localExchange == nil {
		panic("implement not found for interface ExchangeService, forgot register?")
	}
	return localExchange
}

// RegisterExchange 注册兑换记录服务
func RegisterExchange(i ExchangeService) {
	localExchange = i
}
