package internal

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/frame/g"

	"demo/internal/model/do"
)

// ButlerDao 专属管家数据访问对象
type ButlerDao struct{}

// NewButlerDao 创建专属管家DAO
func NewButlerDao() *ButlerDao {
	return &ButlerDao{}
}

// SaveImage 保存专属管家图片
func (d *ButlerDao) SaveImage(ctx context.Context, data *do.ButlerDO) (int64, error) {
	r, err := g.DB().Model("butler").Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

// GetLatestImage 获取最新的专属管家图片
func (d *ButlerDao) GetLatestImage(ctx context.Context) (*do.ButlerDO, error) {
	var butler *do.ButlerDO

	err := g.DB().Model("butler").
		Ctx(ctx).
		Where("status", 1). // 状态为启用
		Order("created_at DESC").
		Limit(1).
		Scan(&butler)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return butler, nil
}
