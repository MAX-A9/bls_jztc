package payment

import (
	"context"

	v1 "demo/api/payment/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ControllerV1 支付控制器V1版本
type ControllerV1 struct{}

// GetConfig 获取支付配置
func (c *ControllerV1) GetConfig(ctx context.Context, req *v1.PaymentConfigReq) (res *v1.PaymentConfigRes, err error) {
	// 检查权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil || role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看支付配置", nil))
	}

	return service.Payment().GetConfig(ctx, req)
}

// SaveConfig 保存支付配置
func (c *ControllerV1) SaveConfig(ctx context.Context, req *v1.SavePaymentConfigReq) (res *v1.SavePaymentConfigRes, err error) {
	// 权限验证：只有管理员才能修改配置
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(403, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限修改支付配置", nil))
	}

	// 调用服务实现保存配置
	res, err = service.Payment().SaveConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	// 记录操作日志
	service.Log().Record(ctx, userId, username, "支付设置", "修改支付配置", 1, "")

	return res, nil
}

// NewV1 创建支付控制器V1版本实例
func NewV1() *ControllerV1 {
	return &ControllerV1{}
}
