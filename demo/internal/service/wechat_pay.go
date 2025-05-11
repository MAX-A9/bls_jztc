package service

import (
	"context"
	"sync"

	v1 "demo/api/payment/v1/wechat"
	"demo/internal/model"
)

// WechatPayService 微信支付服务接口
type WechatPayService interface {
	// UnifiedOrder 统一下单
	UnifiedOrder(ctx context.Context, req *v1.WxPayUnifiedOrderReq) (res *v1.WxPayUnifiedOrderRes, err error)

	// HandleNotify 处理支付回调通知
	HandleNotify(ctx context.Context, notifyData []byte) (res *v1.WxPayNotifyRes, err error)

	// QueryOrder 查询订单
	QueryOrder(ctx context.Context, req *v1.WxPayOrderQueryReq) (res *v1.WxPayOrderQueryRes, err error)

	// 内部使用的方法
	// GetConfig 获取微信支付配置
	GetConfig(ctx context.Context) (*model.WxPayConfig, error)
}

var (
	wechatPayServiceInstance WechatPayService
	wechatPayOnce            sync.Once
)

// WechatPay 获取微信支付服务实例
func WechatPay() WechatPayService {
	if wechatPayServiceInstance == nil {
		panic("微信支付服务未初始化")
	}
	return wechatPayServiceInstance
}

// SetWechatPay 设置微信支付服务实例
func SetWechatPay(service WechatPayService) {
	wechatPayOnce.Do(func() {
		wechatPayServiceInstance = service
	})
}
