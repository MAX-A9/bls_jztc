package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// MiniProgramItem 导航小程序信息
type MiniProgramItem struct {
	Id        int    `json:"id" dc:"导航小程序ID"`
	Name      string `json:"name" dc:"小程序名称"`
	AppId     string `json:"appId" dc:"小程序AppID"`
	Logo      string `json:"logo" dc:"小程序图标"`
	IsEnabled bool   `json:"isEnabled" dc:"是否启用"`
	Order     int    `json:"order" dc:"排序值，数字越小排序越靠前"`
}

// BannerItem 轮播图信息
type BannerItem struct {
	Id        int    `json:"id" dc:"轮播图ID"`
	Image     string `json:"image" dc:"轮播图片地址"`
	LinkType  string `json:"linkType" dc:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl   string `json:"linkUrl" dc:"跳转地址"`
	IsEnabled bool   `json:"isEnabled" dc:"是否启用"`
	Order     int    `json:"order" dc:"排序值，数字越小排序越靠前"`
}

// MiniProgramListReq 获取导航小程序列表请求
type MiniProgramListReq struct {
	g.Meta `path:"/mini-program/list" method:"get" tags:"首页布局" summary:"获取导航小程序列表" security:"Bearer" description:"获取导航小程序列表，需要管理员权限"`
}

// MiniProgramListRes 获取导航小程序列表响应
type MiniProgramListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []MiniProgramItem `json:"list" dc:"导航小程序列表"`
	IsGlobalEnabled bool              `json:"isGlobalEnabled" dc:"导航小程序总开关状态"`
}

// MiniProgramCreateReq 创建导航小程序请求
type MiniProgramCreateReq struct {
	g.Meta    `path:"/mini-program/create" method:"post" tags:"首页布局" summary:"创建导航小程序" security:"Bearer" description:"创建导航小程序，需要管理员权限"`
	Name      string `v:"required#小程序名称不能为空" json:"name" dc:"小程序名称"`
	AppId     string `v:"required#小程序AppID不能为空" json:"appId" dc:"小程序AppID"`
	Logo      string `json:"logo" dc:"小程序图标"`
	IsEnabled bool   `json:"isEnabled" dc:"是否启用"`
	Order     int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，数字越小排序越靠前"`
}

// MiniProgramCreateRes 创建导航小程序响应
type MiniProgramCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"导航小程序ID"`
}

// MiniProgramUpdateReq 更新导航小程序请求
type MiniProgramUpdateReq struct {
	g.Meta    `path:"/mini-program/update" method:"put" tags:"首页布局" summary:"更新导航小程序" security:"Bearer" description:"更新导航小程序，需要管理员权限"`
	Id        int    `v:"required#小程序ID不能为空" json:"id" dc:"导航小程序ID"`
	Name      string `v:"required#小程序名称不能为空" json:"name" dc:"小程序名称"`
	AppId     string `v:"required#小程序AppID不能为空" json:"appId" dc:"小程序AppID"`
	Logo      string `json:"logo" dc:"小程序图标"`
	IsEnabled bool   `json:"isEnabled" dc:"是否启用"`
	Order     int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，数字越小排序越靠前"`
}

// MiniProgramUpdateRes 更新导航小程序响应
type MiniProgramUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// MiniProgramDeleteReq 删除导航小程序请求
type MiniProgramDeleteReq struct {
	g.Meta `path:"/mini-program/delete" method:"delete" tags:"首页布局" summary:"删除导航小程序" security:"Bearer" description:"删除导航小程序，需要管理员权限"`
	Id     int `v:"required#小程序ID不能为空" json:"id" dc:"导航小程序ID"`
}

// MiniProgramDeleteRes 删除导航小程序响应
type MiniProgramDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// MiniProgramStatusUpdateReq 更新导航小程序状态请求
type MiniProgramStatusUpdateReq struct {
	g.Meta    `path:"/mini-program/status/update" method:"put" tags:"首页布局" summary:"更新导航小程序状态" security:"Bearer" description:"更新导航小程序启用状态，需要管理员权限"`
	Id        int  `v:"required#小程序ID不能为空" json:"id" dc:"导航小程序ID"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用"`
}

