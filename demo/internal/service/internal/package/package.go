package package_service

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "demo/api/package/v1"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/service"
)

// 实现服务接口
type sPackageImpl struct{}

// New 创建套餐服务实例
func New() service.PackageService {
	return &sPackageImpl{}
}

// 初始化，注册服务实例
func init() {
	service.SetPackage(New())
}

// List 获取套餐列表
func (s *sPackageImpl) List(ctx context.Context, req *v1.PackageListReq) (res *v1.PackageListRes, err error) {
	// 初始化响应
	res = &v1.PackageListRes{
		List: make([]*v1.Package, 0),
	}

	// 获取系统配置DAO和包DAO
	systemConfigDao := &dao.SystemConfigDao{}
	packageDao := dao.NewPackageDao()
	var packages []*do.PackageDO

	// 根据类型筛选套餐
	if req.Type == v1.PackageTypeTop {
		// 获取置顶套餐总开关状态
		topEnabled, err := systemConfigDao.GetTopPackageEnabled(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取置顶套餐总开关状态失败", err)
			return nil, gerror.New("获取置顶套餐总开关状态失败")
		}
		res.IsGlobalEnabled = topEnabled

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
		if err != nil {
			g.Log().Error(ctx, "获取置顶套餐列表失败", err)
			return nil, gerror.New("获取套餐列表失败")
		}
	} else if req.Type == v1.PackageTypePublish {
		// 获取发布套餐总开关状态
		publishEnabled, err := systemConfigDao.GetPublishPackageEnabled(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取发布套餐总开关状态失败", err)
			return nil, gerror.New("获取发布套餐总开关状态失败")
		}
		res.IsGlobalEnabled = publishEnabled

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
		if err != nil {
			g.Log().Error(ctx, "获取发布套餐列表失败", err)
			return nil, gerror.New("获取套餐列表失败")
		}
	} else {
		// 获取所有套餐
		packages, err = packageDao.FindAll(ctx, req.Sort, req.Order)
		if err != nil {
			g.Log().Error(ctx, "获取套餐列表失败", err)
			return nil, gerror.New("获取套餐列表失败")
		}
	}

	// 转换数据
	for _, p := range packages {
		if p == nil {
			continue
		}

		// 创建套餐对象
		packageItem := &v1.Package{
			Id:           p.Id.(int),
			Title:        p.Title.(string),
			Description:  p.Description.(string),
			Price:        p.Price.(float64),
			Type:         v1.PackageType(p.Type.(string)),
			Duration:     p.Duration.(int),
			DurationType: v1.DurationType(p.DurationType.(string)),
			SortOrder:    p.SortOrder.(int),
		}

		// 设置时间
		if p.CreatedAt != nil {
			packageItem.CreatedAt = p.CreatedAt.Format("2006-01-02 15:04:05")
		}
		if p.UpdatedAt != nil {
			packageItem.UpdatedAt = p.UpdatedAt.Format("2006-01-02 15:04:05")
		}

		res.List = append(res.List, packageItem)
	}

	return res, nil
}

