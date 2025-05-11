package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MiniProgramSettingsDao 小程序基础设置数据访问对象
type MiniProgramSettingsDao struct{}

// Get 获取小程序基础设置
func (d *MiniProgramSettingsDao) Get(ctx context.Context) (*entity.MiniProgramSettings, error) {
	var settings entity.MiniProgramSettings
	err := g.DB().Model("mini_program_settings").Where("id=?", 1).Scan(&settings)
	// 如果没有设置数据，初始化一个空对象
	if err != nil || settings.Id == 0 {
		return &entity.MiniProgramSettings{
			Id:          1,
			Name:        "",
			Description: "",
			Logo:        "",
		}, nil
	}
	return &settings, nil
}

// Save 保存小程序基础设置
func (d *MiniProgramSettingsDao) Save(ctx context.Context, settings *entity.MiniProgramSettings) error {
	now := gtime.Now()
	// 检查是否已有记录
	var count int
	count, err := g.DB().Model("mini_program_settings").Where("id=?", 1).Count()
	if err != nil {
		return err
	}

	// 如果已存在记录则更新，否则创建
	if count > 0 {
		_, err = g.DB().Model("mini_program_settings").
			Data(g.Map{
				"name":        settings.Name,
				"description": settings.Description,
				"logo":        settings.Logo,
				"updated_at":  now,
			}).
			Where("id=?", 1).
			Update()
		return err
	} else {
		_, err = g.DB().Model("mini_program_settings").
			Data(g.Map{
				"id":          1,
				"name":        settings.Name,
				"description": settings.Description,
				"logo":        settings.Logo,
				"created_at":  now,
				"updated_at":  now,
			}).
			Insert()
		return err
	}
}
