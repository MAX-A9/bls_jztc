package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"demo/internal/model/do"
)

// HomeCategoryDao 首页分类表DAO接口
type HomeCategoryDao interface {
	// Insert 插入一条数据
	Insert(ctx context.Context, data *do.HomeCategoryDO) (lastInsertId int64, err error)
	// Update 更新数据
	Update(ctx context.Context, data *do.HomeCategoryDO, id interface{}) (rowsAffected int64, err error)
	// Delete 删除数据
	Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error)
	// FindOne 查询单条数据
	FindOne(ctx context.Context, id interface{}) (*do.HomeCategoryDO, error)
	// FindList 查询列表数据
	FindList(ctx context.Context) (list []*do.HomeCategoryDO, err error)
}

// homeCategoryDao 首页分类表DAO实现
type homeCategoryDao struct{}

// NewHomeCategoryDao 创建首页分类表DAO
func NewHomeCategoryDao() HomeCategoryDao {
	return &homeCategoryDao{}
}

// Insert 插入一条数据
func (d *homeCategoryDao) Insert(ctx context.Context, data *do.HomeCategoryDO) (lastInsertId int64, err error) {
	result, err := g.DB().Model(do.TableHomeCategory).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新数据
func (d *homeCategoryDao) Update(ctx context.Context, data *do.HomeCategoryDO, id interface{}) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TableHomeCategory).Ctx(ctx).Data(data).Where("id", id).Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (d *homeCategoryDao) Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error) {
	// 软删除
	result, err := g.DB().Model(do.TableHomeCategory).Ctx(ctx).
		Data(do.HomeCategoryDO{
			CategoryDO: do.CategoryDO{
				DeletedAt: gtime.Now(),
			},
		}).
		Where("id", id).
		Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindOne 查询单条数据
func (d *homeCategoryDao) FindOne(ctx context.Context, id interface{}) (*do.HomeCategoryDO, error) {
	var category *do.HomeCategoryDO
	err := g.DB().Model(do.TableHomeCategory).Ctx(ctx).
		Where("id", id).
		Where("deleted_at IS NULL").
		Scan(&category)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}
	return category, nil
}

// FindList 查询列表数据
func (d *homeCategoryDao) FindList(ctx context.Context) (list []*do.HomeCategoryDO, err error) {
	var categories []*do.HomeCategoryDO
	err = g.DB().Model(do.TableHomeCategory).Ctx(ctx).
		Where("deleted_at IS NULL").
		Order("sort_order ASC").
		Scan(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// IdleCategoryDao 闲置分类表DAO接口
type IdleCategoryDao interface {
	// Insert 插入一条数据
	Insert(ctx context.Context, data *do.IdleCategoryDO) (lastInsertId int64, err error)
	// Update 更新数据
	Update(ctx context.Context, data *do.IdleCategoryDO, id interface{}) (rowsAffected int64, err error)
	// Delete 删除数据
	Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error)
	// FindOne 查询单条数据
	FindOne(ctx context.Context, id interface{}) (*do.IdleCategoryDO, error)
	// FindList 查询列表数据
	FindList(ctx context.Context) (list []*do.IdleCategoryDO, err error)
}

// idleCategoryDao 闲置分类表DAO实现
type idleCategoryDao struct{}

// NewIdleCategoryDao 创建闲置分类表DAO
func NewIdleCategoryDao() IdleCategoryDao {
	return &idleCategoryDao{}
}

// Insert 插入一条数据
func (d *idleCategoryDao) Insert(ctx context.Context, data *do.IdleCategoryDO) (lastInsertId int64, err error) {
	result, err := g.DB().Model(do.TableIdleCategory).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新数据
func (d *idleCategoryDao) Update(ctx context.Context, data *do.IdleCategoryDO, id interface{}) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TableIdleCategory).Ctx(ctx).Data(data).Where("id", id).Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (d *idleCategoryDao) Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error) {
	// 软删除
	result, err := g.DB().Model(do.TableIdleCategory).Ctx(ctx).
		Data(do.IdleCategoryDO{
			CategoryDO: do.CategoryDO{
				DeletedAt: gtime.Now(),
			},
		}).
		Where("id", id).
		Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindOne 查询单条数据
func (d *idleCategoryDao) FindOne(ctx context.Context, id interface{}) (*do.IdleCategoryDO, error) {
	var category *do.IdleCategoryDO
	err := g.DB().Model(do.TableIdleCategory).Ctx(ctx).
		Where("id", id).
		Where("deleted_at IS NULL").
		Scan(&category)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, nil
	}
	return category, nil
}

// FindList 查询列表数据
func (d *idleCategoryDao) FindList(ctx context.Context) (list []*do.IdleCategoryDO, err error) {
	var categories []*do.IdleCategoryDO
	err = g.DB().Model(do.TableIdleCategory).Ctx(ctx).
		Where("deleted_at IS NULL").
		Order("sort_order ASC").
		Scan(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
