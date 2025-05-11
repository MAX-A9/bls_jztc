package model

// PaymentConfig 支付配置模型
type PaymentConfig struct {
	// 微信支付配置
	AppId     string `json:"appId"`     // 微信支付AppID
	MchId     string `json:"mchId"`     // 微信支付商户号
	ApiKey    string `json:"apiKey"`    // API密钥
	NotifyUrl string `json:"notifyUrl"` // 回调通知地址
	IsEnabled bool   `json:"isEnabled"` // 是否启用微信支付
}
