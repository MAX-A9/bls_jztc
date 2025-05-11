package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ShareSettingsDao 分享设置数据访问对象
type ShareSettingsDao struct{}

// Get 获取分享设置
func (d *ShareSettingsDao) Get(ctx context.Context) (*entity.ShareSettings, error) {
	var settings entity.ShareSettings
	err := g.DB().Model("share_settings").Where("id=?", 1).Scan(&settings)
	// 如果没有设置数据，返回错误，让上层处理
	if err != nil || settings.Id == 0 {
		g.Log().Warning(ctx, "分享设置数据不存在，请先在管理后台配置")
		return nil, gerror.New("分享设置数据不存在")
	}
	return &settings, nil
}

// Save 保存分享设置
func (d *ShareSettingsDao) Save(ctx context.Context, settings *entity.ShareSettings) error {
	now := gtime.Now()
	// 检查是否已有记录
	var count int
	count, err := g.DB().Model("share_settings").Where("id=?", 1).Count()
	if err != nil {
		return err
	}

	// 如果已存在记录则更新，否则创建
	if count > 0 {
		_, err = g.DB().Model("share_settings").
			Data(g.Map{
				"default_share_text":  settings.DefaultShareText,
				"default_share_image": settings.DefaultShareImage,
				"content_share_text":  settings.ContentShareText,
				"content_share_image": settings.ContentShareImage,
				"home_share_text":     settings.HomeShareText,
				"home_share_image":    settings.HomeShareImage,
				"updated_at":          now,
			}).
			Where("id=?", 1).
			Update()
		return err
	} else {
		_, err = g.DB().Model("share_settings").
			Data(g.Map{
				"id":                  1,
				"default_share_text":  settings.DefaultShareText,
				"default_share_image": settings.DefaultShareImage,
				"content_share_text":  settings.ContentShareText,
				"content_share_image": settings.ContentShareImage,
				"home_share_text":     settings.HomeShareText,
				"home_share_image":    settings.HomeShareImage,
				"created_at":          now,
				"updated_at":          now,
			}).
			Insert()
		return err
	}
}
