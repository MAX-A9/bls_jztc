package publisher

import (
	"demo/internal/service"
)

type Controller struct {
	service service.ClientService
}

func New() *Controller {
	return &Controller{
		service: service.Client(),
	}
}

// V1 返回v1版本的控制器
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{
		service: c.service,
	}
}
