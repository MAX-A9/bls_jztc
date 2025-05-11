package dao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"demo/internal/model/do"
)

// PackageDao 套餐表DAO接口
type PackageDao interface {
	// Insert 插入一条数据
	Insert(ctx context.Context, data *do.PackageDO) (lastInsertId int64, err error)
	// InsertWithId 插入一条数据，包含ID
	InsertWithId(ctx context.Context, data *do.PackageDO) (rowsAffected int64, err error)
	// Update 更新数据
	Update(ctx context.Context, data *do.PackageDO, id interface{}) (rowsAffected int64, err error)
	// Delete 删除数据
	Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error)
	// FindOne 查询单条数据
	FindOne(ctx context.Context, id interface{}) (*do.PackageDO, error)
	// FindByType 根据类型查询数据
	FindByType(ctx context.Context, packageType string, sort string, order string) (list []*do.PackageDO, err error)
	// FindAll 查询所有数据
	FindAll(ctx context.Context, sort string, order string) (list []*do.PackageDO, err error)
}

// packageDao 套餐表DAO实现
type packageDao struct{}

// NewPackageDao 创建套餐表DAO
func NewPackageDao() PackageDao {
	return &packageDao{}
}

// Insert 插入一条数据
func (d *packageDao) Insert(ctx context.Context, data *do.PackageDO) (lastInsertId int64, err error) {
	result, err := g.DB().Model(do.TablePackage).Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertWithId 插入一条数据，包含ID
func (d *packageDao) InsertWithId(ctx context.Context, data *do.PackageDO) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TablePackage).Ctx(ctx).Data(data).Save()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Update 更新数据
func (d *packageDao) Update(ctx context.Context, data *do.PackageDO, id interface{}) (rowsAffected int64, err error) {
	result, err := g.DB().Model(do.TablePackage).Ctx(ctx).Data(data).Where("id", id).Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (d *packageDao) Delete(ctx context.Context, id interface{}) (rowsAffected int64, err error) {
	// 软删除
	result, err := g.DB().Model(do.TablePackage).Ctx(ctx).
		Data(do.PackageDO{
			DeletedAt: gtime.Now(),
		}).
		Where("id", id).
		Where("deleted_at IS NULL").
		Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindOne 查询单条数据
func (d *packageDao) FindOne(ctx context.Context, id interface{}) (*do.PackageDO, error) {
	var packageDO *do.PackageDO
	err := g.DB().Model(do.TablePackage).Ctx(ctx).
		Where("id", id).
		Where("deleted_at IS NULL").
		Scan(&packageDO)
	if err != nil {
		return nil, err
	}
	if packageDO == nil {
		return nil, nil
	}
	return packageDO, nil
}

// FindByType 根据类型查询数据
func (d *packageDao) FindByType(ctx context.Context, packageType string, sort string, order string) (list []*do.PackageDO, err error) {
	var packages []*do.PackageDO
	model := g.DB().Model(do.TablePackage).Ctx(ctx).
		Where("type", packageType).
		Where("deleted_at IS NULL")

	// 处理排序
	if sort != "" {
		// 设置排序方向
		orderDirection := "ASC"
		if order == "desc" {
			orderDirection = "DESC"
		}

		// 根据指定字段排序
		switch sort {
		case "price":
			model = model.Order(sort + " " + orderDirection)
		case "duration":
			model = model.Order("duration " + orderDirection)
		default:
			// 默认排序方式
			model = model.Order("duration_type ASC, duration ASC, price ASC")
		}
	} else {
		// 默认排序方式
		model = model.Order("duration_type ASC, duration ASC, price ASC")
	}

	err = model.Scan(&packages)
	if err != nil {
		return nil, err
	}
	return packages, nil
}

// FindAll 查询所有数据
func (d *packageDao) FindAll(ctx context.Context, sort string, order string) (list []*do.PackageDO, err error) {
	var packages []*do.PackageDO
	model := g.DB().Model(do.TablePackage).Ctx(ctx).
		Where("deleted_at IS NULL")

	// 处理排序
	if sort != "" {
		// 设置排序方向
		orderDirection := "ASC"
		if order == "desc" {
			orderDirection = "DESC"
		}

		// 根据指定字段排序
		switch sort {
		case "price":
			model = model.Order(sort + " " + orderDirection)
		case "duration":
			model = model.Order("duration " + orderDirection)
		default:
			// 默认排序方式
			model = model.Order("type ASC, duration_type ASC, duration ASC, price ASC")
		}
	} else {
		// 默认排序方式
		model = model.Order("type ASC, duration_type ASC, duration ASC, price ASC")
	}

	err = model.Scan(&packages)
	if err != nil {
		return nil, err
	}
	return packages, nil
}
