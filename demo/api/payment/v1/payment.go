package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PaymentConfigReq 获取支付配置请求
type PaymentConfigReq struct {
	g.Meta `path:"/config" method:"get" tags:"支付设置" summary:"获取支付配置" security:"Bearer" description:"获取支付配置，需要管理员权限"`
}

// PaymentConfigRes 获取支付配置响应
type PaymentConfigRes struct {
	g.Meta `mime:"application/json" example:"json"`
	// 微信支付配置
	AppId     string `json:"appId" dc:"微信支付AppID"`    // 微信支付AppID
	MchId     string `json:"mchId" dc:"微信支付商户号"`      // 微信支付商户号
	ApiKey    string `json:"apiKey" dc:"API密钥"`       // API密钥
	NotifyUrl string `json:"notifyUrl" dc:"回调通知地址"`   // 回调通知地址
	IsEnabled bool   `json:"isEnabled" dc:"是否启用微信支付"` // 是否启用微信支付
}

// SavePaymentConfigReq 保存支付配置请求
type SavePaymentConfigReq struct {
	g.Meta `path:"/config" method:"post" tags:"支付设置" summary:"保存支付配置" security:"Bearer" description:"保存支付配置，需要管理员权限"`
	// 微信支付配置
	AppId     string `v:"required#请输入微信支付AppID" json:"appId" dc:"微信支付AppID"` // 微信支付AppID
	MchId     string `v:"required#请输入微信支付商户号" json:"mchId" dc:"微信支付商户号"`     // 微信支付商户号
	ApiKey    string `v:"required#请输入API密钥" json:"apiKey" dc:"API密钥"`        // API密钥
	NotifyUrl string `v:"required#请输入回调通知地址" json:"notifyUrl" dc:"回调通知地址"`   // 回调通知地址
	IsEnabled bool   `json:"isEnabled" dc:"是否启用微信支付"`                        // 是否启用微信支付
}

// SavePaymentConfigRes 保存支付配置响应
type SavePaymentConfigRes struct {
	g.Meta `mime:"application/json" example:"json"`
	// 操作结果
	Success bool   `json:"success" dc:"是否成功"` // 是否成功
	Message string `json:"message" dc:"消息"`   // 消息
}
