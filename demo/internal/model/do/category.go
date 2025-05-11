package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// HomeCategory 首页分类表相关常量
const (
	// TableHomeCategory 表名
	TableHomeCategory = "home_category"
	// HomeCategoryColumns 所有字段
	HomeCategoryColumns = "id,name,sort_order,is_active,icon,created_at,updated_at,deleted_at"
)

// IdleCategory 闲置分类表相关常量
const (
	// TableIdleCategory 表名
	TableIdleCategory = "idle_category"
	// IdleCategoryColumns 所有字段
	IdleCategoryColumns = "id,name,sort_order,is_active,icon,created_at,updated_at,deleted_at"
)

// CategoryDO 分类DO基础结构体
type CategoryDO struct {
	Id        interface{} `orm:"id,primary" json:"id"`        // 分类ID
	Name      interface{} `orm:"name" json:"name"`            // 分类名称
	SortOrder interface{} `orm:"sort_order" json:"sortOrder"` // 排序值
	IsActive  interface{} `orm:"is_active" json:"isActive"`   // 是否启用
	Icon      interface{} `orm:"icon" json:"icon"`            // 分类图标
	CreatedAt *gtime.Time `orm:"created_at" json:"createdAt"` // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at" json:"updatedAt"` // 更新时间
	DeletedAt *gtime.Time `orm:"deleted_at" json:"deletedAt"` // 删除时间
}

// HomeCategoryDO 首页分类DO结构体
type HomeCategoryDO struct {
	g.Meta `orm:"table:home_category, do:true"`
	CategoryDO
}

// IdleCategoryDO 闲置分类DO结构体
type IdleCategoryDO struct {
	g.Meta `orm:"table:idle_category, do:true"`
	CategoryDO
}
