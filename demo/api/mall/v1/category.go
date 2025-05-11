package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShopCategoryGetReq 获取商城分类详情请求
type ShopCategoryGetReq struct {
	g.Meta `path:"/shop-category/get" method:"get" tags:"商城分类" summary:"获取商城分类详情" security:"Bearer" description:"获取指定ID的商城分类详情"`
	Id     int `v:"required#商城分类ID不能为空" dc:"商城分类ID" json:"id" query:"id"`
}

// ShopCategoryGetRes 获取商城分类详情响应
type ShopCategoryGetRes struct {
	g.Meta       `mime:"application/json" example:"json"`
	Id           int         `json:"id" dc:"分类ID"`
	Name         string      `json:"name" dc:"分类名称"`
	SortOrder    int         `json:"sortOrder" dc:"排序值"`
	ProductCount int         `json:"productCount" dc:"商品数量"`
	Status       int         `json:"status" dc:"状态：1启用，0禁用"`
	Image        string      `json:"image" dc:"分类图片URL"`
	CreatedAt    *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// ShopCategoryCreateReq 创建商城分类请求
type ShopCategoryCreateReq struct {
	g.Meta       `path:"/shop-category/create" method:"post" tags:"商城分类" summary:"创建商城分类" security:"Bearer" description:"创建新的商城分类"`
	Name         string `v:"required#分类名称不能为空" json:"name" dc:"分类名称"`
	SortOrder    int    `d:"0" json:"sortOrder" dc:"排序值，默认0"`
	ProductCount int    `d:"0" json:"productCount" dc:"商品数量，默认0"`
	Status       int    `d:"1" json:"status" dc:"状态：1启用，0禁用，默认1"`
	Image        string `json:"image" dc:"分类图片URL"`
}

// ShopCategoryCreateRes 创建商城分类响应
type ShopCategoryCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"新创建的商城分类ID"`
}

// ShopCategoryUpdateReq 更新商城分类请求
type ShopCategoryUpdateReq struct {
	g.Meta       `path:"/shop-category/update" method:"put" tags:"商城分类" summary:"更新商城分类" security:"Bearer" description:"更新指定ID的商城分类信息"`
	Id           int    `v:"required#商城分类ID不能为空" json:"id" dc:"商城分类ID"`
	Name         string `v:"required#分类名称不能为空" json:"name" dc:"分类名称"`
	SortOrder    int    `json:"sortOrder" dc:"排序值"`
	ProductCount int    `json:"productCount" dc:"商品数量"`
	Status       int    `json:"status" dc:"状态：1启用，0禁用"`
	Image        string `json:"image" dc:"分类图片URL"`
}

// ShopCategoryUpdateRes 更新商城分类响应
type ShopCategoryUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ShopCategoryDeleteReq 删除商城分类请求
type ShopCategoryDeleteReq struct {
	g.Meta `path:"/shop-category/delete" method:"delete" tags:"商城分类" summary:"删除商城分类" security:"Bearer" description:"删除指定ID的商城分类"`
	Id     int `v:"required#商城分类ID不能为空" json:"id" dc:"商城分类ID"`
}

// ShopCategoryDeleteRes 删除商城分类响应
type ShopCategoryDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ShopCategoryStatusUpdateReq 更新商城分类状态请求
type ShopCategoryStatusUpdateReq struct {
	g.Meta `path:"/shop-category/status/update" method:"put" tags:"商城分类" summary:"更新商城分类状态" security:"Bearer" description:"更新指定ID的商城分类状态"`
	Id     int `v:"required#商城分类ID不能为空" json:"id" dc:"商城分类ID"`
	Status int `v:"required#状态不能为空" json:"status" dc:"状态：1启用，0禁用"`
}

// ShopCategoryStatusUpdateRes 更新商城分类状态响应
type ShopCategoryStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ShopCategoryListReq 获取商城分类列表请求
type ShopCategoryListReq struct {
	g.Meta    `path:"/shop-category/list" method:"get" tags:"商城分类" summary:"获取商城分类列表" security:"Bearer" description:"获取商城分类列表，支持全量获取和分页查询"`
	Page      int    `v:"min:0#页码最小值为0"  dc:"页码，为0时表示不分页获取全部" json:"page" d:"0"`
	Size      int    `v:"max:50#每页最大50条" dc:"每页数量，当page>0时生效" json:"size" d:"10"`
	Name      string `dc:"分类名称, 支持模糊搜索" json:"name"`
	Status    int    `dc:"状态 0:禁用 1:启用" json:"status" d:"-1"`
	SortField string `dc:"排序字段: id, name, sort_order, product_count, created_at" json:"sortField" d:"sort_order"`
	SortOrder string `dc:"排序方式: asc, desc" json:"sortOrder" d:"asc"`
}

// ShopCategoryListRes 获取商城分类列表响应
type ShopCategoryListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ShopCategoryItem `json:"list" dc:"商城分类列表"`
	Total  int                `json:"total,omitempty" dc:"总记录数(分页时返回)"`
	Page   int                `json:"page,omitempty" dc:"当前页码(分页时返回)"`
	Size   int                `json:"size,omitempty" dc:"每页记录数(分页时返回)"`
	Pages  int                `json:"pages,omitempty" dc:"总页数(分页时返回)"`
}

// ShopCategoryItem 商城分类项
type ShopCategoryItem struct {
	Id           int         `json:"id" dc:"分类ID"`
	Name         string      `json:"name" dc:"分类名称"`
	SortOrder    int         `json:"sortOrder" dc:"排序值"`
	ProductCount int         `json:"productCount" dc:"商品数量"`
	Status       int         `json:"status" dc:"状态：1启用，0禁用"`
	Image        string      `json:"image" dc:"分类图片URL"`
	CreatedAt    *gtime.Time `json:"createdAt" dc:"创建时间"`
}

// WxShopCategoryListReq 微信客户端获取商城分类列表请求
type WxShopCategoryListReq struct {
	g.Meta    `path:"/client/shop-category/list" method:"get" tags:"客户端商城分类" summary:"客户端获取商城分类列表" security:"Bearer" description:"客户端获取商城分类列表，支持全量获取和分页查询"`
	Page      int    `v:"min:0#页码最小值为0"  dc:"页码，为0时表示不分页获取全部" json:"page" d:"0"`
	Size      int    `v:"max:50#每页最大50条" dc:"每页数量，当page>0时生效" json:"size" d:"10"`
	SortField string `dc:"排序字段: id, name, sort_order, product_count, created_at" json:"sortField" d:"sort_order"`
	SortOrder string `dc:"排序方式: asc, desc" json:"sortOrder" d:"asc"`
}

// WxShopCategoryListRes 微信客户端获取商城分类列表响应
type WxShopCategoryListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ShopCategoryItem `json:"list" dc:"商城分类列表"`
	Total  int                `json:"total,omitempty" dc:"总记录数(分页时返回)"`
	Page   int                `json:"page,omitempty" dc:"当前页码(分页时返回)"`
	Size   int                `json:"size,omitempty" dc:"每页记录数(分页时返回)"`
	Pages  int                `json:"pages,omitempty" dc:"总页数(分页时返回)"`
}

// ShopCategorySyncProductCountReq 同步商品分类数量请求
type ShopCategorySyncProductCountReq struct {
	g.Meta `path:"/shop-category/sync-product-count" method:"post" tags:"商城分类" summary:"同步商品分类数量" security:"Bearer" description:"同步所有商品分类的商品数量"`
}

// ShopCategorySyncProductCountRes 同步商品分类数量响应
type ShopCategorySyncProductCountRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Status bool `json:"status" dc:"同步状态：true成功，false失败"`
}
