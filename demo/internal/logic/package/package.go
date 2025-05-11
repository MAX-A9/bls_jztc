package package_logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	v1 "demo/api/package/v1"
	"demo/internal/dao"
	"demo/internal/model/do"
	"demo/internal/service"
	"demo/utility/auth"
)

type sPackage struct{}

// New 创建套餐服务实例
func New() service.PackageService {
	return &sPackage{}
}

// List 获取套餐列表
func (s *sPackage) List(ctx context.Context, req *v1.PackageListReq) (res *v1.PackageListRes, err error) {
	res = &v1.PackageListRes{
		List: make([]*v1.Package, 0),
	}

	systemConfigDao := &dao.SystemConfigDao{}
	packageDao := dao.NewPackageDao()
	var packages []*do.PackageDO

	// 根据类型筛选套餐
	if req.Type == v1.PackageTypeTop {
		// 获取置顶套餐总开关状态
		topEnabled, err := systemConfigDao.GetTopPackageEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取置顶套餐总开关状态失败: " + err.Error())
		}
		res.IsGlobalEnabled = topEnabled

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
	} else if req.Type == v1.PackageTypePublish {
		// 获取发布套餐总开关状态
		publishEnabled, err := systemConfigDao.GetPublishPackageEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取发布套餐总开关状态失败: " + err.Error())
		}
		res.IsGlobalEnabled = publishEnabled

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
	} else {
		// 获取所有套餐
		packages, err = packageDao.FindAll(ctx, req.Sort, req.Order)
	}

	if err != nil {
		return nil, gerror.New("获取套餐列表失败: " + err.Error())
	}

	// 转换数据
	for _, p := range packages {
		pkg := &v1.Package{
			Id:           gconv.Int(p.Id),
			Title:        gconv.String(p.Title),
			Description:  gconv.String(p.Description),
			Price:        gconv.Float64(p.Price),
			Type:         v1.PackageType(gconv.String(p.Type)),
			Duration:     gconv.Int(p.Duration),
			DurationType: v1.DurationType(gconv.String(p.DurationType)),
			SortOrder:    gconv.Int(p.SortOrder),
			CreatedAt:    p.CreatedAt.String(),
			UpdatedAt:    p.UpdatedAt.String(),
		}
		res.List = append(res.List, pkg)
	}

	return res, nil
}

// Detail 获取套餐详情
func (s *sPackage) Detail(ctx context.Context, req *v1.PackageDetailReq) (res *v1.PackageDetailRes, err error) {
	res = &v1.PackageDetailRes{}

	packageDao := dao.NewPackageDao()
	p, err := packageDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("获取套餐详情失败: " + err.Error())
	}
	if p == nil {
		return nil, gerror.New("套餐不存在")
	}

	// 转换数据
	res.Package = &v1.Package{
		Id:           gconv.Int(p.Id),
		Title:        gconv.String(p.Title),
		Description:  gconv.String(p.Description),
		Price:        gconv.Float64(p.Price),
		Type:         v1.PackageType(gconv.String(p.Type)),
		Duration:     gconv.Int(p.Duration),
		DurationType: v1.DurationType(gconv.String(p.DurationType)),
		SortOrder:    gconv.Int(p.SortOrder),
		CreatedAt:    p.CreatedAt.String(),
		UpdatedAt:    p.UpdatedAt.String(),
	}

	return res, nil
}

// Create 创建套餐
func (s *sPackage) Create(ctx context.Context, req *v1.PackageCreateReq) (res *v1.PackageCreateRes, err error) {
	res = &v1.PackageCreateRes{}

	// 校验套餐类型
	if req.Type != v1.PackageTypeTop && req.Type != v1.PackageTypePublish {
		return nil, gerror.New("套餐类型错误，只能是置顶套餐或发布套餐")
	}

	// 校验时长单位
	if req.DurationType != v1.DurationTypeHour && req.DurationType != v1.DurationTypeDay && req.DurationType != v1.DurationTypeMonth {
		return nil, gerror.New("时长单位错误，只能是小时、天或月")
	}

	// 构建数据
	data := &do.PackageDO{
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		Type:         string(req.Type),
		Duration:     req.Duration,
		DurationType: string(req.DurationType),
		SortOrder:    req.SortOrder,
		CreatedAt:    gtime.Now(),
		UpdatedAt:    gtime.Now(),
	}

	// 插入数据
	packageDao := dao.NewPackageDao()
	lastInsertId, err := packageDao.Insert(ctx, data)
	if err != nil {
		return nil, gerror.New("创建套餐失败: " + err.Error())
	}

	res.Id = int(lastInsertId)
	return res, nil
}

// Update 更新套餐
func (s *sPackage) Update(ctx context.Context, req *v1.PackageUpdateReq) (res *v1.PackageUpdateRes, err error) {
	res = &v1.PackageUpdateRes{}

	// 校验套餐类型
	if req.Type != v1.PackageTypeTop && req.Type != v1.PackageTypePublish {
		return nil, gerror.New("套餐类型错误，只能是置顶套餐或发布套餐")
	}

	// 校验时长单位
	if req.DurationType != v1.DurationTypeHour && req.DurationType != v1.DurationTypeDay && req.DurationType != v1.DurationTypeMonth {
		return nil, gerror.New("时长单位错误，只能是小时、天或月")
	}

	// 查询套餐
	packageDao := dao.NewPackageDao()
	p, err := packageDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("更新套餐失败: " + err.Error())
	}
	if p == nil {
		return nil, gerror.New("套餐不存在")
	}

	// 构建数据
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

	// 更新数据
	_, err = packageDao.Update(ctx, data, req.Id)
	if err != nil {
		return nil, gerror.New("更新套餐失败: " + err.Error())
	}

	return res, nil
}

