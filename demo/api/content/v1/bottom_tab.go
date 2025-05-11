package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// BottomTabItem 底部导航项信息
type BottomTabItem struct {
	Id           int    `json:"id" dc:"导航项ID"`
	Name         string `json:"name" dc:"Tab名称"`
	Icon         string `json:"icon" dc:"未选中状态图标地址"`
	SelectedIcon string `json:"selectedIcon" dc:"选中状态图标地址"`
	Path         string `json:"path" dc:"页面路径"`
	Order        int    `json:"order" dc:"排序值，越小越靠前"`
	IsEnabled    bool   `json:"isEnabled" dc:"是否启用"`
}

// BottomTabListReq 获取底部导航栏列表请求
type BottomTabListReq struct {
	g.Meta `path:"/bottom-tab/list" method:"get" tags:"底部导航栏" summary:"获取底部导航栏列表" security:"Bearer" description:"获取底部导航栏列表，需要管理员权限"`
}

// BottomTabListRes 获取底部导航栏列表响应
type BottomTabListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []BottomTabItem `json:"list" dc:"底部导航栏列表"`
}

// BottomTabCreateReq 创建底部导航项请求
type BottomTabCreateReq struct {
	g.Meta       `path:"/bottom-tab/create" method:"post" tags:"底部导航栏" summary:"创建底部导航项" security:"Bearer" description:"创建底部导航项，需要管理员权限"`
	Name         string `v:"required#Tab名称不能为空" json:"name" dc:"Tab名称"`
	Icon         string `v:"required#未选中状态图标地址不能为空" json:"icon" dc:"未选中状态图标地址"`
	SelectedIcon string `v:"required#选中状态图标地址不能为空" json:"selectedIcon" dc:"选中状态图标地址"`
	Path         string `v:"required#页面路径不能为空" json:"path" dc:"页面路径"`
	Order        int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，越小越靠前"`
	IsEnabled    bool   `json:"isEnabled" dc:"是否启用"`
}

// BottomTabCreateRes 创建底部导航项响应
type BottomTabCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"导航项ID"`
}

// BottomTabUpdateReq 更新底部导航项请求
type BottomTabUpdateReq struct {
	g.Meta       `path:"/bottom-tab/update" method:"put" tags:"底部导航栏" summary:"更新底部导航项" security:"Bearer" description:"更新底部导航项，需要管理员权限"`
	Id           int    `v:"required#导航项ID不能为空" json:"id" dc:"导航项ID"`
	Name         string `v:"required#Tab名称不能为空" json:"name" dc:"Tab名称"`
	Icon         string `v:"required#未选中状态图标地址不能为空" json:"icon" dc:"未选中状态图标地址"`
	SelectedIcon string `v:"required#选中状态图标地址不能为空" json:"selectedIcon" dc:"选中状态图标地址"`
	Path         string `v:"required#页面路径不能为空" json:"path" dc:"页面路径"`
	Order        int    `v:"min:0#排序值不能小于0" json:"order" dc:"排序值，越小越靠前"`
	IsEnabled    bool   `json:"isEnabled" dc:"是否启用"`
}

// BottomTabUpdateRes 更新底部导航项响应
type BottomTabUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// BottomTabDeleteReq 删除底部导航项请求
type BottomTabDeleteReq struct {
	g.Meta `path:"/bottom-tab/delete" method:"delete" tags:"底部导航栏" summary:"删除底部导航项" security:"Bearer" description:"删除底部导航项，需要管理员权限"`
	Id     int `v:"required#导航项ID不能为空" json:"id" dc:"导航项ID"`
}

// BottomTabDeleteRes 删除底部导航项响应
type BottomTabDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// BottomTabStatusUpdateReq 更新底部导航项状态请求
type BottomTabStatusUpdateReq struct {
	g.Meta    `path:"/bottom-tab/status/update" method:"put" tags:"底部导航栏" summary:"更新底部导航项状态" security:"Bearer" description:"更新底部导航项启用状态，需要管理员权限"`
	Id        int  `v:"required#导航项ID不能为空" json:"id" dc:"导航项ID"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用"`
}

// BottomTabStatusUpdateRes 更新底部导航项状态响应
type BottomTabStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// WxBottomTabListReq 获取底部导航栏列表请求(客户端)
type WxBottomTabListReq struct {
	g.Meta `path:"/wx/bottom-tab/list" method:"get" tags:"客户端底部导航栏" summary:"获取底部导航栏列表"`
}

// WxBottomTabListRes 获取底部导航栏列表响应(客户端)
type WxBottomTabListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []BottomTabItem `json:"list" dc:"底部导航栏列表"`
}
