package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MiniProgram 导航小程序表
type MiniProgram struct {
	Id        int         `json:"id"        description:"导航小程序ID"`
	Name      string      `json:"name"      description:"小程序名称"`
	AppId     string      `json:"appId"     description:"小程序AppID"`
	Logo      string      `json:"logo"      description:"小程序图标URL"`
	IsEnabled int         `json:"isEnabled" description:"是否启用 0:禁用 1:启用"`
	Order     int         `json:"order"     description:"排序值，数字越小排序越靠前"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" description:"更新时间"`
}

// Banner 轮播图表
type Banner struct {
	Id        int         `json:"id"        description:"轮播图ID"`
	Image     string      `json:"image"     description:"轮播图片URL"`
	LinkType  string      `json:"linkType"  description:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl   string      `json:"linkUrl"   description:"跳转地址"`
	IsEnabled int         `json:"isEnabled" description:"是否启用 0:禁用 1:启用"`
	Order     int         `json:"order"     description:"排序值，数字越小排序越靠前"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" description:"更新时间"`
}

// ActivityArea 活动区域表
type ActivityArea struct {
	Id                    int         `json:"id" description:"活动区域ID"`
	TopLeftTitle          string      `json:"topLeftTitle" description:"左上模块标题"`
	TopLeftDescription    string      `json:"topLeftDescription" description:"左上模块描述"`
	TopLeftLinkType       string      `json:"topLeftLinkType" description:"左上模块跳转类型"`
	TopLeftLinkUrl        string      `json:"topLeftLinkUrl" description:"左上模块跳转地址"`
	BottomLeftTitle       string      `json:"bottomLeftTitle" description:"左下模块标题"`
	BottomLeftDescription string      `json:"bottomLeftDescription" description:"左下模块描述"`
	BottomLeftLinkType    string      `json:"bottomLeftLinkType" description:"左下模块跳转类型"`
	BottomLeftLinkUrl     string      `json:"bottomLeftLinkUrl" description:"左下模块跳转地址"`
	RightTitle            string      `json:"rightTitle" description:"右侧模块标题"`
	RightDescription      string      `json:"rightDescription" description:"右侧模块描述"`
	RightLinkType         string      `json:"rightLinkType" description:"右侧模块跳转类型"`
	RightLinkUrl          string      `json:"rightLinkUrl" description:"右侧模块跳转地址"`
	CreatedAt             *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt             *gtime.Time `json:"updatedAt" description:"更新时间"`
}
