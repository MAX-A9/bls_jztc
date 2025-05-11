package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemConfigDao 系统配置数据访问对象
type SystemConfigDao struct{}

// GetByModuleAndKey 根据模块和键名获取配置
func (d *SystemConfigDao) GetByModuleAndKey(ctx context.Context, module, key string) (*entity.SystemConfig, error) {
	var config entity.SystemConfig
	err := g.DB().Model("system_config").
		Where("module=? AND `key`=?", module, key).
		Scan(&config)
	return &config, err
}

// UpdateValue 更新配置值
func (d *SystemConfigDao) UpdateValue(ctx context.Context, module, key, value string) error {
	_, err := g.DB().Model("system_config").
		Data(g.Map{
			"value":      value,
			"updated_at": gtime.Now(),
		}).
		Where("module=? AND `key`=?", module, key).
		Update()
	return err
}

// GetMiniProgramEnabled 获取导航小程序总开关状态
func (d *SystemConfigDao) GetMiniProgramEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "mini_program", "enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdateMiniProgramEnabled 更新导航小程序总开关状态
func (d *SystemConfigDao) UpdateMiniProgramEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "mini_program", "enabled", value)
}

// GetBannerEnabled 获取轮播图总开关状态
func (d *SystemConfigDao) GetBannerEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "banner", "enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdateBannerEnabled 更新轮播图总开关状态
func (d *SystemConfigDao) UpdateBannerEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "banner", "enabled", value)
}

// GetActivityAreaEnabled 获取活动区域总开关状态
func (d *SystemConfigDao) GetActivityAreaEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "activity_area", "enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdateActivityAreaEnabled 更新活动区域总开关状态
func (d *SystemConfigDao) UpdateActivityAreaEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "activity_area", "enabled", value)
}

// GetInnerBannerEnabled 获取内页轮播图总开关状态
func (d *SystemConfigDao) GetInnerBannerEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "inner_banner", "enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdateInnerBannerEnabled 更新内页轮播图总开关状态
func (d *SystemConfigDao) UpdateInnerBannerEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "inner_banner", "enabled", value)
}

// GetHomeInnerBannerEnabled 获取首页内页轮播图总开关状态
func (d *SystemConfigDao) GetHomeInnerBannerEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "inner_banner", "home_enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdateHomeInnerBannerEnabled 更新首页内页轮播图总开关状态
func (d *SystemConfigDao) UpdateHomeInnerBannerEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "inner_banner", "home_enabled", value)
}

// GetIdleInnerBannerEnabled 获取闲置页内页轮播图总开关状态
func (d *SystemConfigDao) GetIdleInnerBannerEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "inner_banner", "idle_enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdateIdleInnerBannerEnabled 更新闲置页内页轮播图总开关状态
func (d *SystemConfigDao) UpdateIdleInnerBannerEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "inner_banner", "idle_enabled", value)
}

// GetTopPackageEnabled 获取置顶套餐总开关状态
func (d *SystemConfigDao) GetTopPackageEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "package", "top_enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdateTopPackageEnabled 更新置顶套餐总开关状态
func (d *SystemConfigDao) UpdateTopPackageEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "package", "top_enabled", value)
}

// GetPublishPackageEnabled 获取发布套餐总开关状态
func (d *SystemConfigDao) GetPublishPackageEnabled(ctx context.Context) (bool, error) {
	config, err := d.GetByModuleAndKey(ctx, "package", "publish_enabled")
	if err != nil {
		return false, err
	}
	return config.Value == "true", nil
}

// UpdatePublishPackageEnabled 更新发布套餐总开关状态
func (d *SystemConfigDao) UpdatePublishPackageEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return d.UpdateValue(ctx, "package", "publish_enabled", value)
}
