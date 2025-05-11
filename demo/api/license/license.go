// Package license 提供授权验证相关API定义
package license

import (
	"github.com/gogf/gf/v2/frame/g"
)

// VerifyReq 授权验证请求
type VerifyReq struct {
	g.Meta      `path:"/verify" method:"post" tags:"License" summary:"验证授权"`
	ClientIP    string `json:"client_ip" dc:"客户端IP地址(可选,默认使用服务器IP)"`
	Domain      string `v:"required" json:"domain" dc:"域名"`
	LicenseCode string `v:"required" json:"license_code" dc:"授权码"`
}

// VerifyRes 授权验证响应
type VerifyRes struct {
	Valid bool `json:"valid" dc:"授权是否有效"`
}

// StatusReq 获取授权状态请求
type StatusReq struct {
	g.Meta `path:"/status" method:"get" tags:"License" summary:"获取授权状态"`
	Domain string `v:"required" json:"domain" dc:"域名"`
}

// StatusRes 获取授权状态响应
type StatusRes struct {
	Valid bool `json:"valid" dc:"授权是否有效"`
}
