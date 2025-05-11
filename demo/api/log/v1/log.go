package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationLogListReq 操作日志列表请求
type OperationLogListReq struct {
	g.Meta    `path:"/list" method:"get" tags:"操作日志" summary:"获取操作日志列表" security:"Bearer" description:"获取操作日志列表，需要管理员权限"`
	Page      int    `json:"page" v:"min:1#页码最小值为1" dc:"页码"`
	PageSize  int    `json:"pageSize" v:"max:100#每页最大100条" dc:"每页数量"`
	Username  string `json:"username" dc:"用户名"`
	Module    string `json:"module" dc:"模块"`
	Action    string `json:"action" dc:"操作类型"`
	Result    string `json:"result" dc:"操作结果 0:失败 1:成功"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword   string `json:"keyword" dc:"关键字"`
}

// OperationLogListRes 操作日志列表响应
type OperationLogListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []OperationLogListItem `json:"list" dc:"日志列表"`
	Total  int                    `json:"total" dc:"总数"`
	Page   int                    `json:"page" dc:"当前页码"`
}

// OperationLogListItem 操作日志列表项
type OperationLogListItem struct {
	Id              int         `json:"id" dc:"日志ID"`
	UserId          int         `json:"userId" dc:"用户ID"`
	Username        string      `json:"username" dc:"用户名"`
	OperationIp     string      `json:"operationIp" dc:"操作IP"`
	OperationTime   *gtime.Time `json:"operationTime" dc:"操作时间"`
	Module          string      `json:"module" dc:"模块"`
	Action          string      `json:"action" dc:"操作类型"`
	OperationResult int         `json:"operationResult" dc:"操作结果 0:失败 1:成功"`
	ResultText      string      `json:"resultText" dc:"操作结果文本"`
}

// OperationLogDeleteReq 删除操作日志请求
type OperationLogDeleteReq struct {
	g.Meta `path:"/delete" method:"delete" tags:"操作日志" summary:"删除操作日志" security:"Bearer" description:"删除操作日志，需要管理员权限"`
	Id     int `json:"id" v:"required#日志ID不能为空" dc:"日志ID"`
}

// OperationLogDeleteRes 删除操作日志响应
type OperationLogDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// OperationLogExportReq 导出操作日志请求
type OperationLogExportReq struct {
	g.Meta    `path:"/export" method:"get" tags:"操作日志" summary:"导出操作日志" security:"Bearer" description:"导出操作日志，需要管理员权限"`
	Username  string `json:"username" dc:"用户名"`
	Module    string `json:"module" dc:"模块"`
	Action    string `json:"action" dc:"操作类型"`
	Result    string `json:"result" dc:"操作结果 0:失败 1:成功"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword   string `json:"keyword" dc:"关键字"`
}

// OperationLogExportRes 导出操作日志响应
type OperationLogExportRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Url    string `json:"url" dc:"下载地址"`
}
