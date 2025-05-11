package client

import (
	"demo/internal/service"
)

// Controller 客户管理控制器
type Controller struct{}

// New 创建客户管理控制器
func New() *Controller {
	return &Controller{}
}

// V1 创建v1版本API控制器
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{
		service: service.Client(),
	}
}

// ControllerV1 客户管理v1版本控制器
type ControllerV1 struct {
	service service.ClientService
}

// NewReward 创建奖励记录控制器
func (c *ControllerV1) NewReward() *Reward {
	return NewReward()
}
