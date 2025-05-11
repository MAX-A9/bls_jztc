package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 统计数据请求
type StatsDataReq struct {
	g.Meta     `path:"/stats/data" method:"get" tags:"统计数据" summary:"获取统计数据" security:"Bearer" description:"获取关键统计数据，包括注册客户数量、兑换数量、发布数量、收益金额"`
	PeriodType string `d:"week" json:"periodType" dc:"周期类型：week-本周，month-本月，year-本年"`
}

// 统计数据响应
type StatsDataRes struct {
	g.Meta        `mime:"application/json" example:"json"`
	ClientCount   int     `json:"clientCount" dc:"注册客户数量"`
	ExchangeCount int     `json:"exchangeCount" dc:"兑换数量"`
	PublishCount  int     `json:"publishCount" dc:"发布数量"`
	RevenueAmount float64 `json:"revenueAmount" dc:"收益金额"`
}

// 趋势分析请求
type StatsTrendReq struct {
	g.Meta     `path:"/stats/trend" method:"get" tags:"统计数据" summary:"获取趋势分析数据" security:"Bearer" description:"获取趋势分析数据，支持按周期类型查询"`
	PeriodType string `d:"week" json:"periodType" dc:"周期类型：week-本周，month-本月，year-本年"`
	DataType   string `v:"required#数据类型不能为空" json:"dataType" dc:"数据类型：clients-客户数量，exchanges-兑换数量，publishes-发布数量，revenue-收益金额，all-所有数据"`
}

// 趋势分析响应
type StatsTrendRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	Labels    []string                 `json:"labels" dc:"时间标签"`
	Values    []interface{}            `json:"values" dc:"数据值"`
	AllValues map[string][]interface{} `json:"allValues" dc:"所有类型的数据值，仅在dataType为all时使用"`
	DataType  string                   `json:"dataType" dc:"数据类型"`
	Period    string                   `json:"period" dc:"周期类型"`
}
