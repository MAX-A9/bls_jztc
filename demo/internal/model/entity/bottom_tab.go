package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BottomTab 底部导航栏实体
type BottomTab struct {
	Id           int         `json:"id"            description:"ID" orm:"id,primary"`
	Name         string      `json:"name"          description:"Tab名称" orm:"name"`
	Icon         string      `json:"icon"          description:"未选中状态图标地址" orm:"icon"`
	SelectedIcon string      `json:"selected_icon" description:"选中状态图标地址" orm:"selected_icon"`
	Path         string      `json:"path"          description:"页面路径" orm:"path"`
	Order        int         `json:"order"         description:"排序值，越小越靠前" orm:"order"`
	IsEnabled    int         `json:"is_enabled"    description:"是否启用：0-禁用，1-启用" orm:"is_enabled"`
	CreatedAt    *gtime.Time `json:"created_at"    description:"创建时间" orm:"created_at"`
	UpdatedAt    *gtime.Time `json:"updated_at"    description:"更新时间" orm:"updated_at"`
}
