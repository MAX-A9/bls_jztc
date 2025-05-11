package v1

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// IRewardController 奖励记录控制器接口
type IRewardController interface {
	RewardRecordList(ctx context.Context, req *RewardRecordListReq) (res *RewardRecordListRes, err error)
	RewardRecordDetail(ctx context.Context, req *RewardRecordDetailReq) (res *RewardRecordDetailRes, err error)
	RewardRecordStat(ctx context.Context, req *RewardRecordStatReq) (res *RewardRecordStatRes, err error)
}

// 获取奖励记录列表请求
type RewardRecordListReq struct {
	g.Meta   `path:"" method:"get" tags:"奖励记录" summary:"获取奖励记录列表" security:"Bearer" description:"获取奖励记录列表，需要管理员权限"`
	Page     int `v:"min:1#页码最小值为1" json:"page" dc:"页码"`
	PageSize int `v:"max:100#每页最大100条" json:"pageSize" dc:"每页数量"`
	Status   int `json:"status" dc:"状态：0-已过期，1-有效，不传则查询全部"`
	ClientId int `json:"clientId" dc:"客户ID，不传则查询全部客户"`
}

// 奖励记录列表响应
type RewardRecordListRes struct {
	g.Meta `mime:"application/json" example:"json"`
	List   []RewardRecordListItem `json:"list" dc:"奖励记录列表"`
	Total  int                    `json:"total" dc:"总数量"`
	Page   int                    `json:"page" dc:"当前页码"`
}

// 奖励记录列表项
type RewardRecordListItem struct {
	Id                 int         `json:"id" dc:"记录ID"`
	ClientId           int         `json:"clientId" dc:"客户ID"`
	ClientName         string      `json:"clientName" dc:"客户名称"`
	RewardMinutes      int         `json:"rewardMinutes" dc:"奖励时长(分钟)"`
	RewardDays         float64     `json:"rewardDays" dc:"奖励时长(天)"`
	IsFirstView        bool        `json:"isFirstView" dc:"是否首次观看：false-否，true-是"`
	RemainingMinutes   int         `json:"remainingMinutes" dc:"剩余时长(分钟)"`
	TotalRewardMinutes int         `json:"totalRewardMinutes" dc:"累计获得奖励(分钟)"`
	UsedMinutes        int         `json:"usedMinutes" dc:"已使用时长(分钟)"`
	Status             int         `json:"status" dc:"状态：0-已过期，1-有效"`
	StatusText         string      `json:"statusText" dc:"状态文本：已过期/有效"`
	ExpireAt           *gtime.Time `json:"expireAt" dc:"过期时间"`
	CreatedAt          *gtime.Time `json:"createdAt" dc:"创建时间"`
}

// 获取奖励记录详情请求
type RewardRecordDetailReq struct {
	g.Meta `path:"" method:"get" tags:"奖励记录" summary:"获取奖励记录详情" security:"Bearer" description:"获取指定ID的奖励记录详情，需要管理员权限"`
	Id     int `v:"required#记录ID不能为空" json:"id" dc:"记录ID"`
}

// 奖励记录详情响应
type RewardRecordDetailRes struct {
	g.Meta             `mime:"application/json" example:"json"`
	Id                 int         `json:"id" dc:"记录ID"`
	ClientId           int         `json:"clientId" dc:"客户ID"`
	ClientName         string      `json:"clientName" dc:"客户名称"`
	RewardMinutes      int         `json:"rewardMinutes" dc:"奖励时长(分钟)"`
	RewardDays         float64     `json:"rewardDays" dc:"奖励时长(天)"`
	IsFirstView        bool        `json:"isFirstView" dc:"是否首次观看：false-否，true-是"`
	RemainingMinutes   int         `json:"remainingMinutes" dc:"剩余时长(分钟)"`
	TotalRewardMinutes int         `json:"totalRewardMinutes" dc:"累计获得奖励(分钟)"`
	UsedMinutes        int         `json:"usedMinutes" dc:"已使用时长(分钟)"`
	Status             int         `json:"status" dc:"状态：0-已过期，1-有效"`
	StatusText         string      `json:"statusText" dc:"状态文本：已过期/有效"`
	ExpireAt           *gtime.Time `json:"expireAt" dc:"过期时间"`
	CreatedAt          *gtime.Time `json:"createdAt" dc:"创建时间"`
	UpdatedAt          *gtime.Time `json:"updatedAt" dc:"更新时间"`
}

// 获取奖励记录统计请求
type RewardRecordStatReq struct {
	g.Meta   `path:"" method:"get" tags:"奖励记录" summary:"获取奖励记录统计" security:"Bearer" description:"获取奖励记录统计信息，需要管理员权限"`
	ClientId int `json:"clientId" dc:"客户ID，不传则统计所有客户"`
}

// 奖励记录统计响应
type RewardRecordStatRes struct {
	g.Meta         `mime:"application/json" example:"json"`
	TotalMinutes   int     `json:"totalMinutes" dc:"有效奖励总分钟数"`
	TotalDays      float64 `json:"totalDays" dc:"有效奖励总天数"`
	TodayCount     int     `json:"todayCount" dc:"今日已领取总次数"`
	RemainingCount int     `json:"remainingCount" dc:"平均剩余可领取次数"`
	ClientCount    int     `json:"clientCount" dc:"有奖励记录的客户数量"`
	ExpiredCount   int     `json:"expiredCount" dc:"已过期奖励记录数量"`
	ValidCount     int     `json:"validCount" dc:"有效奖励记录数量"`
}
