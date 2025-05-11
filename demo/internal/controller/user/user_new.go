package user

import (
	"demo/api/user"
)

// Controller 用户控制器
type Controller struct{}

// New 创建控制器实例
func New() *Controller {
	return &Controller{}
}

// V1 获取V1版本接口控制器
func (c *Controller) V1() user.IUserV1 {
	return &controllerV1{}
}
