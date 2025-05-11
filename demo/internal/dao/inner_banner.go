package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// InnerBannerDao 内页轮播图数据访问对象
type InnerBannerDao struct{}

// GetList 获取内页轮播图列表
// 参数onlyEnabled为true表示只获取已启用的轮播图，false则获取所有轮播图
func (d *InnerBannerDao) GetList(ctx context.Context, bannerType string, onlyEnabled bool) ([]*entity.InnerPageBanner, error) {
	model := g.DB().Model("inner_page_banner").
		Where("banner_type", bannerType).
		OrderAsc("order")

	if onlyEnabled {
		model = model.Where("is_enabled=?", 1)
	}

	var banners []*entity.InnerPageBanner
	err := model.Scan(&banners)

	// 如果是"没有找到行"错误，将其视为正常情况，返回空数组
	if err != nil && err.Error() == "sql: no rows in result set" {
		return []*entity.InnerPageBanner{}, nil
	}

	return banners, err
}

// GetById 根据ID获取内页轮播图
func (d *InnerBannerDao) GetById(ctx context.Context, id int) (*entity.InnerPageBanner, error) {
	var banner entity.InnerPageBanner
	err := g.DB().Model("inner_page_banner").Where("id=?", id).Scan(&banner)
	return &banner, err
}

// Create 创建内页轮播图
func (d *InnerBannerDao) Create(ctx context.Context, banner *entity.InnerPageBanner) (int64, error) {
	result, err := g.DB().Model("inner_page_banner").Insert(g.Map{
		"image":       banner.Image,
		"link_type":   banner.LinkType,
		"link_url":    banner.LinkUrl,
		"is_enabled":  banner.IsEnabled,
		"order":       banner.Order,
		"banner_type": banner.BannerType,
		"created_at":  gtime.Now(),
		"updated_at":  gtime.Now(),
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新内页轮播图
func (d *InnerBannerDao) Update(ctx context.Context, banner *entity.InnerPageBanner) error {
	_, err := g.DB().Model("inner_page_banner").
		Data(g.Map{
			"image":       banner.Image,
			"link_type":   banner.LinkType,
			"link_url":    banner.LinkUrl,
			"is_enabled":  banner.IsEnabled,
			"order":       banner.Order,
			"banner_type": banner.BannerType,
			"updated_at":  gtime.Now(),
		}).
		Where("id=?", banner.Id).
		Update()
	return err
}

// UpdateStatus 更新内页轮播图状态
func (d *InnerBannerDao) UpdateStatus(ctx context.Context, id int, isEnabled int) error {
	_, err := g.DB().Model("inner_page_banner").
		Data(g.Map{
			"is_enabled": isEnabled,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", id).
		Update()
	return err
}

// Delete 删除内页轮播图
func (d *InnerBannerDao) Delete(ctx context.Context, id int) error {
	_, err := g.DB().Model("inner_page_banner").Where("id=?", id).Delete()
	return err
}
