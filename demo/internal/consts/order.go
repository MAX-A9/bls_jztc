package consts

// 订单状态常量
const (
	// OrderStatusWaitPay 待支付
	OrderStatusWaitPay = 0
	// OrderStatusPaid 已支付
	OrderStatusPaid = 1
	// OrderStatusProcessing 进行中
	OrderStatusProcessing = 4
	// OrderStatusCompleted 已完成
	OrderStatusCompleted = 5
	// OrderStatusCancelled 已取消
	OrderStatusCancelled = 2
	// OrderStatusRefunded 已退款
	OrderStatusRefunded = 3
)
