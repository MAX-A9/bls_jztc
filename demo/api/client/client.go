package client

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// ClientApi 客户管理API接口
type ClientApi interface {
	// Register 注册客户管理API
	Register(group *ghttp.RouterGroup)
}

// New 创建客户管理API实例
func New() ClientApi {
	return &controllerV1{}
}

// 私有结构体以控制实例化
type controllerV1 struct{}

// Register 注册客户管理API
func (c *controllerV1) Register(group *ghttp.RouterGroup) {
	group.Group("/client", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		// 客户管理接口
		group.Bind(c)
	})
}
