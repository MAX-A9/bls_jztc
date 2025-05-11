package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdSettings 广告设置表
type AdSettings struct {
	Id                int         `json:"id"                description:"设置ID"`
	EnableWxAd        bool        `json:"enableWxAd"        description:"是否启用微信广告"`
	RewardedVideoAdId string      `json:"rewardedVideoAdId" description:"激励视频广告位ID"`
	CreatedAt         *gtime.Time `json:"createdAt"         description:"创建时间"`
	UpdatedAt         *gtime.Time `json:"updatedAt"         description:"更新时间"`
}
