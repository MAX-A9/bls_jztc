package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"demo/utility/auth"
)

// ClientAuth 客户JWT认证中间件 - 保留旧接口，但使用新实现
func ClientAuth(r *ghttp.Request) {
	// 调用新的客户认证中间件
	auth.ClientAuthMiddleware(r)
}
