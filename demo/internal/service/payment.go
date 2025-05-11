package service

import (
	"context"
	"sync"

	v1 "demo/api/payment/v1"
)

// PaymentService 支付服务接口
type PaymentService interface {
	// GetConfig 获取支付配置
	GetConfig(ctx context.Context, req *v1.PaymentConfigReq) (res *v1.PaymentConfigRes, err error)

	// SaveConfig 保存支付配置
	SaveConfig(ctx context.Context, req *v1.SavePaymentConfigReq) (res *v1.SavePaymentConfigRes, err error)
}

var (
	paymentServiceInstance PaymentService
	paymentOnce            sync.Once
)

// Payment 获取支付服务实例
func Payment() PaymentService {
	if paymentServiceInstance == nil {
		panic("支付服务未初始化")
	}
	return paymentServiceInstance
}

// SetPayment 设置支付服务实例
func SetPayment(service PaymentService) {
	paymentOnce.Do(func() {
		paymentServiceInstance = service
	})
}
