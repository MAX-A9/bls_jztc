package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RewardRecord 奖励记录表
type RewardRecord struct {
	Id                 int         `json:"id"                   description:"记录ID"`
	ClientId           int         `json:"clientId"             description:"客户ID"`
	RewardMinutes      int         `json:"rewardMinutes"        description:"奖励时长(分钟)"`
	RewardDays         float64     `json:"rewardDays"           description:"奖励时长(天)"`
	IsFirstView        bool        `json:"isFirstView"          description:"是否首次观看：0-否，1-是"`
	RemainingMinutes   int         `json:"remainingMinutes"     description:"剩余时长(分钟)"`
	TotalRewardMinutes int         `json:"totalRewardMinutes"   description:"累计获得奖励(分钟)"`
	UsedMinutes        int         `json:"usedMinutes"          description:"已使用时长(分钟)"`
	Status             int         `json:"status"               description:"状态：0-已过期，1-有效"`
	ExpireAt           *gtime.Time `json:"expireAt"             description:"过期时间"`
	CreatedAt          *gtime.Time `json:"createdAt"            description:"创建时间"`
	UpdatedAt          *gtime.Time `json:"updatedAt"            description:"更新时间"`
}
