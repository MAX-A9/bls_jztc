package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"demo/internal/model/do"
)

// ButlerDao 专属管家数据访问对象
type ButlerDao interface {
	// 保存专属管家图片
	SaveImage(ctx context.Context, data *do.ButlerDO) (int64, error)

	// 获取最新的专属管家图片
	GetLatestImage(ctx context.Context) (*do.ButlerDO, error)
}

// butlerDaoImpl 专属管家数据访问对象实现
type butlerDaoImpl struct{}

// NewButlerDao 创建专属管家DAO实例
func NewButlerDao() ButlerDao {
	return &butlerDaoImpl{}
}

// SaveImage 保存专属管家图片
func (d *butlerDaoImpl) SaveImage(ctx context.Context, data *do.ButlerDO) (int64, error) {
	// 先将所有管家图片状态设为禁用 - 添加WHERE条件或使用Safe(false)
	_, err := g.DB().Model(do.TableButler).
		Data(g.Map{"status": 0}).
		Where("status", 1). // 添加条件，只更新状态为1(启用)的记录
		Update()
	if err != nil {
		g.Log().Error(ctx, "禁用所有专属管家图片失败", err)
		return 0, err
	}

	// 插入新图片记录
	res, err := g.DB().Model(do.TableButler).Data(data).Insert()
	if err != nil {
		g.Log().Error(ctx, "插入专属管家图片记录失败", err)
		return 0, err
	}
	return res.LastInsertId()
}

// GetLatestImage 获取最新的专属管家图片
func (d *butlerDaoImpl) GetLatestImage(ctx context.Context) (*do.ButlerDO, error) {
	var butler *do.ButlerDO
	err := g.DB().Model(do.TableButler).
		Where("status", 1).
		Order("id DESC").
		Limit(1).
		Scan(&butler)
	if err != nil {
		g.Log().Error(ctx, "获取最新专属管家图片失败", err)
		return nil, err
	}
	return butler, nil
}
