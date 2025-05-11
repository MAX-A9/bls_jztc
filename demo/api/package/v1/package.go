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
	CreatedAt    string       `json:"createdAt" dc:"创建时间"`
	UpdatedAt    string       `json:"updatedAt" dc:"更新时间"`
}

// PackageListReq 获取套餐列表请求
type PackageListReq struct {
	g.Meta `path:"/package/list" method:"get" tags:"套餐管理" summary:"获取套餐列表" security:"Bearer" description:"获取套餐列表，支持按类型筛选"`
	Type   PackageType `v:"in:top,publish#套餐类型只能是top或publish" json:"type" dc:"套餐类型: top-置顶套餐, publish-发布套餐"`
	Sort   string      `json:"sort" dc:"排序字段，支持price（价格）、duration（时长），默认按类型和时长单位排序"`
	Order  string      `v:"in:asc,desc#排序方式只能是asc或desc" json:"order" dc:"排序方式: asc-升序, desc-降序，默认升序"`
}

// PackageListRes 套餐列表响应
type PackageListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []*Package `json:"list" dc:"套餐列表"`
	IsGlobalEnabled bool       `json:"isGlobalEnabled" dc:"当前类型的套餐总开关状态"`
}

// PackageDetailReq 获取套餐详情请求
type PackageDetailReq struct {
	g.Meta `path:"/package/detail" method:"get" tags:"套餐管理" summary:"获取套餐详情" security:"Bearer" description:"获取套餐详情"`
	Id     int `v:"required|min:1#套餐ID不能为空|套餐ID必须大于0" json:"id" dc:"套餐ID"`
}

// PackageDetailRes 套餐详情响应
type PackageDetailRes struct {
	g.Meta `mime:"application/json" example:"json"`
	*Package
}

// PackageCreateReq 创建套餐请求
type PackageCreateReq struct {
	g.Meta       `path:"/package/create" method:"post" tags:"套餐管理" summary:"创建套餐" security:"Bearer" description:"创建新的套餐，需要管理员权限"`
	Title        string       `v:"required#套餐名称不能为空" json:"title" dc:"套餐名称"`
	Description  string       `v:"required#套餐简介不能为空" json:"description" dc:"套餐简介"`
	Price        float64      `v:"required|min:0#价格不能为空|价格不能小于0" json:"price" dc:"价格(元)"`
	Type         PackageType  `v:"required|in:top,publish#套餐类型不能为空|套餐类型只能是top或publish" json:"type" dc:"套餐类型: top-置顶套餐, publish-发布套餐"`
	Duration     int          `v:"required|min:1#时长不能为空|时长必须大于0" json:"duration" dc:"时长值"`
	DurationType DurationType `v:"required|in:hour,day,month#时长单位不能为空|时长单位只能是hour,day或month" json:"durationType" dc:"时长单位: hour-小时, day-天, month-月"`
	SortOrder    int          `json:"sortOrder" dc:"排序值，数字越小排序越靠前"`
}

// PackageCreateRes 创建套餐响应
type PackageCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"套餐ID"`
}

// PackageUpdateReq 更新套餐请求
type PackageUpdateReq struct {
	g.Meta       `path:"/package/update" method:"put" tags:"套餐管理" summary:"更新套餐" security:"Bearer" description:"更新套餐信息，需要管理员权限"`
	Id           int          `v:"required|min:1#套餐ID不能为空|套餐ID必须大于0" json:"id" dc:"套餐ID"`
	Title        string       `v:"required#套餐名称不能为空" json:"title" dc:"套餐名称"`
	Description  string       `v:"required#套餐简介不能为空" json:"description" dc:"套餐简介"`
	Price        float64      `v:"required|min:0#价格不能为空|价格不能小于0" json:"price" dc:"价格(元)"`
	Type         PackageType  `v:"required|in:top,publish#套餐类型不能为空|套餐类型只能是top或publish" json:"type" dc:"套餐类型: top-置顶套餐, publish-发布套餐"`
	Duration     int          `v:"required|min:1#时长不能为空|时长必须大于0" json:"duration" dc:"时长值"`
	DurationType DurationType `v:"required|in:hour,day,month#时长单位不能为空|时长单位只能是hour,day或month" json:"durationType" dc:"时长单位: hour-小时, day-天, month-月"`
	SortOrder    int          `json:"sortOrder" dc:"排序值，数字越小排序越靠前"`
}

// PackageUpdateRes 更新套餐响应
type PackageUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// PackageDeleteReq 删除套餐请求
type PackageDeleteReq struct {
	g.Meta `path:"/package/delete" method:"delete" tags:"套餐管理" summary:"删除套餐" security:"Bearer" description:"删除套餐，需要管理员权限"`
	Id     int `v:"required|min:1#套餐ID不能为空|套餐ID必须大于0" json:"id" dc:"套餐ID"`
}

// PackageDeleteRes 删除套餐响应
type PackageDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// PackageGlobalStatusReq 获取套餐总开关状态请求
type PackageGlobalStatusReq struct {
	g.Meta `path:"/package/global-status" method:"get" tags:"套餐管理" summary:"获取套餐总开关状态" security:"Bearer" description:"获取套餐总开关状态，包括置顶套餐和发布套餐的总开关状态"`
}

// PackageGlobalStatusRes 获取套餐总开关状态响应
type PackageGlobalStatusRes struct {
	g.Meta         `mime:"application/json" example:"json"`
	TopEnabled     bool `json:"topEnabled" dc:"置顶套餐总开关状态"`
	PublishEnabled bool `json:"publishEnabled" dc:"发布套餐总开关状态"`
}

// TopPackageGlobalStatusUpdateReq 更新置顶套餐总开关请求
type TopPackageGlobalStatusUpdateReq struct {
	g.Meta    `path:"/package/top/global-status/update" method:"put" tags:"套餐管理" summary:"更新置顶套餐总开关" security:"Bearer" description:"更新置顶套餐总开关状态，需要管理员权限"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用总开关"`
}

// TopPackageGlobalStatusUpdateRes 更新置顶套餐总开关响应
type TopPackageGlobalStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// PublishPackageGlobalStatusUpdateReq 更新发布套餐总开关请求
type PublishPackageGlobalStatusUpdateReq struct {
	g.Meta    `path:"/package/publish/global-status/update" method:"put" tags:"套餐管理" summary:"更新发布套餐总开关" security:"Bearer" description:"更新发布套餐总开关状态，需要管理员权限"`
	IsEnabled bool `json:"isEnabled" dc:"是否启用总开关"`
}

// PublishPackageGlobalStatusUpdateRes 更新发布套餐总开关响应
type PublishPackageGlobalStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// WxPackageListReq 客户端获取套餐列表请求
type WxPackageListReq struct {
	g.Meta `path:"/wx/package/list" method:"get" tags:"客户端套餐" summary:"获取套餐列表" description:"获取套餐列表，支持按类型筛选"`
	Type   PackageType `v:"in:top,publish#套餐类型只能是top或publish" json:"type" dc:"套餐类型: top-置顶套餐, publish-发布套餐"`
	Sort   string      `json:"sort" dc:"排序字段，支持price（价格）、duration（时长），默认按类型和时长单位排序"`
	Order  string      `v:"in:asc,desc#排序方式只能是asc或desc" json:"order" dc:"排序方式: asc-升序, desc-降序，默认升序"`
}

// WxPackageListRes 客户端套餐列表响应
type WxPackageListRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	List            []*Package `json:"list" dc:"套餐列表"`
	IsGlobalEnabled bool       `json:"isGlobalEnabled" dc:"当前类型的套餐总开关状态"`
}
