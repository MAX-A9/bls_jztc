package wechat

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"demo/internal/middleware"
)

// WechatPayApi 微信支付API接口
type WechatPayApi interface {
	// Register 注册微信支付API
	Register(group *ghttp.RouterGroup)
}

// New 创建微信支付API实例
func New() WechatPayApi {
	return &controllerV1{}
}

// 私有结构体以控制实例化
type controllerV1 struct{}

// Register 注册微信支付API
func (c *controllerV1) Register(group *ghttp.RouterGroup) {
	// 客户端接口路由，使用wx前缀
	wxGroup := group.Group("/wx/pay")

	// 需要认证的API
	authGroup := wxGroup.Group("/")
	// 添加客户端鉴权中间件
	authGroup.Middleware(middleware.ClientAuth)
	authGroup.Middleware(ghttp.MiddlewareHandlerResponse)

	// 注册需要认证的接口
	authGroup.POST("/unified-order", NewClientV1().UnifiedOrder)
	authGroup.POST("/query-order", NewClientV1().QueryOrder)

	// 不需要认证的API
	// 单独注册回调接口（不需要认证或响应中间件）
	wxGroup.POST("/notify", NewV1().Notify)
}
