package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// RewardSettingsReq 获取奖励设置请求
type RewardSettingsReq struct {
	g.Meta `path:"/reward/settings" method:"get" tags:"广告设置" summary:"获取奖励设置" security:"Bearer" description:"获取广告奖励设置信息，需要管理员权限"`
}

// RewardSettingsRes 获取奖励设置响应
type RewardSettingsRes struct {
	g.Meta                 `mime:"application/json" example:"json"`
	EnableReward           bool    `json:"enableReward" dc:"是否启用奖励功能"`
	FirstViewMinRewardMin  int     `json:"firstViewMinRewardMin" dc:"首次观看广告最小奖励(分钟)"`
	FirstViewMaxRewardDay  float64 `json:"firstViewMaxRewardDay" dc:"首次观看广告最大奖励(天)"`
	SingleAdMinRewardMin   int     `json:"singleAdMinRewardMin" dc:"单次广告最小奖励(分钟)"`
	SingleAdMaxRewardDay   float64 `json:"singleAdMaxRewardDay" dc:"单次广告最大奖励(天)"`
	DailyRewardLimit       int     `json:"dailyRewardLimit" dc:"每日奖励次数上限"`
	DailyMaxAccumulatedDay float64 `json:"dailyMaxAccumulatedDay" dc:"每日最大累计奖励(天)"`
	RewardExpirationDays   int     `json:"rewardExpirationDays" dc:"奖励过期天数"`
}

// RewardSettingsSaveReq 保存奖励设置请求
type RewardSettingsSaveReq struct {
	g.Meta                 `path:"/reward/settings/save" method:"post" tags:"广告设置" summary:"保存奖励设置" security:"Bearer" description:"保存广告奖励设置信息，需要管理员权限"`
	EnableReward           bool    `json:"enableReward" dc:"是否启用奖励功能"`
	FirstViewMinRewardMin  int     `json:"firstViewMinRewardMin" v:"required|min:1#首次观看最小奖励分钟不能为空|首次观看最小奖励必须大于0" dc:"首次观看广告最小奖励(分钟)"`
	FirstViewMaxRewardDay  float64 `json:"firstViewMaxRewardDay" v:"required|min:0.1#首次观看最大奖励天数不能为空|首次观看最大奖励必须大于0" dc:"首次观看广告最大奖励(天)"`
	SingleAdMinRewardMin   int     `json:"singleAdMinRewardMin" v:"required|min:1#单次广告最小奖励分钟不能为空|单次广告最小奖励必须大于0" dc:"单次广告最小奖励(分钟)"`
	SingleAdMaxRewardDay   float64 `json:"singleAdMaxRewardDay" v:"required|min:0.1#单次广告最大奖励天数不能为空|单次广告最大奖励必须大于0" dc:"单次广告最大奖励(天)"`
	DailyRewardLimit       int     `json:"dailyRewardLimit" v:"required|min:1#每日奖励次数上限不能为空|每日奖励次数上限必须大于0" dc:"每日奖励次数上限"`
	DailyMaxAccumulatedDay float64 `json:"dailyMaxAccumulatedDay" v:"required|min:0.1#每日最大累计奖励不能为空|每日最大累计奖励必须大于0" dc:"每日最大累计奖励(天)"`
	RewardExpirationDays   int     `json:"rewardExpirationDays" v:"required|min:1#奖励过期天数不能为空|奖励过期天数必须大于0" dc:"奖励过期天数"`
}

// RewardSettingsSaveRes 保存奖励设置响应
type RewardSettingsSaveRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	IsSuccess bool `json:"isSuccess" dc:"是否成功"`
}
