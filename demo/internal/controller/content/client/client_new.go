package client

import (
	"demo/api/content/client"
)

// 创建控制器实例
func New() *Controller {
	return &Controller{}
}

// Controller 控制器结构体
type Controller struct{}

// V1 创建V1版本控制器
func (c *Controller) V1() client.IClientV1 {
	return &controllerV1{}
}
