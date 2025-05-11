package dao

import (
	"context"
	"demo/internal/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserDao 用户数据访问对象
type UserDao struct{}

// New 创建DAO对象
func (dao *UserDao) New() *UserDao {
	return &UserDao{}
}

// Model 获取ORM模型
func (dao *UserDao) Model(ctx context.Context) *gdb.Model {
	return g.DB().Model(do.TableUser).Safe().Ctx(ctx)
}

// FindOne 根据条件查询单条数据
func (dao *UserDao) FindOne(ctx context.Context, where ...interface{}) (user *do.UserDO, err error) {
	err = dao.Model(ctx).Where(where).Scan(&user)
	return
}

// FindById 根据ID查询
func (dao *UserDao) FindById(ctx context.Context, id int) (user *do.UserDO, err error) {
	err = dao.Model(ctx).Where("id", id).Scan(&user)
	return
}

// FindByUsername 根据用户名查询
func (dao *UserDao) FindByUsername(ctx context.Context, username string) (user *do.UserDO, err error) {
	err = dao.Model(ctx).Where("username", username).Scan(&user)
	return
}

// FindList 查询列表
func (dao *UserDao) FindList(ctx context.Context, filter map[string]interface{}, page, size int) (list []*do.UserDO, total int, err error) {
	model := dao.Model(ctx)

	// 添加查询条件
	if v, ok := filter["username"]; ok && v != "" {
		model = model.WhereLike("username", "%"+v.(string)+"%")
	}
	if v, ok := filter["nickname"]; ok && v != "" {
		model = model.WhereLike("nickname", "%"+v.(string)+"%")
	}
	if v, ok := filter["status"]; ok && v != "" {
		model = model.Where("status", v)
	}

	// 查询总数
	total, err = model.Count()
	if err != nil {
		return
	}

	// 查询数据列表
	err = model.Page(page, size).Order("id DESC").Scan(&list)
	return
}

// Insert 插入数据
func (dao *UserDao) Insert(ctx context.Context, data *do.UserDO) (id int64, err error) {
	result, err := dao.Model(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新数据
func (dao *UserDao) Update(ctx context.Context, data *do.UserDO, id int) (affected int64, err error) {
	result, err := dao.Model(ctx).Data(data).Where("id", id).Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据
func (dao *UserDao) Delete(ctx context.Context, id int) (affected int64, err error) {
	result, err := dao.Model(ctx).Where("id", id).Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
