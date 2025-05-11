package middleware

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"

	"demo/utility/auth"
)

// Auth 通用JWT认证中间件 - 保留旧接口，但使用新实现
func Auth(r *ghttp.Request) {
	// 调用新的管理员认证中间件
	auth.AdminAuthMiddleware(r)
}

// GetCurrentUser 从上下文中获取当前用户ID - 保留旧接口，但使用新实现
func GetCurrentUser(ctx context.Context) (userId int, username string, role string, err error) {
	// 调用auth包中的GetUserInfo方法
	return auth.GetUserInfo(ctx)
}

// 保留ContextKey类型和常量，保持向后兼容
type ContextKey string

// 导出上下文键常量
const (
	RoleKey     ContextKey = "role"
	UserIdKey   ContextKey = "userId"
	UsernameKey ContextKey = "username"
)
