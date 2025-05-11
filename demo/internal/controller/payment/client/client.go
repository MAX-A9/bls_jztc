package client

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"demo/internal/middleware"
)

// ClientAPI 客户端API接口
type ClientAPI interface {
	// Register 注册客户端API
	Register(group *ghttp.RouterGroup)
}

// NewAPI 创建客户端API实例
func NewAPI() ClientAPI {
	return &controllerV1{}
}

// 私有结构体以控制实例化
type controllerV1 struct{}

// Register 注册客户端订单API
func (c *controllerV1) Register(group *ghttp.RouterGroup) {
	// 微信客户端接口组
	wxClientGroup := group.Group("/wx/client")

	// 需要客户端认证的接口
	wxClientGroup.Middleware(middleware.ClientAuth)
	wxClientGroup.Middleware(ghttp.MiddlewareHandlerResponse)

	// 注册客户端订单接口
	orderController := New()
	wxClientGroup.GET("/order/list", orderController.List)
	wxClientGroup.GET("/order/detail", orderController.Detail)
	wxClientGroup.POST("/order/pay", orderController.Pay)
	wxClientGroup.POST("/order/cancel", orderController.Cancel)
}

// Order 获取客户端订单控制器
func (c *controllerV1) Order() *OrderController {
	return New()
}
