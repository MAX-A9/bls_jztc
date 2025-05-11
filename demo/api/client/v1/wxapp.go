package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// WxappConfigGetReq 获取微信小程序配置请求
type WxappConfigGetReq struct {
	g.Meta `path:"/config" method:"get" tags:"小程序管理" summary:"获取微信小程序配置" security:"Bearer" description:"获取微信小程序配置信息，需要管理员权限"`
}

// WxappConfigGetRes 获取微信小程序配置响应
type WxappConfigGetRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	AppId     string `json:"appId" dc:"小程序AppID"`
	AppSecret string `json:"appSecret" dc:"小程序AppSecret"`
	Enabled   bool   `json:"enabled" dc:"是否启用微信小程序功能"`
}

// WxappConfigSaveReq 保存微信小程序配置请求
type WxappConfigSaveReq struct {
	g.Meta    `path:"/config" method:"post" tags:"小程序管理" summary:"保存微信小程序配置" security:"Bearer" description:"保存微信小程序配置信息，需要管理员权限"`
	AppId     string `v:"required#小程序AppID不能为空" json:"appId" dc:"小程序AppID"`
	AppSecret string `v:"required#小程序AppSecret不能为空" json:"appSecret" dc:"小程序AppSecret"`
	Enabled   bool   `json:"enabled" dc:"是否启用微信小程序功能"`
}

// WxappConfigSaveRes 保存微信小程序配置响应
type WxappConfigSaveRes struct {
	g.Meta `mime:"application/json" example:"json"`
}
