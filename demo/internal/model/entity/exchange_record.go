package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExchangeRecord 兑换记录表
type ExchangeRecord struct {
	Id              int         `json:"id"             description:"兑换记录ID"`
	ClientId        int         `json:"clientId"       description:"客户ID"`
	ClientName      string      `json:"clientName"     description:"客户名称"`
	RechargeAccount string      `json:"rechargeAccount" description:"充值账号"`
	ProductName     string      `json:"productName"    description:"商品名称"`
	Duration        int         `json:"duration"       description:"消耗时长(分钟)"`
	ExchangeTime    *gtime.Time `json:"exchangeTime"   description:"兑换时间"`
	Status          string      `json:"status"         description:"状态：processing(处理中)、completed(已完成)、failed(失败)"`
	Remark          string      `json:"remark"         description:"备注"`
	CreatedAt       *gtime.Time `json:"createdAt"      description:"创建时间"`
	UpdatedAt       *gtime.Time `json:"updatedAt"      description:"更新时间"`
}
