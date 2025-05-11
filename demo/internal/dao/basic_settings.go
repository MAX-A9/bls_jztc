package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MiniProgramDao 导航小程序数据访问对象
type MiniProgramDao struct{}

// GetList 获取导航小程序列表
// 参数onlyEnabled为true表示只获取已启用的小程序，false则获取所有小程序
// 参数checkGlobal为true表示检查总开关状态，false则忽略总开关状态
func (d *MiniProgramDao) GetList(ctx context.Context, onlyEnabled bool, checkGlobal ...bool) ([]*entity.MiniProgram, error) {
	// 默认检查总开关状态
	needCheckGlobal := true
	if len(checkGlobal) > 0 && !checkGlobal[0] {
		needCheckGlobal = false
	}

	// 如果需要检查总开关，则先检查总开关状态
	if needCheckGlobal && onlyEnabled {
		systemConfigDao := &SystemConfigDao{}
		enabled, err := systemConfigDao.GetMiniProgramEnabled(ctx)
		if err != nil {
			return nil, err
		}

		// 如果总开关已关闭且请求只查询已启用的小程序，则直接返回空列表
		if !enabled {
			return []*entity.MiniProgram{}, nil
		}
	}

	model := g.DB().Model("mini_program").OrderAsc("order")
	if onlyEnabled {
		model = model.Where("is_enabled=?", 1)
	}
	var miniPrograms []*entity.MiniProgram
	err := model.Scan(&miniPrograms)
	return miniPrograms, err
}

// GetById 根据ID获取导航小程序
func (d *MiniProgramDao) GetById(ctx context.Context, id int) (*entity.MiniProgram, error) {
	var miniProgram entity.MiniProgram
	err := g.DB().Model("mini_program").Where("id=?", id).Scan(&miniProgram)
	return &miniProgram, err
}

// Create 创建导航小程序
func (d *MiniProgramDao) Create(ctx context.Context, miniProgram *entity.MiniProgram) (int64, error) {
	result, err := g.DB().Model("mini_program").Insert(g.Map{
		"name":       miniProgram.Name,
		"app_id":     miniProgram.AppId,
		"logo":       miniProgram.Logo,
		"is_enabled": miniProgram.IsEnabled,
		"order":      miniProgram.Order,
		"created_at": gtime.Now(),
		"updated_at": gtime.Now(),
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新导航小程序
func (d *MiniProgramDao) Update(ctx context.Context, miniProgram *entity.MiniProgram) error {
	_, err := g.DB().Model("mini_program").
		Data(g.Map{
			"name":       miniProgram.Name,
			"app_id":     miniProgram.AppId,
			"logo":       miniProgram.Logo,
			"is_enabled": miniProgram.IsEnabled,
			"order":      miniProgram.Order,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", miniProgram.Id).
		Update()
	return err
}

// UpdateStatus 更新导航小程序状态
func (d *MiniProgramDao) UpdateStatus(ctx context.Context, id int, isEnabled int) error {
	_, err := g.DB().Model("mini_program").
		Data(g.Map{
			"is_enabled": isEnabled,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", id).
		Update()
	return err
}

// Delete 删除导航小程序
func (d *MiniProgramDao) Delete(ctx context.Context, id int) error {
	_, err := g.DB().Model("mini_program").Where("id=?", id).Delete()
	return err
}

// BannerDao 轮播图数据访问对象
type BannerDao struct{}

// GetList 获取轮播图列表
// 参数onlyEnabled为true表示只获取已启用的轮播图，false则获取所有轮播图
// 参数checkGlobal为true表示检查总开关状态，false则忽略总开关状态
func (d *BannerDao) GetList(ctx context.Context, onlyEnabled bool, checkGlobal ...bool) ([]*entity.Banner, error) {
	// 默认检查总开关状态
	needCheckGlobal := true
	if len(checkGlobal) > 0 && !checkGlobal[0] {
		needCheckGlobal = false
	}

	// 如果需要检查总开关，则先检查总开关状态
	if needCheckGlobal && onlyEnabled {
		systemConfigDao := &SystemConfigDao{}
		enabled, err := systemConfigDao.GetBannerEnabled(ctx)
		if err != nil {
			return nil, err
		}

		// 如果总开关已关闭且请求只查询已启用的轮播图，则直接返回空列表
		if !enabled {
			return []*entity.Banner{}, nil
		}
	}

	model := g.DB().Model("banner").OrderAsc("order")
	if onlyEnabled {
		model = model.Where("is_enabled=?", 1)
	}
	var banners []*entity.Banner
	err := model.Scan(&banners)
	return banners, err
}

// GetById 根据ID获取轮播图
func (d *BannerDao) GetById(ctx context.Context, id int) (*entity.Banner, error) {
	var banner entity.Banner
	err := g.DB().Model("banner").Where("id=?", id).Scan(&banner)
	return &banner, err
}

// Create 创建轮播图
func (d *BannerDao) Create(ctx context.Context, banner *entity.Banner) (int64, error) {
	result, err := g.DB().Model("banner").Insert(g.Map{
		"image":      banner.Image,
		"link_type":  banner.LinkType,
		"link_url":   banner.LinkUrl,
		"is_enabled": banner.IsEnabled,
		"order":      banner.Order,
		"created_at": gtime.Now(),
		"updated_at": gtime.Now(),
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新轮播图
func (d *BannerDao) Update(ctx context.Context, banner *entity.Banner) error {
	_, err := g.DB().Model("banner").
		Data(g.Map{
			"image":      banner.Image,
			"link_type":  banner.LinkType,
			"link_url":   banner.LinkUrl,
			"is_enabled": banner.IsEnabled,
			"order":      banner.Order,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", banner.Id).
		Update()
	return err
}

// UpdateStatus 更新轮播图状态
func (d *BannerDao) UpdateStatus(ctx context.Context, id int, isEnabled int) error {
	_, err := g.DB().Model("banner").
		Data(g.Map{
			"is_enabled": isEnabled,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", id).
		Update()
	return err
}

// Delete 删除轮播图
func (d *BannerDao) Delete(ctx context.Context, id int) error {
	_, err := g.DB().Model("banner").Where("id=?", id).Delete()
	return err
}

// ActivityAreaDao 活动区域数据访问对象
type ActivityAreaDao struct{}

// Get 获取活动区域
func (d *ActivityAreaDao) Get(ctx context.Context) (*entity.ActivityArea, error) {
	var activityArea entity.ActivityArea
	err := g.DB().Model("activity_area").Where("id=?", 1).Scan(&activityArea)
	return &activityArea, err
}

// Save 保存活动区域
func (d *ActivityAreaDao) Save(ctx context.Context, data *entity.ActivityArea) error {
	_, err := g.DB().Model("activity_area").
		Data(g.Map{
			"top_left_title":          data.TopLeftTitle,
			"top_left_description":    data.TopLeftDescription,
			"top_left_link_type":      data.TopLeftLinkType,
			"top_left_link_url":       data.TopLeftLinkUrl,
			"bottom_left_title":       data.BottomLeftTitle,
			"bottom_left_description": data.BottomLeftDescription,
			"bottom_left_link_type":   data.BottomLeftLinkType,
			"bottom_left_link_url":    data.BottomLeftLinkUrl,
			"right_title":             data.RightTitle,
			"right_description":       data.RightDescription,
			"right_link_type":         data.RightLinkType,
			"right_link_url":          data.RightLinkUrl,
			"updated_at":              gtime.Now(),
		}).
		Where("id=?", 1).
		Update()
	return err
}
