package consts

// ContextKey 自定义上下文键类型，避免与其他包的 context key 冲突
type ContextKey string

// ConfigKeys 配置键常量
const (
	// 订单配置
	ConfigOrderExpireTime           = "order.expireTime"           // 订单过期时间(分钟)
	ConfigOrderDefaultPaymentMethod = "order.defaultPaymentMethod" // 订单默认支付方式
	ConfigOrderDefaultExpireDays    = "order.defaultExpireDays"    // 订单服务默认过期天数
)
