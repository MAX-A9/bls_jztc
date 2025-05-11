package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdSettingsDao 广告设置数据访问对象
type AdSettingsDao struct{}

// Get 获取广告设置
func (d *AdSettingsDao) Get(ctx context.Context) (*entity.AdSettings, error) {
	var settings entity.AdSettings
	err := g.DB().Model("ad_settings").Where("id=?", 1).Scan(&settings)
	// 如果没有设置数据，初始化一个空对象
	if err != nil || settings.Id == 0 {
		return &entity.AdSettings{
			Id:                1,
			EnableWxAd:        false,
			RewardedVideoAdId: "",
		}, nil
	}
	return &settings, nil
}

// Save 保存广告设置
func (d *AdSettingsDao) Save(ctx context.Context, settings *entity.AdSettings) error {
	now := gtime.Now()
	// 检查是否已有记录
	var count int
	count, err := g.DB().Model("ad_settings").Where("id=?", 1).Count()
	if err != nil {
		return err
	}

	// 如果已存在记录则更新，否则创建
	if count > 0 {
		_, err = g.DB().Model("ad_settings").
			Data(g.Map{
				"enable_wx_ad":         settings.EnableWxAd,
				"rewarded_video_ad_id": settings.RewardedVideoAdId,
				"updated_at":           now,
			}).
			Where("id=?", 1).
			Update()
		return err
	} else {
		_, err = g.DB().Model("ad_settings").
			Data(g.Map{
				"id":                   1,
				"enable_wx_ad":         settings.EnableWxAd,
				"rewarded_video_ad_id": settings.RewardedVideoAdId,
				"created_at":           now,
				"updated_at":           now,
			}).
			Insert()
		return err
	}
}