// Detail 获取套餐详情
func (s *sPackageImpl) Detail(ctx context.Context, req *v1.PackageDetailReq) (res *v1.PackageDetailRes, err error) {
	// 查询数据
	packageDO, err := dao.NewPackageDao().FindOne(ctx, req.Id)
	if err != nil {
		g.Log().Error(ctx, "获取套餐详情失败", err)
		return nil, gerror.New("获取套餐详情失败")
	}

	if packageDO == nil {
		return nil, gerror.New("套餐不存在")
	}

	// 构建响应
	res = &v1.PackageDetailRes{
		Package: &v1.Package{
			Id:           packageDO.Id.(int),
			Title:        packageDO.Title.(string),
			Description:  packageDO.Description.(string),
			Price:        packageDO.Price.(float64),
			Type:         v1.PackageType(packageDO.Type.(string)),
			Duration:     packageDO.Duration.(int),
			DurationType: v1.DurationType(packageDO.DurationType.(string)),
			SortOrder:    packageDO.SortOrder.(int),
		},
	}

	// 设置时间
	if packageDO.CreatedAt != nil {
		res.Package.CreatedAt = packageDO.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if packageDO.UpdatedAt != nil {
		res.Package.UpdatedAt = packageDO.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return res, nil
}

// Create 创建套餐
func (s *sPackageImpl) Create(ctx context.Context, req *v1.PackageCreateReq) (res *v1.PackageCreateRes, err error) {
	// 创建数据
	now := gtime.Now()
	data := &do.PackageDO{
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		Type:         string(req.Type),
		Duration:     req.Duration,
		DurationType: string(req.DurationType),
		SortOrder:    req.SortOrder,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 插入数据
	id, err := dao.NewPackageDao().Insert(ctx, data)
	if err != nil {
		g.Log().Error(ctx, "创建套餐失败", err)
		return nil, gerror.New("创建套餐失败")
	}

	// 构建响应
	res = &v1.PackageCreateRes{
		Id: int(id),
	}

	return res, nil
}

// Update 更新套餐
func (s *sPackageImpl) Update(ctx context.Context, req *v1.PackageUpdateReq) (res *v1.PackageUpdateRes, err error) {
	// 查询是否存在
	packageDO, err := dao.NewPackageDao().FindOne(ctx, req.Id)
	if err != nil {
		g.Log().Error(ctx, "更新套餐查询失败", err)
		return nil, gerror.New("更新套餐失败")
	}

	if packageDO == nil {
		return nil, gerror.New("套餐不存在")
	}

	// 更新数据
	data := &do.PackageDO{
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		Type:         string(req.Type),
		Duration:     req.Duration,
		DurationType: string(req.DurationType),
		SortOrder:    req.SortOrder,
		UpdatedAt:    gtime.Now(),
	}

	// 执行更新
	_, err = dao.NewPackageDao().Update(ctx, data, req.Id)
	if err != nil {
		g.Log().Error(ctx, "更新套餐失败", err)
		return nil, gerror.New("更新套餐失败")
	}

	// 构建响应
	res = &v1.PackageUpdateRes{}

	return res, nil
}

// Delete 删除套餐
func (s *sPackageImpl) Delete(ctx context.Context, req *v1.PackageDeleteReq) (res *v1.PackageDeleteRes, err error) {
	// 查询是否存在
	packageDO, err := dao.NewPackageDao().FindOne(ctx, req.Id)
	if err != nil {
		g.Log().Error(ctx, "删除套餐查询失败", err)
		return nil, gerror.New("删除套餐失败")
	}

	if packageDO == nil {
		return nil, gerror.New("套餐不存在")
	}

	// 执行删除
	_, err = dao.NewPackageDao().Delete(ctx, req.Id)
	if err != nil {
		g.Log().Error(ctx, "删除套餐失败", err)
		return nil, gerror.New("删除套餐失败")
	}

	// 构建响应
	res = &v1.PackageDeleteRes{}

	return res, nil
}

// GetGlobalStatus 获取套餐总开关状态
func (s *sPackageImpl) GetGlobalStatus(ctx context.Context, req *v1.PackageGlobalStatusReq) (res *v1.PackageGlobalStatusRes, err error) {
	res = &v1.PackageGlobalStatusRes{}

	// 获取系统配置DAO
	systemConfigDao := &dao.SystemConfigDao{}

	// 获取置顶套餐总开关状态
	topEnabled, err := systemConfigDao.GetTopPackageEnabled(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取置顶套餐总开关状态失败", err)
		return nil, gerror.New("获取置顶套餐总开关状态失败")
	}
	res.TopEnabled = topEnabled

	// 获取发布套餐总开关状态
	publishEnabled, err := systemConfigDao.GetPublishPackageEnabled(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取发布套餐总开关状态失败", err)
		return nil, gerror.New("获取发布套餐总开关状态失败")
	}
	res.PublishEnabled = publishEnabled

	return res, nil
}

// UpdateTopPackageGlobalStatus 更新置顶套餐总开关状态
func (s *sPackageImpl) UpdateTopPackageGlobalStatus(ctx context.Context, req *v1.TopPackageGlobalStatusUpdateReq) (res *v1.TopPackageGlobalStatusUpdateRes, err error) {
	res = &v1.TopPackageGlobalStatusUpdateRes{}

	// 获取系统配置DAO
	systemConfigDao := &dao.SystemConfigDao{}

	// 更新置顶套餐总开关状态
	if err := systemConfigDao.UpdateTopPackageEnabled(ctx, req.IsEnabled); err != nil {
		g.Log().Error(ctx, "更新置顶套餐总开关状态失败", err)
		return nil, gerror.New("更新置顶套餐总开关状态失败")
	}

	return res, nil
}

// UpdatePublishPackageGlobalStatus 更新发布套餐总开关状态
func (s *sPackageImpl) UpdatePublishPackageGlobalStatus(ctx context.Context, req *v1.PublishPackageGlobalStatusUpdateReq) (res *v1.PublishPackageGlobalStatusUpdateRes, err error) {
	res = &v1.PublishPackageGlobalStatusUpdateRes{}

	// 获取系统配置DAO
	systemConfigDao := &dao.SystemConfigDao{}

	// 更新发布套餐总开关状态
	if err := systemConfigDao.UpdatePublishPackageEnabled(ctx, req.IsEnabled); err != nil {
		g.Log().Error(ctx, "更新发布套餐总开关状态失败", err)
		return nil, gerror.New("更新发布套餐总开关状态失败")
	}

	return res, nil
}

// WxList 客户端获取套餐列表
func (s *sPackageImpl) WxList(ctx context.Context, req *v1.WxPackageListReq) (res *v1.WxPackageListRes, err error) {
	res = &v1.WxPackageListRes{
		List: make([]*v1.Package, 0),
	}

	// 获取系统配置DAO
	systemConfigDao := &dao.SystemConfigDao{}
	packageDao := dao.NewPackageDao()
	var packages []*do.PackageDO

	// 根据类型筛选套餐
	if req.Type == v1.PackageTypeTop {
		// 获取置顶套餐总开关状态
		topEnabled, err := systemConfigDao.GetTopPackageEnabled(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取置顶套餐总开关状态失败", err)
			return nil, gerror.New("获取置顶套餐总开关状态失败")
		}
		res.IsGlobalEnabled = topEnabled

		// 如果总开关已关闭，则直接返回空列表
		if !topEnabled {
			return res, nil
		}

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
		if err != nil {
			g.Log().Error(ctx, "获取置顶套餐列表失败", err)
			return nil, gerror.New("获取套餐列表失败")
		}
	} else if req.Type == v1.PackageTypePublish {
		// 获取发布套餐总开关状态
		publishEnabled, err := systemConfigDao.GetPublishPackageEnabled(ctx)
		if err != nil {
			g.Log().Error(ctx, "获取发布套餐总开关状态失败", err)
			return nil, gerror.New("获取发布套餐总开关状态失败")
		}
		res.IsGlobalEnabled = publishEnabled

		// 如果总开关已关闭，则直接返回空列表
		if !publishEnabled {
			return res, nil
		}

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
		if err != nil {
			g.Log().Error(ctx, "获取发布套餐列表失败", err)
			return nil, gerror.New("获取套餐列表失败")
		}
	} else {
		// 获取所有套餐
		packages, err = packageDao.FindAll(ctx, req.Sort, req.Order)
	}

	if err != nil {
		g.Log().Error(ctx, "获取套餐列表失败", err)
		return nil, gerror.New("获取套餐列表失败")
	}

	// 转换数据
	for _, p := range packages {
		if p == nil {
			continue
		}

		// 创建套餐对象
		packageItem := &v1.Package{
			Id:           p.Id.(int),
			Title:        p.Title.(string),
			Description:  p.Description.(string),
			Price:        p.Price.(float64),
			Type:         v1.PackageType(p.Type.(string)),
			Duration:     p.Duration.(int),
			DurationType: v1.DurationType(p.DurationType.(string)),
			SortOrder:    p.SortOrder.(int),
		}

		// 设置时间
		if p.CreatedAt != nil {
			packageItem.CreatedAt = p.CreatedAt.Format("2006-01-02 15:04:05")
		}
		if p.UpdatedAt != nil {
			packageItem.UpdatedAt = p.UpdatedAt.Format("2006-01-02 15:04:05")
		}

		res.List = append(res.List, packageItem)
	}

	return res, nil
}
