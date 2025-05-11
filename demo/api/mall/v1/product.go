package v1

import (
	"demo/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// 商品列表请求
type ProductListReq struct {
	g.Meta     `path:"/product/list" method:"get" tags:"商品管理" summary:"获取商品列表" security:"Bearer" description:"获取商品列表，需要管理员权限"`
	Page       int    `v:"min:1#页码最小值为1"  dc:"页码" json:"page" d:"1"`
	Size       int    `v:"max:50#每页最大50条" dc:"每页数量" json:"size" d:"10"`
	Name       string `dc:"商品名称, 支持模糊搜索" json:"name"`
	CategoryId int    `dc:"分类ID" json:"categoryId"`
	Status     int    `dc:"状态 0:未上架 1:已上架 2:已售罄" json:"status"`
	Duration   int    `dc:"所需时长(天)" json:"duration"`
	Stock      int    `dc:"库存" json:"stock"`
	Tags       string `dc:"商品标签，多个标签用逗号分隔" json:"tags"`
	SortField  string `dc:"排序字段: id, name, price, stock, duration, sales, sort_order, created_at" json:"sortField" d:"sort_order"`
	SortOrder  string `dc:"排序方式: asc, desc" json:"sortOrder" d:"asc"`
}

// 商品列表响应
type ProductListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.Product `json:"list"`  // 商品列表
	Total  int              `json:"total"` // 总数
	Page   int              `json:"page"`  // 页码
	Size   int              `json:"size"`  // 每页数量
}

// 商品详情请求
type ProductDetailReq struct {
	g.Meta `path:"/product/detail" method:"get" tags:"商品管理" summary:"获取商品详情" security:"Bearer" description:"获取商品详情，需要管理员权限"`
	Id     int `v:"required#商品ID不能为空" dc:"商品ID" json:"id"`
}

// 商品详情响应
type ProductDetailRes struct {
	g.Meta  `mime:"application/json"`
	Product *model.Product `json:"product"` // 商品信息
}

// 创建商品请求
type ProductCreateReq struct {
	g.Meta       `path:"/product/create" method:"post" tags:"商品管理" summary:"创建商品" security:"Bearer" description:"创建商品，需要管理员权限"`
	Name         string  `v:"required#商品名称不能为空" dc:"商品名称" json:"name"`
	CategoryId   int     `v:"required#分类ID不能为空" dc:"分类ID" json:"categoryId"`
	CategoryName string  `v:"required#分类名称不能为空" dc:"分类名称" json:"categoryName"`
	Price        float64 `v:"required|min:0.01#价格不能为空|价格必须大于0" dc:"商品价格" json:"price"`
	Duration     int     `v:"required|min:1#时长不能为空|时长必须大于0" dc:"所需时长(分钟)" json:"duration"`
	Stock        int     `v:"required|min:0#库存不能为空|库存不能小于0" dc:"库存数量" json:"stock"`
	Status       int     `v:"required|in:0,1,2#状态不能为空|状态值只能是0,1,2" dc:"状态 0:未上架 1:已上架 2:已售罄" json:"status"`
	SortOrder    int     `dc:"排序值，数字越小排序越靠前" json:"sortOrder" d:"0"`
	Description  string  `dc:"商品描述" json:"description"`
	Thumbnail    string  `dc:"缩略图URL" json:"thumbnail"`
	Images       string  `dc:"商品图片URLs，JSON格式" json:"images"`
	Tags         string  `dc:"商品标签，多个标签用逗号分隔" json:"tags"`
}

// 创建商品响应
type ProductCreateRes struct {
	g.Meta `mime:"application/json"`
	Id     int `json:"id"` // 新创建的商品ID
}

// 更新商品请求
type ProductUpdateReq struct {
	g.Meta       `path:"/product/update" method:"post" tags:"商品管理" summary:"更新商品" security:"Bearer" description:"更新商品，需要管理员权限"`
	Id           int     `v:"required#商品ID不能为空" dc:"商品ID" json:"id"`
	Name         string  `dc:"商品名称" json:"name"`
	CategoryId   int     `dc:"分类ID" json:"categoryId"`
	CategoryName string  `dc:"分类名称" json:"categoryName"`
	Price        float64 `v:"min:0.01#价格必须大于0" dc:"商品价格" json:"price"`
	Duration     int     `v:"min:1#时长必须大于0" dc:"所需时长(分钟)" json:"duration"`
	Stock        int     `v:"min:0#库存不能小于0" dc:"库存数量" json:"stock"`
	Status       int     `v:"in:0,1,2#状态值只能是0,1,2" dc:"状态 0:未上架 1:已上架 2:已售罄" json:"status"`
	SortOrder    int     `dc:"排序值，数字越小排序越靠前" json:"sortOrder"`
	Description  string  `dc:"商品描述" json:"description"`
	Thumbnail    string  `dc:"缩略图URL" json:"thumbnail"`
	Images       string  `dc:"商品图片URLs，JSON格式" json:"images"`
	Tags         string  `dc:"商品标签，多个标签用逗号分隔" json:"tags"`
}

// 更新商品响应
type ProductUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// 删除商品请求
type ProductDeleteReq struct {
	g.Meta `path:"/product/delete" method:"post" tags:"商品管理" summary:"删除商品" security:"Bearer" description:"删除商品，需要管理员权限"`
	Id     int `v:"required#商品ID不能为空" dc:"商品ID" json:"id"`
}

// 删除商品响应
type ProductDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// 商品上/下架请求
type ProductStatusUpdateReq struct {
	g.Meta `path:"/product/status" method:"post" tags:"商品管理" summary:"更新商品状态(上架/下架)" security:"Bearer" description:"更新商品状态，需要管理员权限"`
	Id     int `v:"required#商品ID不能为空" dc:"商品ID" json:"id"`
	Status int `v:"required|in:0,1,2#状态不能为空|状态值只能是0,1,2" dc:"状态 0:未上架 1:已上架 2:已售罄" json:"status"`
}

// 商品上/下架响应
type ProductStatusUpdateRes struct {
	g.Meta `mime:"application/json"`
}
