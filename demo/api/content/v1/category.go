package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// HomeCategoryListReq 首页分类列表请求
type HomeCategoryListReq struct {
	g.Meta `path:"/home-category/list" method:"get" tags:"内容管理" summary:"获取首页分类列表" security:"Bearer" description:"获取首页分类列表，需要管理员权限"`
}

// HomeCategoryListRes 首页分类列表响应
type HomeCategoryListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []CategoryItem `json:"list" dc:"分类列表"`
}

// IdleCategoryListReq 闲置分类列表请求
type IdleCategoryListReq struct {
	g.Meta `path:"/idle-category/list" method:"get" tags:"内容管理" summary:"获取闲置分类列表" security:"Bearer" description:"获取闲置分类列表，需要管理员权限"`
}

// IdleCategoryListRes 闲置分类列表响应
type IdleCategoryListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []CategoryItem `json:"list" dc:"分类列表"`
}

// CategoryItem 分类项
type CategoryItem struct {
	Id        int    `json:"id" dc:"分类ID"`
	Name      string `json:"name" dc:"分类名称"`
	SortOrder int    `json:"sortOrder" dc:"排序值"`
	IsActive  bool   `json:"isActive" dc:"是否启用"`
	Icon      string `json:"icon" dc:"分类图标URL"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
	UpdatedAt string `json:"updatedAt" dc:"更新时间"`
}

// HomeCategoryCreateReq 创建首页分类请求
type HomeCategoryCreateReq struct {
	g.Meta    `path:"/home-category/create" method:"post" tags:"内容管理" summary:"创建首页分类" security:"Bearer" description:"创建首页分类，需要管理员权限"`
	Name      string `v:"required#分类名称不能为空" json:"name" dc:"分类名称"`
	SortOrder int    `v:"min:0#排序值不能小于0" json:"sortOrder" dc:"排序值"`
	IsActive  bool   `json:"isActive" dc:"是否启用"`
	Icon      string `json:"icon" dc:"分类图标URL"`
}

// HomeCategoryCreateRes 创建首页分类响应
type HomeCategoryCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"分类ID"`
}

// IdleCategoryCreateReq 创建闲置分类请求
type IdleCategoryCreateReq struct {
	g.Meta    `path:"/idle-category/create" method:"post" tags:"内容管理" summary:"创建闲置分类" security:"Bearer" description:"创建闲置分类，需要管理员权限"`
	Name      string `v:"required#分类名称不能为空" json:"name" dc:"分类名称"`
	SortOrder int    `v:"min:0#排序值不能小于0" json:"sortOrder" dc:"排序值"`
	IsActive  bool   `json:"isActive" dc:"是否启用"`
	Icon      string `json:"icon" dc:"分类图标URL"`
}

// IdleCategoryCreateRes 创建闲置分类响应
type IdleCategoryCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"分类ID"`
}

// HomeCategoryUpdateReq 更新首页分类请求
type HomeCategoryUpdateReq struct {
	g.Meta    `path:"/home-category/update" method:"put" tags:"内容管理" summary:"更新首页分类" security:"Bearer" description:"更新首页分类，需要管理员权限"`
	Id        int    `v:"required#分类ID不能为空" json:"id" dc:"分类ID"`
	Name      string `v:"required#分类名称不能为空" json:"name" dc:"分类名称"`
	SortOrder int    `v:"min:0#排序值不能小于0" json:"sortOrder" dc:"排序值"`
	IsActive  bool   `json:"isActive" dc:"是否启用"`
	Icon      string `json:"icon" dc:"分类图标URL"`
}

// HomeCategoryUpdateRes 更新首页分类响应
type HomeCategoryUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// IdleCategoryUpdateReq 更新闲置分类请求
type IdleCategoryUpdateReq struct {
	g.Meta    `path:"/idle-category/update" method:"put" tags:"内容管理" summary:"更新闲置分类" security:"Bearer" description:"更新闲置分类，需要管理员权限"`
	Id        int    `v:"required#分类ID不能为空" json:"id" dc:"分类ID"`
	Name      string `v:"required#分类名称不能为空" json:"name" dc:"分类名称"`
	SortOrder int    `v:"min:0#排序值不能小于0" json:"sortOrder" dc:"排序值"`
	IsActive  bool   `json:"isActive" dc:"是否启用"`
	Icon      string `json:"icon" dc:"分类图标URL"`
}

// IdleCategoryUpdateRes 更新闲置分类响应
type IdleCategoryUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// HomeCategoryDeleteReq 删除首页分类请求
type HomeCategoryDeleteReq struct {
	g.Meta `path:"/home-category/delete" method:"delete" tags:"内容管理" summary:"删除首页分类" security:"Bearer" description:"删除首页分类，需要管理员权限"`
	Id     int `v:"required#分类ID不能为空" json:"id" dc:"分类ID"`
}

// HomeCategoryDeleteRes 删除首页分类响应
type HomeCategoryDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// IdleCategoryDeleteReq 删除闲置分类请求
type IdleCategoryDeleteReq struct {
	g.Meta `path:"/idle-category/delete" method:"delete" tags:"内容管理" summary:"删除闲置分类" security:"Bearer" description:"删除闲置分类，需要管理员权限"`
	Id     int `v:"required#分类ID不能为空" json:"id" dc:"分类ID"`
}

// IdleCategoryDeleteRes 删除闲置分类响应
type IdleCategoryDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}
