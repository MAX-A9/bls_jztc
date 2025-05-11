package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// AdSettingsReq 获取广告设置请求
type AdSettingsReq struct {
	g.Meta `path:"/ad/settings" method:"get" tags:"广告设置" summary:"获取广告设置" security:"Bearer" description:"获取广告设置信息，需要管理员权限"`
}

// AdSettingsRes 获取广告设置响应
type AdSettingsRes struct {
	g.Meta            `mime:"application/json" example:"json"`
	EnableWxAd        bool   `json:"enableWxAd" dc:"是否启用微信广告"`
	RewardedVideoAdId string `json:"rewardedVideoAdId" dc:"激励视频广告位ID"`
}

// AdSettingsSaveReq 保存广告设置请求
type AdSettingsSaveReq struct {
	g.Meta            `path:"/ad/settings/save" method:"post" tags:"广告设置" summary:"保存广告设置" security:"Bearer" description:"保存广告设置信息，需要管理员权限"`
	EnableWxAd        bool   `json:"enableWxAd" dc:"是否启用微信广告"`
	RewardedVideoAdId string `json:"rewardedVideoAdId" v:"required#激励视频广告位ID不能为空" dc:"激励视频广告位ID"`
}

// AdSettingsSaveRes 保存广告设置响应
type AdSettingsSaveRes struct {
	g.Meta    `mime:"application/json" example:"json"`
	IsSuccess bool `json:"isSuccess" dc:"是否成功"`
}

// WxGetAdSettingsReq 微信客户端获取广告设置请求
type WxGetAdSettingsReq struct {
	g.Meta `path:"/wx/ad/settings" method:"get" tags:"客户端基础设置" summary:"获取广告设置" description:"微信客户端获取广告设置信息"`
}

// WxGetAdSettingsRes 微信客户端获取广告设置响应
type WxGetAdSettingsRes struct {
	g.Meta            `mime:"application/json" example:"json"`
	EnableWxAd        bool   `json:"enableWxAd" dc:"是否启用微信广告"`
	RewardedVideoAdId string `json:"rewardedVideoAdId" dc:"激励视频广告位ID"`
}

// WxRewardAdViewedReq 微信客户端广告观看完成请求
type WxRewardAdViewedReq struct {
	g.Meta `path:"/ad/reward/viewed" method:"post" tags:"客户端基础设置" summary:"广告观看完成" security:"Bearer" description:"微信客户端上报广告观看完成，获取奖励"`
}

// WxRewardAdViewedRes 微信客户端广告观看完成响应
type WxRewardAdViewedRes struct {
	g.Meta         `mime:"application/json" example:"json"`
	Success        bool    `json:"success" dc:"是否成功"`
	RewardMinutes  int     `json:"rewardMinutes" dc:"奖励时长(分钟)"`
	RewardDays     float64 `json:"rewardDays" dc:"奖励时长(天)"`
	ExpirationDate string  `json:"expirationDate" dc:"奖励过期日期，格式：YYYY-MM-DD"`
	Message        string  `json:"message" dc:"提示消息"`
}
