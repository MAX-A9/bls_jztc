package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 套餐表相关常量
const (
	// TablePackage 表名
	TablePackage = "package"
	// PackageColumns 所有字段
	PackageColumns = "id,title,description,price,type,duration,duration_type,sort_order,created_at,updated_at,deleted_at"
)

// 套餐类型常量
const (
	// PackageTypeTop 置顶套餐
	PackageTypeTop = "top"
	// PackageTypePublish 发布套餐
	PackageTypePublish = "publish"
)

// 时长单位类型常量
const (
	// DurationTypeHour 按小时
	DurationTypeHour = "hour"
	// DurationTypeDay 按天
	DurationTypeDay = "day"
	// DurationTypeMonth 按月
	DurationTypeMonth = "month"
)

// PackageDO 套餐DO结构体
type PackageDO struct {
	g.Meta       `orm:"table:package, do:true"`
	Id           interface{} `orm:"id,primary"     json:"id"`           // 套餐ID
	Title        interface{} `orm:"title"          json:"title"`        // 套餐名称
	Description  interface{} `orm:"description"    json:"description"`  // 套餐简介
	Price        interface{} `orm:"price"          json:"price"`        // 价格(元)
	Type         interface{} `orm:"type"           json:"type"`         // 套餐类型: top-置顶套餐, publish-发布套餐
	Duration     interface{} `orm:"duration"       json:"duration"`     // 时长值
	DurationType interface{} `orm:"duration_type"  json:"durationType"` // 时长单位: hour-小时, day-天, month-月
	SortOrder    interface{} `orm:"sort_order"     json:"sortOrder"`    // 排序值，数字越小排序越靠前
	CreatedAt    *gtime.Time `orm:"created_at"     json:"createdAt"`    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"     json:"updatedAt"`    // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"     json:"deletedAt"`    // 删除时间
}
