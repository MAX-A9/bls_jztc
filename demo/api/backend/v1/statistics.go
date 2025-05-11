package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// StatisticsDataReq 统计数据请求
type StatisticsDataReq struct {
	g.Meta    `path:"/statistics/data" method:"get" tags:"统计" summary:"获取统计数据" security:"Bearer" description:"获取统计数据，包含注册客户数量、兑换数量、发布数量、收益金额"`
	TimeRange string `v:"required|in:week,month,year#时间范围不能为空|时间范围参数错误" json:"timeRange" dc:"时间范围：week-本周 month-本月 year-本年"`
}

// StatisticsDataRes 统计数据响应
type StatisticsDataRes struct {
	g.Meta        `mime:"application/json" example:"json"`
	ClientCount   int     `json:"clientCount" dc:"注册客户数量"`
	ExchangeCount int     `json:"exchangeCount" dc:"兑换数量"`
	PublishCount  int     `json:"publishCount" dc:"发布数量"`
	RevenueAmount float64 `json:"revenueAmount" dc:"收益金额"`
	TimeRange     string  `json:"timeRange" dc:"时间范围"`
	StartTime     string  `json:"startTime" dc:"开始时间"`
	EndTime       string  `json:"endTime" dc:"结束时间"`
}
