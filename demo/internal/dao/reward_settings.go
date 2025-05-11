package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RewardSettingsDao 广告奖励设置数据访问对象
type RewardSettingsDao struct{}

// Get 获取奖励设置
func (d *RewardSettingsDao) Get(ctx context.Context) (*entity.RewardSettings, error) {
	var settings entity.RewardSettings
	err := g.DB().Model("reward_settings").Where("id=?", 1).Scan(&settings)
	// 如果没有设置数据，初始化一个空对象
	if err != nil || settings.Id == 0 {
		return &entity.RewardSettings{
			Id:                     1,
			EnableReward:           false,
			FirstViewMinRewardMin:  60,  // 默认60分钟
			FirstViewMaxRewardDay:  7.0, // 默认7天
			SingleAdMinRewardMin:   5,   // 默认5分钟
			SingleAdMaxRewardDay:   1.0, // 默认1天
			DailyRewardLimit:       10,  // 默认10次
			DailyMaxAccumulatedDay: 2.0, // 默认2天
			RewardExpirationDays:   30,  // 默认30天
		}, nil
	}
	return &settings, nil
}

// Save 保存奖励设置
func (d *RewardSettingsDao) Save(ctx context.Context, settings *entity.RewardSettings) error {
	now := gtime.Now()
	// 检查是否已有记录
	var count int
	count, err := g.DB().Model("reward_settings").Where("id=?", 1).Count()
	if err != nil {
		return err
	}

	// 如果已存在记录则更新，否则创建
	if count > 0 {
		_, err = g.DB().Model("reward_settings").
			Data(g.Map{
				"enable_reward":             settings.EnableReward,
				"first_view_min_reward_min": settings.FirstViewMinRewardMin,
				"first_view_max_reward_day": settings.FirstViewMaxRewardDay,
				"single_ad_min_reward_min":  settings.SingleAdMinRewardMin,
				"single_ad_max_reward_day":  settings.SingleAdMaxRewardDay,
				"daily_reward_limit":        settings.DailyRewardLimit,
				"daily_max_accumulated_day": settings.DailyMaxAccumulatedDay,
				"reward_expiration_days":    settings.RewardExpirationDays,
				"updated_at":                now,
			}).
			Where("id=?", 1).
			Update()
		return err
	} else {
		_, err = g.DB().Model("reward_settings").
			Data(g.Map{
				"id":                        1,
				"enable_reward":             settings.EnableReward,
				"first_view_min_reward_min": settings.FirstViewMinRewardMin,
				"first_view_max_reward_day": settings.FirstViewMaxRewardDay,
				"single_ad_min_reward_min":  settings.SingleAdMinRewardMin,
				"single_ad_max_reward_day":  settings.SingleAdMaxRewardDay,
				"daily_reward_limit":        settings.DailyRewardLimit,
				"daily_max_accumulated_day": settings.DailyMaxAccumulatedDay,
				"reward_expiration_days":    settings.RewardExpirationDays,
				"created_at":                now,
				"updated_at":                now,
			}).
			Insert()
		return err
	}
}
