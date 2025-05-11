package dao

import (
	"context"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BottomTabDao 底部导航栏数据访问对象
type BottomTabDao struct{}

// GetList 获取底部导航栏列表
// 参数onlyEnabled为true表示只获取已启用的导航项，false则获取所有导航项
func (d *BottomTabDao) GetList(ctx context.Context, onlyEnabled bool) ([]*entity.BottomTab, error) {
	model := g.DB().Model("bottom_tab").OrderAsc("order")

	if onlyEnabled {
		model = model.Where("is_enabled=?", 1)
	}

	var tabs []*entity.BottomTab
	err := model.Scan(&tabs)

	// 如果是"没有找到行"错误，将其视为正常情况，返回空数组
	if err != nil && err.Error() == "sql: no rows in result set" {
		return []*entity.BottomTab{}, nil
	}

	return tabs, err
}

// GetById 根据ID获取底部导航项
func (d *BottomTabDao) GetById(ctx context.Context, id int) (*entity.BottomTab, error) {
	var tab entity.BottomTab
	err := g.DB().Model("bottom_tab").Where("id=?", id).Scan(&tab)
	return &tab, err
}

// Create 创建底部导航项
func (d *BottomTabDao) Create(ctx context.Context, tab *entity.BottomTab) (int64, error) {
	result, err := g.DB().Model("bottom_tab").Insert(g.Map{
		"name":          tab.Name,
		"icon":          tab.Icon,
		"selected_icon": tab.SelectedIcon,
		"path":          tab.Path,
		"order":         tab.Order,
		"is_enabled":    tab.IsEnabled,
		"created_at":    gtime.Now(),
		"updated_at":    gtime.Now(),
	})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新底部导航项
func (d *BottomTabDao) Update(ctx context.Context, tab *entity.BottomTab) error {
	_, err := g.DB().Model("bottom_tab").
		Data(g.Map{
			"name":          tab.Name,
			"icon":          tab.Icon,
			"selected_icon": tab.SelectedIcon,
			"path":          tab.Path,
			"order":         tab.Order,
			"is_enabled":    tab.IsEnabled,
			"updated_at":    gtime.Now(),
		}).
		Where("id=?", tab.Id).
		Update()
	return err
}

// UpdateStatus 更新底部导航项状态
func (d *BottomTabDao) UpdateStatus(ctx context.Context, id int, isEnabled int) error {
	_, err := g.DB().Model("bottom_tab").
		Data(g.Map{
			"is_enabled": isEnabled,
			"updated_at": gtime.Now(),
		}).
		Where("id=?", id).
		Update()
	return err
}

// Delete 删除底部导航项
func (d *BottomTabDao) Delete(ctx context.Context, id int) error {
	_, err := g.DB().Model("bottom_tab").Where("id=?", id).Delete()
	return err
}
