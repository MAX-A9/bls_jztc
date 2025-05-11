package v1

import (
	"demo/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// 微信端商品列表请求
type WxProductListReq struct {
	g.Meta     `path:"/product/list" method:"get" tags:"微信小程序-商品管理" summary:"获取商品列表"`
	Page       int    `v:"min:1#页码最小值为1"  dc:"页码" json:"page" d:"1"`
	Size       int    `v:"max:50#每页最大50条" dc:"每页数量" json:"size" d:"10"`
	Name       string `dc:"商品名称, 支持模糊搜索" json:"name"`
	CategoryId int    `dc:"分类ID" json:"categoryId"`
	Duration   int    `dc:"所需时长(天)" json:"duration"`
	Tags       string `dc:"商品标签，多个标签用逗号分隔" json:"tags"`
	SortField  string `dc:"排序字段: id, name, price, duration, sales, sort_order, created_at" json:"sortField" d:"sort_order"`
	SortOrder  string `dc:"排序方式: asc, desc" json:"sortOrder" d:"asc"`
}

// 微信端商品列表响应
type WxProductListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.Product `json:"list"`  // 商品列表
	Total  int              `json:"total"` // 总数
	Page   int              `json:"page"`  // 页码
	Size   int              `json:"size"`  // 每页数量
}

// 微信端商品详情请求
type WxProductDetailReq struct {
	g.Meta `path:"/product/detail" method:"get" tags:"微信小程序-商品管理" summary:"获取商品详情"`
	Id     int `v:"required#商品ID不能为空" dc:"商品ID" json:"id"`
}

// 微信端商品详情响应
type WxProductDetailRes struct {
	g.Meta  `mime:"application/json"`
	Product *model.Product `json:"product"` // 商品信息
}
