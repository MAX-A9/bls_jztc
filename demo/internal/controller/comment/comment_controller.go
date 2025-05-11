package comment

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"demo/internal/middleware"
)

// CommentApi 评论管理API接口
type CommentApi interface {
	// Register 注册评论管理API
	Register(group *ghttp.RouterGroup)
}

// New 创建评论管理API实例
func New() CommentApi {
	return &controllerV1{}
}

// 私有结构体以控制实例化
type controllerV1 struct{}

// V1 创建评论管理API实例V1版本
func (c *controllerV1) V1() *Controller {
	return &Controller{}
}

// Client 创建微信客户端评论控制器
func (c *controllerV1) Client() *ClientController {
	return &ClientController{}
}

// Register 注册评论管理API
func (c *controllerV1) Register(group *ghttp.RouterGroup) {
	// 创建评论路由组
	routeGroup := group.Group("/comment")

	// 添加认证中间件
	routeGroup.Middleware(middleware.Auth)

	// 添加响应处理中间件
	routeGroup.Middleware(ghttp.MiddlewareHandlerResponse)

	// 绑定评论管理接口
	routeGroup.Bind(
		c.V1().List,
		c.V1().Detail,
		c.V1().ContentComments,
		c.V1().Create,
		c.V1().Update,
		c.V1().Delete,
		c.V1().UpdateStatus,
	)

	// 注册微信客户端评论接口
	wxGroup := group.Group("/")
	// 评论列表接口不需要认证
	wxGroup.Bind(
		c.Client().WxClientList,
	)

	// 创建评论接口需要客户端认证
	wxAuthGroup := group.Group("/")
	wxAuthGroup.Middleware(middleware.ClientAuth)
	wxAuthGroup.Bind(
		c.Client().WxClientCreate,
	)
}
