// Package middleware 提供中间件
package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"

	"demo/utility/license"
)

// 授权验证白名单路径
var licenseAuthWhitelist = []string{
	"/license/verify", // 授权验证接口
	"/license/status", // 授权状态查询接口
	"/swagger",        // Swagger文档
	"/api.json",       // OpenAPI定义
	"/resource",       // 静态资源
}

// LicenseAuth 授权验证中间件
func LicenseAuth(r *ghttp.Request) {
	// 获取请求域名
	host := r.GetHost()
	domain := license.ExtractDomainFromHost(host)

	// 检查是否为本地环境
	if license.IsLocalhost(domain) {
		r.Middleware.Next()
		return
	}

	// 检查当前路径是否在白名单中
	path := r.URL.Path
	if isWhitelistPath(path) {
		r.Middleware.Next()
		return
	}

	// 检查域名授权状态
	if !license.GetLicenseStatus(domain) {
		glog.Warning(r.GetCtx(), "未授权访问:", path, "域名:", domain)
		r.Response.WriteJson(g.Map{
			"code":    403,
			"message": "未授权访问，请先获取授权",
		})
		r.Exit()
		return
	}

	// 通过授权验证，继续处理请求
	r.Middleware.Next()
}

// 检查是否为白名单路径
func isWhitelistPath(path string) bool {
	for _, whitelistPath := range licenseAuthWhitelist {
		if gstr.HasPrefix(path, whitelistPath) {
			return true
		}
	}
	return false
}
