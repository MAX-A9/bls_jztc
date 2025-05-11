package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// WxMyPublishListReq 客户端获取我的发布列表请求
type WxMyPublishListReq struct {
	g.Meta   `path:"/client/publish/list" method:"get" tags:"客户端内容" summary:"获取我的发布列表" security:"Bearer" description:"获取当前登录客户发布的内容列表，支持分页和筛选"`
	Page     int    `v:"min:1#页码最小值为1" d:"1" json:"page" query:"page" dc:"页码，默认1"`
	PageSize int    `v:"max:50#每页最大50条" d:"10" json:"pageSize" query:"pageSize" dc:"每页数量，默认10"`
	Type     int    `d:"0" json:"type" query:"type" dc:"内容类型，0:全部，1:普通信息，2:闲置"`
	Status   string `json:"status" query:"status" dc:"状态筛选：已发布、待审核、已下架"`
}

// WxMyPublishListRes 客户端获取我的发布列表响应
type WxMyPublishListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []MyPublishItem `json:"list" dc:"发布列表"`
	Total  int             `json:"total" dc:"总记录数"`
	Page   int             `json:"page" dc:"当前页码"`
	Pages  int             `json:"pages" dc:"总页数"`
}

// MyPublishItem 我的发布列表项
type MyPublishItem struct {
	Id          int         `json:"id" dc:"内容ID"`
	Title       string      `json:"title" dc:"标题"`
	Category    string      `json:"category" dc:"分类"`
	Status      string      `json:"status" dc:"状态：已发布、待审核、已下架"`
	PublishedAt *gtime.Time `json:"publishedAt" dc:"发布时间"`
}

// WxMyPublishCountReq 客户端获取我的发布数量请求
type WxMyPublishCountReq struct {
	g.Meta `path:"/client/publish/count" method:"get" tags:"客户端内容" summary:"获取我的发布数量" security:"Bearer" description:"获取当前登录客户发布的内容数量"`
}

// WxMyPublishCountRes 客户端获取我的发布数量响应
type WxMyPublishCountRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Total  int `json:"total" dc:"总数量"`
}
