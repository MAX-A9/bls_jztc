package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExchangeRecordGetReq 获取兑换记录详情请求
type ExchangeRecordGetReq struct {
	g.Meta `path:"/exchange-record/get" method:"get" tags:"兑换记录" summary:"获取兑换记录详情" security:"Bearer" description:"获取指定ID的兑换记录详情"`
	Id     int `v:"required#兑换记录ID不能为空" dc:"兑换记录ID" json:"id" query:"id"`
}

// ExchangeRecordGetRes 获取兑换记录详情响应
type ExchangeRecordGetRes struct {
	g.Meta          `mime:"application/json" example:"json"`
	Id              int         `json:"id" dc:"兑换记录ID"`
	ClientId        int         `json:"clientId" dc:"客户ID"`
	ClientName      string      `json:"clientName" dc:"客户名称"`
	RechargeAccount string      `json:"rechargeAccount" dc:"充值账号"`
	ProductName     string      `json:"productName" dc:"商品名称"`
	Duration        int         `json:"duration" dc:"消耗时长(分钟)"`
	ExchangeTime    *gtime.Time `json:"exchangeTime" dc:"兑换时间"`
	Status          string      `json:"status" dc:"状态"`
	Remark          string      `json:"remark" dc:"备注"`
	CreatedAt       *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// ExchangeRecordCreateReq 创建兑换记录请求
type ExchangeRecordCreateReq struct {
	g.Meta          `path:"/exchange-record/create" method:"post" tags:"兑换记录" summary:"创建兑换记录" security:"Bearer" description:"创建新的兑换记录"`
	ClientId        int         `v:"required#客户ID不能为空" json:"clientId" dc:"客户ID"`
	ClientName      string      `v:"required#客户名称不能为空" json:"clientName" dc:"客户名称"`
	RechargeAccount string      `v:"required#充值账号不能为空" json:"rechargeAccount" dc:"充值账号"`
	ProductName     string      `v:"required#商品名称不能为空" json:"productName" dc:"商品名称"`
	Duration        int         `v:"required#消耗时长不能为空" json:"duration" dc:"消耗时长(分钟)"`
	ExchangeTime    *gtime.Time `json:"exchangeTime" dc:"兑换时间，默认为当前时间"`
	Status          string      `json:"status" dc:"状态：processing、completed、failed，默认为processing"`
	Remark          string      `json:"remark" dc:"备注"`
}

// ExchangeRecordCreateRes 创建兑换记录响应
type ExchangeRecordCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"新创建的兑换记录ID"`
}

// ExchangeRecordUpdateReq 更新兑换记录请求
type ExchangeRecordUpdateReq struct {
	g.Meta          `path:"/exchange-record/update" method:"put" tags:"兑换记录" summary:"更新兑换记录" security:"Bearer" description:"更新指定ID的兑换记录信息"`
	Id              int         `v:"required#兑换记录ID不能为空" json:"id" dc:"兑换记录ID"`
	ClientId        int         `v:"required#客户ID不能为空" json:"clientId" dc:"客户ID"`
	ClientName      string      `v:"required#客户名称不能为空" json:"clientName" dc:"客户名称"`
	RechargeAccount string      `v:"required#充值账号不能为空" json:"rechargeAccount" dc:"充值账号"`
	ProductName     string      `v:"required#商品名称不能为空" json:"productName" dc:"商品名称"`
	Duration        int         `v:"required#消耗时长不能为空" json:"duration" dc:"消耗时长(分钟)"`
	ExchangeTime    *gtime.Time `json:"exchangeTime" dc:"兑换时间"`
	Status          string      `json:"status" dc:"状态：processing、completed、failed"`
	Remark          string      `json:"remark" dc:"备注"`
}

// ExchangeRecordUpdateRes 更新兑换记录响应
type ExchangeRecordUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ExchangeRecordDeleteReq 删除兑换记录请求
type ExchangeRecordDeleteReq struct {
	g.Meta `path:"/exchange-record/delete" method:"delete" tags:"兑换记录" summary:"删除兑换记录" security:"Bearer" description:"删除指定ID的兑换记录"`
	Id     int `v:"required#兑换记录ID不能为空" json:"id" dc:"兑换记录ID"`
}

// ExchangeRecordDeleteRes 删除兑换记录响应
type ExchangeRecordDeleteRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ExchangeRecordStatusUpdateReq 更新兑换记录状态请求
type ExchangeRecordStatusUpdateReq struct {
	g.Meta `path:"/exchange-record/status/update" method:"put" tags:"兑换记录" summary:"更新兑换记录状态" security:"Bearer" description:"更新指定ID的兑换记录状态"`
	Id     int    `v:"required#兑换记录ID不能为空" json:"id" dc:"兑换记录ID"`
	Status string `v:"required#状态不能为空" json:"status" dc:"状态：processing、completed、failed"`
}

// ExchangeRecordStatusUpdateRes 更新兑换记录状态响应
type ExchangeRecordStatusUpdateRes struct {
	g.Meta `mime:"application/json" example:"json"`
}

// ExchangeRecordListReq 获取兑换记录列表请求
type ExchangeRecordListReq struct {
	g.Meta   `path:"/exchange-record/list" method:"get" tags:"兑换记录" summary:"获取兑换记录列表" security:"Bearer" description:"获取兑换记录列表，支持分页和筛选"`
	Page     int    `d:"1" json:"page" query:"page" dc:"当前页码，默认1"`
	Size     int    `d:"10" json:"size" query:"size" dc:"每页记录数，默认10"`
	Id       int    `json:"id" query:"id" dc:"记录ID，可选筛选条件"`
	ClientId int    `json:"clientId" query:"clientId" dc:"客户ID，可选筛选条件"`
	Status   string `json:"status" query:"status" dc:"状态，可选筛选条件：processing、completed、failed"`
}

// ExchangeRecordListRes 获取兑换记录列表响应
type ExchangeRecordListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ExchangeRecordItem `json:"list" dc:"兑换记录列表"`
	Total  int                  `json:"total" dc:"总记录数"`
	Page   int                  `json:"page" dc:"当前页码"`
	Size   int                  `json:"size" dc:"每页记录数"`
	Pages  int                  `json:"pages" dc:"总页数"`
}

// ExchangeRecordItem 兑换记录项
type ExchangeRecordItem struct {
	Id              int         `json:"id" dc:"兑换记录ID"`
	ClientId        int         `json:"clientId" dc:"客户ID"`
	ClientName      string      `json:"clientName" dc:"客户名称"`
	RechargeAccount string      `json:"rechargeAccount" dc:"充值账号"`
	ProductName     string      `json:"productName" dc:"商品名称"`
	Duration        int         `json:"duration" dc:"消耗时长(分钟)"`
	ExchangeTime    *gtime.Time `json:"exchangeTime" dc:"兑换时间"`
	Status          string      `json:"status" dc:"状态"`
	Remark          string      `json:"remark" dc:"备注"`
	CreatedAt       *gtime.Time `json:"createdAt" dc:"创建时间"`
}

// WxExchangeRecordPageReq 微信客户端分页获取兑换记录列表请求
type WxExchangeRecordPageReq struct {
	g.Meta `path:"/client/exchange-record/page" method:"get" tags:"客户端兑换记录" summary:"客户端分页获取兑换记录列表" security:"Bearer" description:"客户端分页获取当前登录客户的兑换记录列表"`
	Page   int `d:"1" json:"page" query:"page" dc:"当前页码，默认1"`
	Size   int `d:"10" json:"size" query:"size" dc:"每页记录数，默认10"`
}

// WxExchangeRecordPageRes 微信客户端分页获取兑换记录列表响应
type WxExchangeRecordPageRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ExchangeRecordItem `json:"list" dc:"兑换记录列表"`
	Total  int                  `json:"total" dc:"总记录数"`
	Page   int                  `json:"page" dc:"当前页码"`
	Size   int                  `json:"size" dc:"每页记录数"`
	Pages  int                  `json:"pages" dc:"总页数"`
}

// WxExchangeRecordCreateReq 微信客户端创建兑换记录请求
type WxExchangeRecordCreateReq struct {
	g.Meta          `path:"/client/exchange-record/create" method:"post" tags:"客户端兑换记录" summary:"客户端创建兑换记录" security:"Bearer" description:"客户端创建新的兑换记录"`
	RechargeAccount string `v:"required#充值账号不能为空" json:"rechargeAccount" dc:"充值账号"`
	ProductName     string `v:"required#商品名称不能为空" json:"productName" dc:"商品名称"`
	Duration        int    `v:"required#消耗时长不能为空" json:"duration" dc:"消耗时长(分钟)"`
}

// WxExchangeRecordCreateRes 微信客户端创建兑换记录响应
type WxExchangeRecordCreateRes struct {
	g.Meta `mime:"application/json" example:"json"`
	Id     int `json:"id" dc:"新创建的兑换记录ID"`
}

// WxExchangeRecordPublicListReq 微信客户端获取最新兑换记录列表请求（公开）
type WxExchangeRecordPublicListReq struct {
	g.Meta `path:"/wx/exchange-record/list" method:"get" tags:"客户端兑换记录" summary:"客户端获取最新兑换记录列表" description:"获取最新的兑换记录列表，客户名称中间4位用*号屏蔽"`
	Limit  int `d:"10" json:"limit" query:"limit" dc:"获取记录数量，默认10条"`
}

// WxExchangeRecordPublicListRes 微信客户端获取最新兑换记录列表响应（公开）
type WxExchangeRecordPublicListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []ExchangeRecordPublicItem `json:"list" dc:"兑换记录列表"`
}

// ExchangeRecordPublicItem 兑换记录项（公开展示，客户名称中间屏蔽）
type ExchangeRecordPublicItem struct {
	ClientName   string      `json:"clientName" dc:"客户名称（中间4位用*号屏蔽）"`
	ProductName  string      `json:"productName" dc:"商品名称"`
	ExchangeTime *gtime.Time `json:"exchangeTime" dc:"兑换时间"`
}