// MiniProgramStatusUpdateRes 更新导航小程序状态响应
type MiniProgramStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// MiniProgramGlobalStatusUpdateReq 更新导航小程序总开关请求
type MiniProgramGlobalStatusUpdateReq struct {
	g.Meta    `path:"/mini-program/global-status/update" method:"put" tags:"首页布局" summary:"更新导航小程序总开关" security:"Bearer" description:"更新导航小程序总开关状态，需要管理员权限"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用总开关"`
}

// MiniProgramGlobalStatusUpdateRes 更新导航小程序总开关响应
type MiniProgramGlobalStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// BannerListReq 获取轮播图列表请求
type BannerListReq struct {
	g.Meta `path:"/banner/list" method:"get" tags:"首页布局" summary:"获取轮播图列表" security:"Bearer" description:"获取轮播图列表，需要管理员权限"`
}

// BannerListRes 获取轮播图列表响应
type BannerListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []BannerItem `json:"list" dc:"轮播图列表"`
	IsGlobalEnabled bool         `json:"isGlobalEnabled" dc:"轮播图总开关状态"`
}

// BannerCreateReq 创建轮播图请求
type BannerCreateReq struct {
	g.Meta    `path:"/banner/create" method:"post" tags:"首页布局" summary:"创建轮播图" security:"Bearer" description:"创建轮播图，需要管理员权限"`
	Image     string `v:"required#轮播图片不能为空" json:"image" dc:"轮播图片地址"`
	LinkType  string `v:"required#跳转类型不能为空" json:"linkType" dc:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl   string `v:"required#跳转地址不能为空" json:"linkUrl" dc:"跳转地址"`
	IsEnabled bool   `json:"isEnabled" dc:"是否启用"`
	Order     int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，数字越小排序越靠前"`
}

// BannerCreateRes 创建轮播图响应
type BannerCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"轮播图ID"`
}

// BannerUpdateReq 更新轮播图请求
type BannerUpdateReq struct {
	g.Meta    `path:"/banner/update" method:"put" tags:"首页布局" summary:"更新轮播图" security:"Bearer" description:"更新轮播图，需要管理员权限"`
	Id        int    `v:"required#轮播图ID不能为空" json:"id" dc:"轮播图ID"`
	Image     string `v:"required#轮播图片不能为空" json:"image" dc:"轮播图片地址"`
	LinkType  string `v:"required#跳转类型不能为空" json:"linkType" dc:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl   string `v:"required#跳转地址不能为空" json:"linkUrl" dc:"跳转地址"`
	IsEnabled bool   `json:"isEnabled" dc:"是否启用"`
	Order     int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，数字越小排序越靠前"`
}

// BannerUpdateRes 更新轮播图响应
type BannerUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// BannerDeleteReq 删除轮播图请求
type BannerDeleteReq struct {
	g.Meta `path:"/banner/delete" method:"delete" tags:"首页布局" summary:"删除轮播图" security:"Bearer" description:"删除轮播图，需要管理员权限"`
	Id     int `v:"required#轮播图ID不能为空" json:"id" dc:"轮播图ID"`
}

// BannerDeleteRes 删除轮播图响应
type BannerDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// BannerStatusUpdateReq 更新轮播图状态请求
type BannerStatusUpdateReq struct {
	g.Meta    `path:"/banner/status/update" method:"put" tags:"首页布局" summary:"更新轮播图状态" security:"Bearer" description:"更新轮播图启用状态，需要管理员权限"`
	Id        int  `v:"required#轮播图ID不能为空" json:"id" dc:"轮播图ID"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用"`
}

// BannerStatusUpdateRes 更新轮播图状态响应
type BannerStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// BannerGlobalStatusUpdateReq 更新轮播图总开关请求
type BannerGlobalStatusUpdateReq struct {
	g.Meta    `path:"/banner/global-status/update" method:"put" tags:"首页布局" summary:"更新轮播图总开关" security:"Bearer" description:"更新轮播图总开关状态，需要管理员权限"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用总开关"`
}

