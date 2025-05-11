package dao

import (
	"context"
	"strings"

	"demo/internal/model"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FileDao 文件数据访问对象
type FileDao struct{}

// New 创建DAO对象
func (dao *FileDao) New() *FileDao {
	return &FileDao{}
}

// Model 获取ORM模型
func (dao *FileDao) Model(ctx context.Context) *gdb.Model {
	// 使用Unscoped禁用软删除功能
	return g.DB().Model(model.TableFile).Safe().Unscoped().Ctx(ctx)
}

// Insert 插入文件记录
func (dao *FileDao) Insert(ctx context.Context, data *model.FileDO) (id int64, err error) {
	// 确保ID字段为0，让数据库自动生成
	data.Id = 0

	result, err := dao.Model(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Delete 删除文件记录
func (dao *FileDao) Delete(ctx context.Context, id int) (affected int64, err error) {
	result, err := dao.Model(ctx).Where("id", id).Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// FindById 根据ID查询文件
func (dao *FileDao) FindById(ctx context.Context, id int) (file *model.FileDO, err error) {
	err = dao.Model(ctx).Where("id", id).Scan(&file)
	return
}

// FindByPath 根据路径查询文件
func (dao *FileDao) FindByPath(ctx context.Context, path string) (file *model.FileDO, err error) {
	err = dao.Model(ctx).Where("path", path).Scan(&file)
	return
}

// FindList 查询文件列表
func (dao *FileDao) FindList(ctx context.Context, filter map[string]interface{}, page, size int) (list []*model.FileDO, total int, err error) {
	model := dao.Model(ctx)

	// 添加查询条件
	if v, ok := filter["keyword"]; ok && v != "" {
		keyword := v.(string)
		model = model.Where("name LIKE ? OR path LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if v, ok := filter["type"]; ok && v != "" {
		model = model.Where("type", v)
	}

	if v, ok := filter["isPublic"]; ok {
		model = model.Where("is_public", v)
	}

	if v, ok := filter["userId"]; ok && v != 0 {
		model = model.Where("user_id", v)
	}

	// 扩展名筛选
	if v, ok := filter["extension"]; ok && v != "" {
		extensions := strings.Split(v.(string), ",")
		if len(extensions) > 0 {
			model = model.WhereIn("extension", extensions)
		}
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

// UpdatePublicStatus 更新文件公开状态
func (dao *FileDao) UpdatePublicStatus(ctx context.Context, id int, isPublic bool) (affected int64, err error) {
	result, err := dao.Model(ctx).Data(g.Map{
		"is_public": isPublic,
	}).Where("id", id).Update()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// BatchDelete 批量删除文件
func (dao *FileDao) BatchDelete(ctx context.Context, ids []int) (affected int64, err error) {
	result, err := dao.Model(ctx).WhereIn("id", ids).Delete()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
