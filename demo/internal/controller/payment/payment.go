package payment

// New 创建支付控制器
func New() *Controller {
	return &Controller{}
}

// Controller 支付控制器
type Controller struct{}

// V1 获取支付控制器V1版本
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{}
}
