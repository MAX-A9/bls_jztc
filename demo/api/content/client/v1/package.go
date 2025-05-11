package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PackageType 套餐类型
type PackageType string

const (
	PackageTypeTop     PackageType = "top"     // 置顶套餐
	PackageTypePublish PackageType = "publish" // 发布套餐
)

// DurationType 时长单位类型
type DurationType string

const (
	DurationTypeHour  DurationType = "hour"  // 按小时
	DurationTypeDay   DurationType = "day"   // 按天
	DurationTypeMonth DurationType = "month" // 按月
)

// Package 套餐基础结构
type Package struct {
	Id           int          `json:"id" dc:"套餐ID"`
	Title        string       `json:"title" dc:"套餐名称"`
	Description  string       `json:"description" dc:"套餐简介"`
	Price        float64      `json:"price" dc:"价格(元)"`
	Type         PackageType  `json:"type" dc:"套餐类型: top-置顶套餐, publish-发布套餐"`
	Duration     int          `json:"duration" dc:"时长值"`
	DurationType DurationType `json:"durationType" dc:"时长单位: hour-小时, day-天, month-月"`
	SortOrder    int          `json:"sortOrder" dc:"排序值，数字越小排序越靠前"`
}

// WxClientPackageListReq 客户端获取套餐列表请求
type WxClientPackageListReq struct {
	g.Meta `path:"/wx/client/package/list" method:"get" tags:"客户端套餐" summary:"获取套餐列表" description:"获取所有可用的套餐列表"`
	Sort   string `json:"sort" dc:"排序字段，支持price（价格）、duration（时长），默认按时长单位排序"`
	Order  string `v:"in:asc,desc#排序方式只能是asc或desc" json:"order" dc:"排序方式: asc-升序, desc-降序，默认升序"`
}

// WxClientPackageListRes 客户端套餐列表响应
type WxClientPackageListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	TopPackages     []*Package `json:"topPackages" dc:"置顶套餐列表"`
	PublishPackages []*Package `json:"publishPackages" dc:"发布套餐列表"`
	TopEnabled      bool       `json:"topEnabled" dc:"置顶套餐总开关状态"`
	PublishEnabled  bool       `json:"publishEnabled" dc:"发布套餐总开关状态"`
}
