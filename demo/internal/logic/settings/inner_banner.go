package settings

import (
	"context"
	v1 "demo/api/content/v1"
	"demo/internal/model/entity"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
)

// GetInnerBannerList 获取内页轮播图列表
func (s *sSettings) GetInnerBannerList(ctx context.Context, req *v1.InnerBannerListReq) (res *v1.InnerBannerListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看内页轮播图列表")
	}

	// 初始化响应对象
	res = &v1.InnerBannerListRes{
		List: make([]v1.InnerBannerItem, 0),
	}

	// 根据轮播图类型获取对应的总开关状态
	var typeEnabled bool
	if req.BannerType == "home" {
		typeEnabled, err = s.systemConfigDao.GetHomeInnerBannerEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取首页内页轮播图总开关状态失败: " + err.Error())
		}
	} else if req.BannerType == "idle" {
		typeEnabled, err = s.systemConfigDao.GetIdleInnerBannerEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取闲置页内页轮播图总开关状态失败: " + err.Error())
		}
	}

	// 返回开关状态
	res.IsGlobalEnabled = typeEnabled

	// 获取内页轮播图列表
	banners, err := s.innerBannerDao.GetList(ctx, req.BannerType, false)
	if err != nil {
		return nil, gerror.New("获取内页轮播图列表失败: " + err.Error())
	}

	// 转换为API响应格式
	for _, b := range banners {
		res.List = append(res.List, v1.InnerBannerItem{
			Id:         b.Id,
			Image:      b.Image,
			LinkType:   b.LinkType,
			LinkUrl:    b.LinkUrl,
			IsEnabled:  b.IsEnabled == 1,
			Order:      b.Order,
			BannerType: b.BannerType,
		})
	}

	return res, nil
}

// CreateInnerBanner 创建内页轮播图
func (s *sSettings) CreateInnerBanner(ctx context.Context, req *v1.InnerBannerCreateReq) (res *v1.InnerBannerCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建内页轮播图")
	}

	// 准备保存的数据
	banner := &entity.InnerPageBanner{
		Image:      req.Image,
		LinkType:   req.LinkType,
		LinkUrl:    req.LinkUrl,
		IsEnabled:  0,
		Order:      req.Order,
		BannerType: req.BannerType,
	}

	if req.IsEnabled {
		banner.IsEnabled = 1
	}

	// 创建内页轮播图
	id, err := s.innerBannerDao.Create(ctx, banner)
	if err != nil {
		return nil, gerror.New("创建内页轮播图失败: " + err.Error())
	}

	return &v1.InnerBannerCreateRes{
		Id: int(id),
	}, nil
}

// UpdateInnerBanner 更新内页轮播图
func (s *sSettings) UpdateInnerBanner(ctx context.Context, req *v1.InnerBannerUpdateReq) (res *v1.InnerBannerUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新内页轮播图")
	}

	// 检查内页轮播图是否存在
	_, err = s.innerBannerDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("内页轮播图不存在或已被删除")
	}

	// 准备更新的数据
	banner := &entity.InnerPageBanner{
		Id:         req.Id,
		Image:      req.Image,
		LinkType:   req.LinkType,
		LinkUrl:    req.LinkUrl,
		IsEnabled:  0,
		Order:      req.Order,
		BannerType: req.BannerType,
	}

	if req.IsEnabled {
		banner.IsEnabled = 1
	}

	// 更新内页轮播图
	if err := s.innerBannerDao.Update(ctx, banner); err != nil {
		return nil, gerror.New("更新内页轮播图失败: " + err.Error())
	}

	return &v1.InnerBannerUpdateRes{}, nil
}

// DeleteInnerBanner 删除内页轮播图
func (s *sSettings) DeleteInnerBanner(ctx context.Context, req *v1.InnerBannerDeleteReq) (res *v1.InnerBannerDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除内页轮播图")
	}

	// 检查内页轮播图是否存在
	_, err = s.innerBannerDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("内页轮播图不存在或已被删除")
	}

	// 删除内页轮播图
	if err := s.innerBannerDao.Delete(ctx, req.Id); err != nil {
		return nil, gerror.New("删除内页轮播图失败: " + err.Error())
	}

	return &v1.InnerBannerDeleteRes{}, nil
}

