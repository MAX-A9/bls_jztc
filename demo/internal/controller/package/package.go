package package_controller

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"demo/internal/middleware"
)

// 套餐控制器实例
var (
	Package = controllerNew()
)

// controllerNew 创建套餐控制器
func controllerNew() *Controller {
	return &Controller{}
}

// Register 注册路由
func (c *Controller) Register(group *ghttp.RouterGroup) {
	// 需要权限验证的路由组
	packageGroup := group.Group("/")
	packageGroup.Middleware(middleware.Auth)

	// 绑定套餐接口
	packageGroup.Bind(
		c.V1().List,
		c.V1().Detail,
		c.V1().Create,
		c.V1().Update,
		c.V1().Delete,
		c.V1().GetGlobalStatus,
		c.V1().UpdateTopPackageGlobalStatus,
		c.V1().UpdatePublishPackageGlobalStatus,
	)

	// 绑定客户端接口(不需要验证)
	group.Bind(
		c.V1().WxList,
	)
}
