package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// MiniProgramBaseSettingsReq 获取小程序基础设置请求
type MiniProgramBaseSettingsReq struct {
	g.Meta `path:"/mini-program/base/settings" method:"get" tags:"基础设置" summary:"获取小程序基础设置" security:"Bearer" description:"获取小程序基础设置信息，需要管理员权限"`
}

// MiniProgramBaseSettingsRes 获取小程序基础设置响应
type MiniProgramBaseSettingsRes struct {
	g.Meta      `mime:"application/json" example:"json"`
	Name        string `json:"name" dc:"小程序名称"`
	Description string `json:"description" dc:"小程序描述"`
	Logo        string `json:"logo" dc:"小程序Logo"`
}

// MiniProgramBaseSettingsSaveReq 保存小程序基础设置请求
type MiniProgramBaseSettingsSaveReq struct {
	g.Meta      `path:"/mini-program/base/settings/save" method:"post" tags:"基础设置" summary:"保存小程序基础设置" security:"Bearer" description:"保存小程序基础设置信息，需要管理员权限"`
	Name        string `v:"required#小程序名称不能为空" json:"name" dc:"小程序名称"`
	Description string `json:"description" dc:"小程序描述"`
	Logo        string `v:"required#小程序Logo不能为空" json:"logo" dc:"小程序Logo"`
}

// MiniProgramBaseSettingsSaveRes 保存小程序基础设置响应
type MiniProgramBaseSettingsSaveRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	IsSuccess bool `json:"isSuccess" dc:"是否成功"`
}

// WxMiniProgramBaseSettingsReq 微信客户端获取小程序基础设置请求
type WxMiniProgramBaseSettingsReq struct {
	g.Meta `path:"/wx/mini-program/base/settings" method:"get" tags:"客户端基础设置" summary:"获取小程序基础设置" description:"微信客户端获取小程序基础设置信息"`
}

// WxMiniProgramBaseSettingsRes 微信客户端获取小程序基础设置响应
type WxMiniProgramBaseSettingsRes struct {
	g.Meta      `mime:"application/json" example:"json"`
	Name        string `json:"name" dc:"小程序名称"`
	Description string `json:"description" dc:"小程序描述"`
	Logo        string `json:"logo" dc:"小程序Logo"`
}

// BannerGlobalStatusUpdateReq 轮播图总开关更新请求
type BannerGlobalStatusUpdateReq struct {
	g.Meta `path:"/banner/global/status" method:"post" tags:"总设置" summary:"更新轮播图总开关状态"`
	Enable *bool `json:"enable" v:"required#开关状态不能为空" dc:"是否启用"`
}

// BannerGlobalStatusUpdateRes 轮播图总开关更新响应
type BannerGlobalStatusUpdateRes struct {
	IsSuccess bool `json:"isSuccess" dc:"是否成功"`
}

// ActivityAreaGlobalStatusUpdateReq 活动区域总开关更新请求
type ActivityAreaGlobalStatusUpdateReq struct {
	g.Meta `path:"/activity-area/global/status" method:"post" tags:"总设置" summary:"更新活动区域总开关状态"`
	Enable *bool `json:"enable" v:"required#开关状态不能为空" dc:"是否启用"`
}

// ActivityAreaGlobalStatusUpdateRes 活动区域总开关更新响应
type ActivityAreaGlobalStatusUpdateRes struct {
	IsSuccess bool `json:"isSuccess" dc:"是否成功"`
}