// UpdateInnerBannerStatus 更新内页轮播图状态
func (s *sSettings) UpdateInnerBannerStatus(ctx context.Context, req *v1.InnerBannerStatusUpdateReq) (res *v1.InnerBannerStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新内页轮播图状态")
	}

	// 检查内页轮播图是否存在
	_, err = s.innerBannerDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("内页轮播图不存在或已被删除")
	}

	// 更新状态
	isEnabled := 0
	if req.IsEnabled {
		isEnabled = 1
	}

	if err := s.innerBannerDao.UpdateStatus(ctx, req.Id, isEnabled); err != nil {
		return nil, gerror.New("更新内页轮播图状态失败: " + err.Error())
	}

	return &v1.InnerBannerStatusUpdateRes{}, nil
}

// UpdateHomeInnerBannerGlobalStatus 更新首页内页轮播图总开关状态
func (s *sSettings) UpdateHomeInnerBannerGlobalStatus(ctx context.Context, req *v1.HomeInnerBannerGlobalStatusUpdateReq) (res *v1.HomeInnerBannerGlobalStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新首页内页轮播图总开关状态")
	}

	// 更新总开关状态
	if err := s.systemConfigDao.UpdateHomeInnerBannerEnabled(ctx, req.IsEnabled); err != nil {
		return nil, gerror.New("更新首页内页轮播图总开关状态失败: " + err.Error())
	}

	return &v1.HomeInnerBannerGlobalStatusUpdateRes{}, nil
}

// UpdateIdleInnerBannerGlobalStatus 更新闲置页内页轮播图总开关状态
func (s *sSettings) UpdateIdleInnerBannerGlobalStatus(ctx context.Context, req *v1.IdleInnerBannerGlobalStatusUpdateReq) (res *v1.IdleInnerBannerGlobalStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新闲置页内页轮播图总开关状态")
	}

	// 更新总开关状态
	if err := s.systemConfigDao.UpdateIdleInnerBannerEnabled(ctx, req.IsEnabled); err != nil {
		return nil, gerror.New("更新闲置页内页轮播图总开关状态失败: " + err.Error())
	}

	return &v1.IdleInnerBannerGlobalStatusUpdateRes{}, nil
}

// WxGetInnerBannerList 微信客户端获取内页轮播图列表
func (s *sSettings) WxGetInnerBannerList(ctx context.Context, req *v1.WxInnerBannerListReq) (res *v1.WxInnerBannerListRes, err error) {
	// 初始化响应对象
	res = &v1.WxInnerBannerListRes{
		List: make([]v1.InnerBannerItem, 0),
	}

	// 根据轮播图类型获取对应的总开关状态
	var typeEnabled bool
	if req.BannerType == "home" {
		typeEnabled, err = s.systemConfigDao.GetHomeInnerBannerEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取首页内页轮播图总开关状态失败: " + err.Error())
		}
	} else if req.BannerType == "idle" {
		typeEnabled, err = s.systemConfigDao.GetIdleInnerBannerEnabled(ctx)
		if err != nil {
			return nil, gerror.New("获取闲置页内页轮播图总开关状态失败: " + err.Error())
		}
	}

	res.IsGlobalEnabled = typeEnabled

	// 如果开关已关闭，则直接返回空列表
	if !typeEnabled {
		return res, nil
	}

	// 获取内页轮播图列表，只获取已启用的项目
	banners, err := s.innerBannerDao.GetList(ctx, req.BannerType, true)
	if err != nil {
		return nil, gerror.New("获取内页轮播图列表失败: " + err.Error())
	}

	// 转换为API响应格式
	for _, b := range banners {
		res.List = append(res.List, v1.InnerBannerItem{
			Id:         b.Id,
			Image:      b.Image,
			LinkType:   b.LinkType,
			LinkUrl:    b.LinkUrl,
			IsEnabled:  b.IsEnabled == 1,
			Order:      b.Order,
			BannerType: b.BannerType,
		})
	}

	return res, nil
}
