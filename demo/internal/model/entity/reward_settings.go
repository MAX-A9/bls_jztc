package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RewardSettings 广告奖励设置表
type RewardSettings struct {
	Id                     int         `json:"id"                   description:"设置ID"`
	EnableReward           bool        `json:"enableReward"         description:"是否启用奖励功能"`
	FirstViewMinRewardMin  int         `json:"firstViewMinRewardMin"    description:"首次观看广告最小奖励(分钟)"`
	FirstViewMaxRewardDay  float64     `json:"firstViewMaxRewardDay"    description:"首次观看广告最大奖励(天)"`
	SingleAdMinRewardMin   int         `json:"singleAdMinRewardMin"     description:"单次广告最小奖励(分钟)"`
	SingleAdMaxRewardDay   float64     `json:"singleAdMaxRewardDay"     description:"单次广告最大奖励(天)"`
	DailyRewardLimit       int         `json:"dailyRewardLimit"         description:"每日奖励次数上限"`
	DailyMaxAccumulatedDay float64     `json:"dailyMaxAccumulatedDay"   description:"每日最大累计奖励(天)"`
	RewardExpirationDays   int         `json:"rewardExpirationDays"     description:"奖励过期天数"`
	CreatedAt              *gtime.Time `json:"createdAt"             description:"创建时间"`
	UpdatedAt              *gtime.Time `json:"updatedAt"             description:"更新时间"`
}
