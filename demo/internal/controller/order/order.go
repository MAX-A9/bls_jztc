package order

import (
	"demo/internal/middleware"

	"github.com/gogf/gf/v2/net/ghttp"
)

// OrderApi 订单管理API接口
type OrderApi interface {
	// Register 注册订单管理API
	Register(group *ghttp.RouterGroup)
}

// New 创建订单管理API实例
func New() OrderApi {
	return &controllerV1{}
}

// 私有结构体以控制实例化
type controllerV1 struct{}

// Register 注册订单管理API
func (c *controllerV1) Register(group *ghttp.RouterGroup) {
	// 直接使用 /order 路径，不在下级再创建 /order 组
	// 防止形成 /order/order/xxx 这样的重复路径
	routeGroup := group.Group("/order")

	// 添加认证中间件
	routeGroup.Middleware(middleware.Auth)

	// 添加响应处理中间件
	routeGroup.Middleware(ghttp.MiddlewareHandlerResponse)

	routeGroup.Bind(NewV1())
}
