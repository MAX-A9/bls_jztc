package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// InnerBannerItem 内页轮播图信息
type InnerBannerItem struct {
	Id         int    `json:"id" dc:"内页轮播图ID"`
	Image      string `json:"image" dc:"轮播图片地址"`
	LinkType   string `json:"linkType" dc:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl    string `json:"linkUrl" dc:"跳转地址"`
	IsEnabled  bool   `json:"isEnabled" dc:"是否启用"`
	Order      int    `json:"order" dc:"排序值，数字越小排序越靠前"`
	BannerType string `json:"bannerType" dc:"轮播图类型：home-首页轮播，idle-闲置轮播"`
}

// InnerBannerListReq 获取内页轮播图列表请求
type InnerBannerListReq struct {
	g.Meta     `path:"/inner-banner/list" method:"get" tags:"内页轮播图" summary:"获取内页轮播图列表" security:"Bearer" description:"获取内页轮播图列表，需要管理员权限"`
	BannerType string `v:"required|in:home,idle#轮播图类型不能为空|轮播图类型必须是home或idle" json:"bannerType" dc:"轮播图类型：home-首页轮播，idle-闲置轮播"`
}

// InnerBannerListRes 获取内页轮播图列表响应
type InnerBannerListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []InnerBannerItem `json:"list" dc:"内页轮播图列表"`
	IsGlobalEnabled bool              `json:"isGlobalEnabled" dc:"内页轮播图总开关状态"`
}

// InnerBannerCreateReq 创建内页轮播图请求
type InnerBannerCreateReq struct {
	g.Meta     `path:"/inner-banner/create" method:"post" tags:"内页轮播图" summary:"创建内页轮播图" security:"Bearer" description:"创建内页轮播图，需要管理员权限"`
	Image      string `v:"required#轮播图片不能为空" json:"image" dc:"轮播图片地址"`
	LinkType   string `v:"required#跳转类型不能为空" json:"linkType" dc:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl    string `v:"required#跳转地址不能为空" json:"linkUrl" dc:"跳转地址"`
	IsEnabled  bool   `json:"isEnabled" dc:"是否启用"`
	Order      int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，数字越小排序越靠前"`
	BannerType string `v:"required|in:home,idle#轮播图类型不能为空|轮播图类型必须是home或idle" json:"bannerType" dc:"轮播图类型：home-首页轮播，idle-闲置轮播"`
}

// InnerBannerCreateRes 创建内页轮播图响应
type InnerBannerCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"内页轮播图ID"`
}

// InnerBannerUpdateReq 更新内页轮播图请求
type InnerBannerUpdateReq struct {
	g.Meta     `path:"/inner-banner/update" method:"put" tags:"内页轮播图" summary:"更新内页轮播图" security:"Bearer" description:"更新内页轮播图，需要管理员权限"`
	Id         int    `v:"required#轮播图ID不能为空" json:"id" dc:"内页轮播图ID"`
	Image      string `v:"required#轮播图片不能为空" json:"image" dc:"轮播图片地址"`
	LinkType   string `v:"required#跳转类型不能为空" json:"linkType" dc:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl    string `v:"required#跳转地址不能为空" json:"linkUrl" dc:"跳转地址"`
	IsEnabled  bool   `json:"isEnabled" dc:"是否启用"`
	Order      int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，数字越小排序越靠前"`
	BannerType string `v:"required|in:home,idle#轮播图类型不能为空|轮播图类型必须是home或idle" json:"bannerType" dc:"轮播图类型：home-首页轮播，idle-闲置轮播"`
}

// InnerBannerUpdateRes 更新内页轮播图响应
type InnerBannerUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// InnerBannerDeleteReq 删除内页轮播图请求
type InnerBannerDeleteReq struct {
	g.Meta `path:"/inner-banner/delete" method:"delete" tags:"内页轮播图" summary:"删除内页轮播图" security:"Bearer" description:"删除内页轮播图，需要管理员权限"`
	Id     int `v:"required#轮播图ID不能为空" json:"id" dc:"内页轮播图ID"`
}

// InnerBannerDeleteRes 删除内页轮播图响应
type InnerBannerDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// InnerBannerStatusUpdateReq 更新内页轮播图状态请求
type InnerBannerStatusUpdateReq struct {
	g.Meta    `path:"/inner-banner/status/update" method:"put" tags:"内页轮播图" summary:"更新内页轮播图状态" security:"Bearer" description:"更新内页轮播图启用状态，需要管理员权限"`
	Id        int  `v:"required#轮播图ID不能为空" json:"id" dc:"内页轮播图ID"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用"`
}

// InnerBannerStatusUpdateRes 更新内页轮播图状态响应
type InnerBannerStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// HomeInnerBannerGlobalStatusUpdateReq 更新首页内页轮播图总开关请求
type HomeInnerBannerGlobalStatusUpdateReq struct {
	g.Meta    `path:"/inner-banner/home/global-status/update" method:"put" tags:"内页轮播图" summary:"更新首页内页轮播图总开关" security:"Bearer" description:"更新首页内页轮播图总开关状态，需要管理员权限"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用总开关"`
}

// HomeInnerBannerGlobalStatusUpdateRes 更新首页内页轮播图总开关响应
type HomeInnerBannerGlobalStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// IdleInnerBannerGlobalStatusUpdateReq 更新闲置页内页轮播图总开关请求
type IdleInnerBannerGlobalStatusUpdateReq struct {
	g.Meta    `path:"/inner-banner/idle/global-status/update" method:"put" tags:"内页轮播图" summary:"更新闲置页内页轮播图总开关" security:"Bearer" description:"更新闲置页内页轮播图总开关状态，需要管理员权限"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用总开关"`
}

// IdleInnerBannerGlobalStatusUpdateRes 更新闲置页内页轮播图总开关响应
type IdleInnerBannerGlobalStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// WxInnerBannerListReq 获取内页轮播图列表请求(客户端)
type WxInnerBannerListReq struct {
	g.Meta     `path:"/wx/inner-banner/list" method:"get" tags:"客户端内页轮播图" summary:"获取内页轮播图列表"`
	BannerType string `v:"required|in:home,idle#轮播图类型不能为空|轮播图类型必须是home或idle" json:"bannerType" dc:"轮播图类型：home-首页轮播，idle-闲置轮播"`
}

// WxInnerBannerListRes 获取内页轮播图列表响应(客户端)
type WxInnerBannerListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []InnerBannerItem `json:"list" dc:"内页轮播图列表"`
	IsGlobalEnabled bool              `json:"isGlobalEnabled" dc:"内页轮播图总开关状态"`
}
