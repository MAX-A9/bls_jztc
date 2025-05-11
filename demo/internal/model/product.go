package model

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// 商品基础信息结构
type Product struct {
	Id           int         `json:"id"`           // 商品ID
	Name         string      `json:"name"`         // 商品名称
	CategoryId   int         `json:"categoryId"`   // 分类ID
	CategoryName string      `json:"categoryName"` // 分类名称
	Price        float64     `json:"price"`        // 商品价格
	Duration     int         `json:"duration"`     // 所需时长(天)
	Stock        int         `json:"stock"`        // 库存数量
	Status       int         `json:"status"`       // 状态 0:未上架 1:已上架 2:已售罄
	SortOrder    int         `json:"sortOrder"`    // 排序值
	Description  string      `json:"description"`  // 商品描述
	Thumbnail    string      `json:"thumbnail"`    // 缩略图URL
	Images       string      `json:"images"`       // 商品图片URLs，JSON格式
	Tags         string      `json:"tags"`         // 商品标签，多个标签用逗号分隔
	Sales        int         `json:"sales"`        // 销量
	CreatedAt    *gtime.Time `json:"createdAt"`    // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"`    // 更新时间
}

// 商品过滤条件
type ProductFilter struct {
	Name        string `json:"name"`        // 商品名称，模糊搜索
	CategoryId  int    `json:"categoryId"`  // 分类ID
	Status      int    `json:"status"`      // 状态
	MinDuration int    `json:"minDuration"` // 最小时长
	MaxDuration int    `json:"maxDuration"` // 最大时长
	MinStock    int    `json:"minStock"`    // 最小库存
	SortField   string `json:"sortField"`   // 排序字段
	SortOrder   string `json:"sortOrder"`   // 排序方式: asc, desc
	Tags        string `json:"tags"`        // 标签过滤
}

// 商品状态常量
const (
	ProductStatusUnlisted = 0 // 未上架
	ProductStatusListed   = 1 // 已上架
	ProductStatusSoldOut  = 2 // 已售罄
)

// 商品状态文本映射
var ProductStatusText = map[int]string{
	ProductStatusUnlisted: "未上架",
	ProductStatusListed:   "已上架",
	ProductStatusSoldOut:  "已售罄",
}
