package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 分享设置请求
type ShareSettingsReq struct {
	g.Meta `path:"/share/settings" method:"get" tags:"系统设置" summary:"获取分享设置" description:"获取系统分享相关设置"`
}

// 分享设置响应
type ShareSettingsRes struct {
	g.Meta            `mime:"application/json" example:"json"`
	DefaultShareText  string `json:"default_share_text" dc:"默认分享语"`
	DefaultShareImage string `json:"default_share_image" dc:"默认分享图片"`
	ContentShareText  string `json:"content_share_text" dc:"内容页分享语"`
	ContentShareImage string `json:"content_share_image" dc:"内容默认分享图片"`
	HomeShareText     string `json:"home_share_text" dc:"首页分享语"`
	HomeShareImage    string `json:"home_share_image" dc:"首页默认分享图片"`
}

// 保存分享设置请求
type SaveShareSettingsReq struct {
	g.Meta            `path:"/share/settings/save" method:"post" tags:"系统设置" summary:"保存分享设置" description:"保存系统分享相关设置"`
	DefaultShareText  string `v:"required#默认分享语不能为空" json:"default_share_text" dc:"默认分享语"`
	DefaultShareImage string `v:"required#默认分享图片不能为空" json:"default_share_image" dc:"默认分享图片"`
	ContentShareText  string `v:"required#内容页分享语不能为空" json:"content_share_text" dc:"内容页分享语"`
	ContentShareImage string `v:"required#内容默认分享图片不能为空" json:"content_share_image" dc:"内容默认分享图片"`
	HomeShareText     string `v:"required#首页分享语不能为空" json:"home_share_text" dc:"首页分享语"`
	HomeShareImage    string `v:"required#首页默认分享图片不能为空" json:"home_share_image" dc:"首页默认分享图片"`
}

// 保存分享设置响应
type SaveShareSettingsRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool `json:"success" dc:"是否保存成功"`
}

// 微信客户端获取分享设置请求
type WxShareSettingsReq struct {
	g.Meta `path:"/wx/share/settings" method:"get" tags:"系统设置" summary:"微信客户端获取分享设置" description:"微信客户端获取系统分享相关设置"`
}

// 微信客户端获取分享设置响应
type WxShareSettingsRes struct {
	g.Meta            `mime:"application/json" example:"json"`
	DefaultShareText  string `json:"default_share_text" dc:"默认分享语"`
	DefaultShareImage string `json:"default_share_image" dc:"默认分享图片"`
	ContentShareText  string `json:"content_share_text" dc:"内容页分享语"`
	ContentShareImage string `json:"content_share_image" dc:"内容默认分享图片"`
	HomeShareText     string `json:"home_share_text" dc:"首页分享语"`
	HomeShareImage    string `json:"home_share_image" dc:"首页默认分享图片"`
}
