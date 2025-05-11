package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ShopCategory 商城分类表
type ShopCategory struct {
	Id           int         `json:"id"           description:"分类ID"`
	Name         string      `json:"name"         description:"分类名称"`
	SortOrder    int         `json:"sortOrder"    description:"排序值"`
	ProductCount int         `json:"productCount" description:"商品数量"`
	Status       int         `json:"status"       description:"状态：1启用，0禁用"`
	Image        string      `json:"image"        description:"分类图片URL"`
	CreatedAt    *gtime.Time `json:"createdAt"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    description:"更新时间"`
}