// BannerGlobalStatusUpdateRes 更新轮播图总开关响应
type BannerGlobalStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ActivityAreaModule 活动区域模块信息
type ActivityAreaModule struct {
	Title       string `json:"title" dc:"模块标题"`
	Description string `json:"description" dc:"模块描述"`
	LinkType    string `json:"linkType" dc:"跳转类型：page-小程序页面, miniprogram-其他小程序, webview-网页"`
	LinkUrl     string `json:"linkUrl" dc:"跳转地址"`
	Position    string `json:"position" dc:"位置：topLeft-左上, bottomLeft-左下, right-右侧"`
}

// ActivityAreaSaveReq 保存活动区域请求
type ActivityAreaSaveReq struct {
	g.Meta     `path:"/activity-area/save" method:"post" tags:"首页布局" summary:"保存活动区域" security:"Bearer" description:"保存首页活动区域设置，需要管理员权限"`
	TopLeft    ActivityAreaModule `v:"required#左上模块不能为空" json:"topLeft" dc:"左上模块"`
	BottomLeft ActivityAreaModule `v:"required#左下模块不能为空" json:"bottomLeft" dc:"左下模块"`
	Right      ActivityAreaModule `v:"required#右侧模块不能为空" json:"right" dc:"右侧模块"`
}

// ActivityAreaSaveRes 保存活动区域响应
type ActivityAreaSaveRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ActivityAreaGetReq 获取活动区域请求
type ActivityAreaGetReq struct {
	g.Meta `path:"/activity-area/get" method:"get" tags:"首页布局" summary:"获取活动区域" security:"Bearer" description:"获取首页活动区域设置"`
}

// ActivityAreaGetRes 获取活动区域响应
type ActivityAreaGetRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	TopLeft         ActivityAreaModule `json:"topLeft" dc:"左上模块"`
	BottomLeft      ActivityAreaModule `json:"bottomLeft" dc:"左下模块"`
	Right           ActivityAreaModule `json:"right" dc:"右侧模块"`
	IsGlobalEnabled bool               `json:"isGlobalEnabled" dc:"活动区域总开关状态"`
}

// WxActivityAreaGetReq 微信客户端获取活动区域请求
type WxActivityAreaGetReq struct {
	g.Meta `path:"/wx/activity-area/get" method:"get" tags:"客户端首页布局" summary:"获取活动区域"`
}

// WxActivityAreaGetRes 微信客户端获取活动区域响应
type WxActivityAreaGetRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []ActivityAreaModule `json:"list" dc:"活动区域模块列表"`
	IsGlobalEnabled bool                 `json:"isGlobalEnabled" dc:"活动区域总开关状态"`
}

// ActivityAreaGlobalStatusUpdateReq 更新活动区域总开关请求
type ActivityAreaGlobalStatusUpdateReq struct {
	g.Meta    `path:"/activity-area/global-status/update" method:"put" tags:"首页布局" summary:"更新活动区域总开关" security:"Bearer" description:"更新活动区域总开关状态，需要管理员权限"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用总开关"`
}

// ActivityAreaGlobalStatusUpdateRes 更新活动区域总开关响应
type ActivityAreaGlobalStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// WxMiniProgramListReq 获取导航小程序列表请求(客户端)
type WxMiniProgramListReq struct {
	g.Meta `path:"/wx/mini-program/list" method:"get" tags:"客户端首页布局" summary:"获取导航小程序列表"`
}

// WxMiniProgramListRes 获取导航小程序列表响应(客户端)
type WxMiniProgramListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []MiniProgramItem `json:"list" dc:"导航小程序列表"`
	IsGlobalEnabled bool              `json:"isGlobalEnabled" dc:"导航小程序总开关状态"`
}

// WxBannerListReq 获取轮播图列表请求(客户端)
type WxBannerListReq struct {
	g.Meta `path:"/wx/banner/list" method:"get" tags:"客户端首页布局" summary:"获取轮播图列表"`
}

// WxBannerListRes 获取轮播图列表响应(客户端)
type WxBannerListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []BannerItem `json:"list" dc:"轮播图列表"`
	IsGlobalEnabled bool         `json:"isGlobalEnabled" dc:"轮播图总开关状态"`
}
