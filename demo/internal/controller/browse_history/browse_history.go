package browse_history

import (
	"demo/api/browse_history"
)

// Controller 浏览历史记录控制器
type Controller struct{}

// New 创建浏览历史记录控制器
func New() *Controller {
	return &Controller{}
}

// V1 创建v1版本API控制器
func (c *Controller) V1() browse_history.IBrowseHistoryV1 {
	return &ControllerV1{}
}

// ControllerV1 浏览历史记录v1版本控制器
type ControllerV1 struct{}
