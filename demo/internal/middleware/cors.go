package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

// CorsMiddleware CORS跨域请求中间件
func CorsMiddleware(r *ghttp.Request) {
	// 获取配置的允许域名列表
	allowedOrigins := g.Cfg().MustGet(r.GetCtx(), "security.cors.allowedOrigins").Strings()

	// 如果配置中没有设置允许的域名，默认允许本地开发环境
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"localhost", "127.0.0.1"}
	}

	// 获取当前请求的Origin
	origin := r.GetHeader("Origin")

	// 检查Origin是否在允许列表中
	allowed := false
	if origin != "" {
		// 从origin中提取主机部分（去掉协议部分）
		hostPart := origin
		if strings.Contains(origin, "://") {
			hostPart = strings.Split(origin, "://")[1]
		}

		for _, allowedOrigin := range allowedOrigins {
			// 完全匹配
			if allowedOrigin == "*" || allowedOrigin == origin || allowedOrigin == hostPart {
				allowed = true
				break
			}

			// 域名部分匹配（支持带端口号的情况）
			if strings.Contains(hostPart, allowedOrigin) {
				allowed = true
				break
			}

			// 通配符匹配
			if gstr.Contains(allowedOrigin, "*") {
				wildcardBase := gstr.Replace(allowedOrigin, "*", "")
				if gstr.Contains(hostPart, wildcardBase) {
					allowed = true
					break
				}
			}
		}

		if allowed {
			// 设置CORS头信息
			r.Response.Header().Set("Access-Control-Allow-Origin", origin)
			r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
			r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			r.Response.Header().Set("Access-Control-Max-Age", "86400") // 24小时
		} else {
			// 打印调试信息
			g.Log().Debugf(r.GetCtx(), "CORS拒绝访问: Origin=%s, 允许的域名=%v", origin, allowedOrigins)

			// 非法域名请求，返回403禁止访问
			r.Response.WriteStatusExit(403, g.Map{
				"code":    403,
				"message": "禁止访问: 域名未授权",
			})
			return
		}
	}

	// 处理预检请求
	if r.Method == "OPTIONS" {
		r.Response.WriteStatusExit(200, "")
		return
	}

	// 继续处理请求
	r.Middleware.Next()
}
