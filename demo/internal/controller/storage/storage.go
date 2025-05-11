package storage

// New 创建存储控制器
func New() *Controller {
	return &Controller{}
}

// Controller 存储控制器
type Controller struct{}

// V1 获取存储控制器V1版本
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{
		Client: Client,
	}
}
