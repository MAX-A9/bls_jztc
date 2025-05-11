package content

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "demo/api/content/v1"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/service"
)

type sCategory struct{}

// NewCategory 创建分类服务实例
func NewCategory() service.CategoryService {
	return &sCategory{}
}

// HomeList 获取首页分类列表
func (s *sCategory) HomeList(ctx context.Context, req *v1.HomeCategoryListReq) (res *v1.HomeCategoryListRes, err error) {
	res = &v1.HomeCategoryListRes{
		List: make([]v1.CategoryItem, 0),
	}

	// 查询数据
	homeCategoryDao := dao.NewHomeCategoryDao()
	list, err := homeCategoryDao.FindList(ctx)
	if err != nil {
		return nil, gerror.New("获取首页分类列表失败: " + err.Error())
	}

	// 转换数据
	for _, item := range list {
		categoryItem := v1.CategoryItem{
			Id:        gconv.Int(item.Id),
			Name:      gconv.String(item.Name),
			SortOrder: gconv.Int(item.SortOrder),
			IsActive:  gconv.Bool(item.IsActive),
			Icon:      gconv.String(item.Icon),
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		}
		res.List = append(res.List, categoryItem)
	}

	return res, nil
}

// HomeCreate 创建首页分类
func (s *sCategory) HomeCreate(ctx context.Context, req *v1.HomeCategoryCreateReq) (res *v1.HomeCategoryCreateRes, err error) {
	res = &v1.HomeCategoryCreateRes{}

	// 构建数据
	data := &do.HomeCategoryDO{
		CategoryDO: do.CategoryDO{
			Name:      req.Name,
			SortOrder: req.SortOrder,
			IsActive:  req.IsActive,
			Icon:      req.Icon,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		},
	}

	// 插入数据
	homeCategoryDao := dao.NewHomeCategoryDao()
	lastInsertId, err := homeCategoryDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建首页分类失败: " + err.Error())
	}

	res.Id = int(lastInsertId)
	return res, nil
}

// HomeUpdate 更新首页分类
func (s *sCategory) HomeUpdate(ctx context.Context, req *v1.HomeCategoryUpdateReq) (res *v1.HomeCategoryUpdateRes, err error) {
	res = &v1.HomeCategoryUpdateRes{}

	// 查询原有数据
	homeCategoryDao := dao.NewHomeCategoryDao()
	exists, err := homeCategoryDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询首页分类失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("首页分类不存在")
	}

	// 构建更新数据
	data := &do.HomeCategoryDO{
		CategoryDO: do.CategoryDO{
			Name:      req.Name,
			SortOrder: req.SortOrder,
			IsActive:  req.IsActive,
			Icon:      req.Icon,
			UpdatedAt: gtime.Now(),
		},
	}

	// 更新数据
	_, err = homeCategoryDao.Update(ctx, data, req.Id)
	if err != nil {
		return nil, gerror.New("更新首页分类失败: " + err.Error())
	}

	return res, nil
}

// HomeDelete 删除首页分类
func (s *sCategory) HomeDelete(ctx context.Context, req *v1.HomeCategoryDeleteReq) (res *v1.HomeCategoryDeleteRes, err error) {
	res = &v1.HomeCategoryDeleteRes{}

	// 查询原有数据
	homeCategoryDao := dao.NewHomeCategoryDao()
	exists, err := homeCategoryDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询首页分类失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("首页分类不存在")
	}

	// 删除数据（软删除）
	_, err = homeCategoryDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除首页分类失败: " + err.Error())
	}

	return res, nil
}

// IdleList 获取闲置分类列表
func (s *sCategory) IdleList(ctx context.Context, req *v1.IdleCategoryListReq) (res *v1.IdleCategoryListRes, err error) {
	res = &v1.IdleCategoryListRes{
		List: make([]v1.CategoryItem, 0),
	}

	// 查询数据
	idleCategoryDao := dao.NewIdleCategoryDao()
	list, err := idleCategoryDao.FindList(ctx)
	if err != nil {
		return nil, gerror.New("获取闲置分类列表失败: " + err.Error())
	}

	// 转换数据
	for _, item := range list {
		categoryItem := v1.CategoryItem{
			Id:        gconv.Int(item.Id),
			Name:      gconv.String(item.Name),
			SortOrder: gconv.Int(item.SortOrder),
			IsActive:  gconv.Bool(item.IsActive),
			Icon:      gconv.String(item.Icon),
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		}
		res.List = append(res.List, categoryItem)
	}

	return res, nil
}

// IdleCreate 创建闲置分类
func (s *sCategory) IdleCreate(ctx context.Context, req *v1.IdleCategoryCreateReq) (res *v1.IdleCategoryCreateRes, err error) {
	res = &v1.IdleCategoryCreateRes{}

	// 构建数据
	data := &do.IdleCategoryDO{
		CategoryDO: do.CategoryDO{
			Name:      req.Name,
			SortOrder: req.SortOrder,
			IsActive:  req.IsActive,
			Icon:      req.Icon,
			CreatedAt: gtime.Now(),
			UpdatedAt: gtime.Now(),
		},
	}

	// 插入数据
	idleCategoryDao := dao.NewIdleCategoryDao()
	lastInsertId, err := idleCategoryDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建闲置分类失败: " + err.Error())
	}

	res.Id = int(lastInsertId)
	return res, nil
}

// IdleUpdate 更新闲置分类
func (s *sCategory) IdleUpdate(ctx context.Context, req *v1.IdleCategoryUpdateReq) (res *v1.IdleCategoryUpdateRes, err error) {
	res = &v1.IdleCategoryUpdateRes{}

	// 查询原有数据
	idleCategoryDao := dao.NewIdleCategoryDao()
	exists, err := idleCategoryDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询闲置分类失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("闲置分类不存在")
	}

	// 构建更新数据
	data := &do.IdleCategoryDO{
		CategoryDO: do.CategoryDO{
			Name:      req.Name,
			SortOrder: req.SortOrder,
			IsActive:  req.IsActive,
			Icon:      req.Icon,
			UpdatedAt: gtime.Now(),
		},
	}

	// 更新数据
	_, err = idleCategoryDao.Update(ctx, data, req.Id)
	if err != nil {
		return nil, gerror.New("更新闲置分类失败: " + err.Error())
	}

	return res, nil
}

// IdleDelete 删除闲置分类
func (s *sCategory) IdleDelete(ctx context.Context, req *v1.IdleCategoryDeleteReq) (res *v1.IdleCategoryDeleteRes, err error) {
	res = &v1.IdleCategoryDeleteRes{}

	// 查询原有数据
	idleCategoryDao := dao.NewIdleCategoryDao()
	exists, err := idleCategoryDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("查询闲置分类失败: " + err.Error())
	}
	if exists == nil {
		return nil, gerror.New("闲置分类不存在")
	}

	// 删除数据（软删除）
	_, err = idleCategoryDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除闲置分类失败: " + err.Error())
	}

	return res, nil
}
