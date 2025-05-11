package v1

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// 浏览历史记录筛选时间类型
const (
	BrowseHistoryTimeTypeAll       = "all"        // 全部
	BrowseHistoryTimeTypeToday     = "today"      // 今天
	BrowseHistoryTimeTypeYesterday = "yesterday"  // 昨天
	BrowseHistoryTimeTypeThisWeek  = "this_week"  // 本周
	BrowseHistoryTimeTypeThisMonth = "this_month" // 本月
)

// 浏览历史记录列表请求
type BrowseHistoryListReq struct {
	g.Meta   `path:"/client/browse-history/list" method:"get" tags:"浏览历史" summary:"获取浏览历史记录" security:"Bearer" description:"客户获取自己的浏览历史记录，支持按时间筛选"`
	TimeType string `v:"in:all,today,yesterday,this_week,this_month#时间类型参数错误" json:"timeType" dc:"时间筛选类型，可选值：all-全部，today-今天，yesterday-昨天，this_week-本周，this_month-本月"`
	Page     int    `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int    `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
}

// 浏览历史记录列表响应
type BrowseHistoryListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []BrowseHistoryItem `json:"list" dc:"浏览历史记录列表"`
	Total  int                 `json:"total" dc:"总数量"`
	Page   int                 `json:"page" dc:"当前页码"`
}

// 浏览历史记录项
type BrowseHistoryItem struct {
	Id            int     `json:"id" dc:"历史记录ID"`
	ContentId     int     `json:"contentId" dc:"内容ID"`
	ContentType   string  `json:"contentType" dc:"内容类型，如article、idle等"`
	ContentTitle  string  `json:"contentTitle" dc:"内容标题"`
	ContentCover  string  `json:"contentCover" dc:"内容封面图"`
	ContentStatus int     `json:"contentStatus" dc:"内容状态"`
	BrowseTime    string  `json:"browseTime" dc:"浏览时间"`
	Category      string  `json:"category" dc:"内容分类"`
	Price         float64 `json:"price" dc:"价格（仅对闲置物品有效）"`
}

// 清空浏览历史记录请求
type BrowseHistoryClearReq struct {
	g.Meta   `path:"/client/browse-history/clear" method:"post" tags:"浏览历史" summary:"清空浏览历史记录" security:"Bearer" description:"客户清空自己的浏览历史记录"`
	TimeType string `v:"in:all,today,yesterday,this_week,this_month#时间类型参数错误" json:"timeType" dc:"时间筛选类型，可选值：all-全部，today-今天，yesterday-昨天，this_week-本周，this_month-本月"`
}

// 清空浏览历史记录响应
type BrowseHistoryClearRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 添加浏览历史记录请求
type BrowseHistoryAddReq struct {
	g.Meta      `path:"/client/browse-history/add" method:"post" tags:"浏览历史" summary:"添加浏览历史记录" security:"Bearer" description:"添加浏览历史记录，通常由前端自动调用"`
	ContentId   int    `v:"required#内容ID不能为空" json:"contentId" dc:"内容ID"`
	ContentType string `v:"required|in:article,idle,info#内容类型不能为空|内容类型必须是article、idle或info" json:"contentType" dc:"内容类型，可选值：article、idle、info等"`
}

// 添加浏览历史记录响应
type BrowseHistoryAddRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// 浏览历史记录数量请求
type BrowseHistoryCountReq struct {
	g.Meta `path:"/client/browse-history/count" method:"get" tags:"浏览历史" summary:"获取浏览历史记录数量" security:"Bearer" description:"客户获取自己的所有浏览历史记录总数"`
}

// 浏览历史记录数量响应
type BrowseHistoryCountRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Count  int `json:"count" dc:"浏览历史记录数量"`
}

// IBrowseHistory 浏览历史API
type IBrowseHistory interface {
	// ... existing code ...

	// Count 获取浏览历史记录数量
	Count(ctx context.Context, req *BrowseHistoryCountReq) (res *BrowseHistoryCountRes, err error)
}