// Delete 删除套餐
func (s *sPackage) Delete(ctx context.Context, req *v1.PackageDeleteReq) (res *v1.PackageDeleteRes, err error) {
	res = &v1.PackageDeleteRes{}

	// 查询套餐
	packageDao := dao.NewPackageDao()
	p, err := packageDao.FindOne(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除套餐失败: " + err.Error())
	}
	if p == nil {
		return nil, gerror.New("套餐不存在")
	}

	// 删除数据
	_, err = packageDao.Delete(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("删除套餐失败: " + err.Error())
	}

	return res, nil
}

// GetGlobalStatus 获取套餐总开关状态
func (s *sPackage) GetGlobalStatus(ctx context.Context, req *v1.PackageGlobalStatusReq) (res *v1.PackageGlobalStatusRes, err error) {
	res = &v1.PackageGlobalStatusRes{}

	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看套餐总开关状态")
	}

	// 获取置顶套餐总开关状态
	systemConfigDao := &dao.SystemConfigDao{}
	topEnabled, err := systemConfigDao.GetTopPackageEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取置顶套餐总开关状态失败: " + err.Error())
	}
	res.TopEnabled = topEnabled

	// 获取发布套餐总开关状态
	publishEnabled, err := systemConfigDao.GetPublishPackageEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取发布套餐总开关状态失败: " + err.Error())
	}
	res.PublishEnabled = publishEnabled

	return res, nil
}

// UpdateTopPackageGlobalStatus 更新置顶套餐总开关状态
func (s *sPackage) UpdateTopPackageGlobalStatus(ctx context.Context, req *v1.TopPackageGlobalStatusUpdateReq) (res *v1.TopPackageGlobalStatusUpdateRes, err error) {
	res = &v1.TopPackageGlobalStatusUpdateRes{}

	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新置顶套餐总开关状态")
	}

	// 更新置顶套餐总开关状态
	systemConfigDao := &dao.SystemConfigDao{}
	if err := systemConfigDao.UpdateTopPackageEnabled(ctx, req.IsEnabled); err != nil {
		return nil, gerror.New("更新置顶套餐总开关状态失败: " + err.Error())
	}

	return res, nil
}

// UpdatePublishPackageGlobalStatus 更新发布套餐总开关状态
func (s *sPackage) UpdatePublishPackageGlobalStatus(ctx context.Context, req *v1.PublishPackageGlobalStatusUpdateReq) (res *v1.PublishPackageGlobalStatusUpdateRes, err error) {
	res = &v1.PublishPackageGlobalStatusUpdateRes{}

	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新发布套餐总开关状态")
	}

	// 更新发布套餐总开关状态
	systemConfigDao := &dao.SystemConfigDao{}
	if err := systemConfigDao.UpdatePublishPackageEnabled(ctx, req.IsEnabled); err != nil {
		return nil, gerror.New("更新发布套餐总开关状态失败: " + err.Error())
	}

	return res, nil
}

// WxList 客户端获取套餐列表
func (s *sPackage) WxList(ctx context.Context, req *v1.WxPackageListReq) (res *v1.WxPackageListRes, err error) {
	res = &v1.WxPackageListRes{
		List: make([]*v1.Package, 0),
	}

	systemConfigDao := &dao.SystemConfigDao{}
	packageDao := dao.NewPackageDao()
	var packages []*do.PackageDO

	// 根据类型筛选套餐
	if req.Type == v1.PackageTypeTop {
		// 获取置顶套餐总开关状态
		topEnabled, err := systemConfigDao.GetTopPackageEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取置顶套餐总开关状态失败: " + err.Error())
		}
		res.IsGlobalEnabled = topEnabled

		// 如果总开关已关闭，则直接返回空列表
		if !topEnabled {
			return res, nil
		}

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
	} else if req.Type == v1.PackageTypePublish {
		// 获取发布套餐总开关状态
		publishEnabled, err := systemConfigDao.GetPublishPackageEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取发布套餐总开关状态失败: " + err.Error())
		}
		res.IsGlobalEnabled = publishEnabled

		// 如果总开关已关闭，则直接返回空列表
		if !publishEnabled {
			return res, nil
		}

		packages, err = packageDao.FindByType(ctx, string(req.Type), req.Sort, req.Order)
	} else {
		// 获取所有套餐
		packages, err = packageDao.FindAll(ctx, req.Sort, req.Order)
	}

	if err != nil {
		return nil, gerror.New("获取套餐列表失败: " + err.Error())
	}

	// 转换数据
	for _, p := range packages {
		pkg := &v1.Package{
			Id:           gconv.Int(p.Id),
			Title:        gconv.String(p.Title),
			Description:  gconv.String(p.Description),
			Price:        gconv.Float64(p.Price),
			Type:         v1.PackageType(gconv.String(p.Type)),
			Duration:     gconv.Int(p.Duration),
			DurationType: v1.DurationType(gconv.String(p.DurationType)),
			SortOrder:    gconv.Int(p.SortOrder),
			CreatedAt:    p.CreatedAt.String(),
			UpdatedAt:    p.UpdatedAt.String(),
		}
		res.List = append(res.List, pkg)
	}

	return res, nil
}
