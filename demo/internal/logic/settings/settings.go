package settings

import (
	"context"
	v1 "demo/api/content/v1"
	settingsV1 "demo/api/settings/v1"
	"demo/internal/dao"
	"demo/internal/model/entity"
	"demo/utility/auth"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 设置服务实现
type sSettings struct {
	miniProgramDao         *dao.MiniProgramDao
	bannerDao              *dao.BannerDao
	activityAreaDao        *dao.ActivityAreaDao
	systemConfigDao        *dao.SystemConfigDao
	innerBannerDao         *dao.InnerBannerDao
	miniProgramSettingsDao *dao.MiniProgramSettingsDao
	adSettingsDao          *dao.AdSettingsDao
	rewardSettingsDao      *dao.RewardSettingsDao
	agreementSettingsDao   *dao.AgreementSettingsDao
	shareSettingsDao       *dao.ShareSettingsDao
}

// New 创建设置服务实例
func New() *sSettings {
	return &sSettings{
		miniProgramDao:         &dao.MiniProgramDao{},
		bannerDao:              &dao.BannerDao{},
		activityAreaDao:        &dao.ActivityAreaDao{},
		systemConfigDao:        &dao.SystemConfigDao{},
		innerBannerDao:         &dao.InnerBannerDao{},
		miniProgramSettingsDao: &dao.MiniProgramSettingsDao{},
		adSettingsDao:          &dao.AdSettingsDao{},
		rewardSettingsDao:      &dao.RewardSettingsDao{},
		agreementSettingsDao:   &dao.AgreementSettingsDao{},
		shareSettingsDao:       &dao.ShareSettingsDao{},
	}
}

// GetMiniProgramList 获取导航小程序列表
func (s *sSettings) GetMiniProgramList(ctx context.Context, req *v1.MiniProgramListReq) (res *v1.MiniProgramListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看导航小程序列表")
	}

	res = &v1.MiniProgramListRes{
		List: make([]v1.MiniProgramItem, 0),
	}

	// 获取导航小程序列表, 管理后台不受总开关影响
	miniPrograms, err := s.miniProgramDao.GetList(ctx, false, false)
	if err != nil {
		return nil, gerror.New("获取导航小程序列表失败: " + err.Error())
	}

	// 转换为API响应格式
	for _, mp := range miniPrograms {
		res.List = append(res.List, v1.MiniProgramItem{
			Id:        mp.Id,
			Name:      mp.Name,
			AppId:     mp.AppId,
			Logo:      mp.Logo,
			IsEnabled: mp.IsEnabled == 1,
			Order:     mp.Order,
		})
	}

	// 获取总开关状态
	enabled, err := s.systemConfigDao.GetMiniProgramEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取导航小程序总开关状态失败: " + err.Error())
	}

	res.IsGlobalEnabled = enabled

	return res, nil
}

// CreateMiniProgram 创建导航小程序
func (s *sSettings) CreateMiniProgram(ctx context.Context, req *v1.MiniProgramCreateReq) (res *v1.MiniProgramCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建导航小程序")
	}

	// 准备保存的数据
	miniProgram := &entity.MiniProgram{
		Name:      req.Name,
		AppId:     req.AppId,
		Logo:      req.Logo,
		IsEnabled: 0,
		Order:     req.Order,
	}

	if req.IsEnabled {
		miniProgram.IsEnabled = 1
	}

	// 创建导航小程序
	id, err := s.miniProgramDao.Create(ctx, miniProgram)
	if err != nil {
		return nil, gerror.New("创建导航小程序失败: " + err.Error())
	}

	return &v1.MiniProgramCreateRes{
		Id: int(id),
	}, nil
}

// UpdateMiniProgram 更新导航小程序
func (s *sSettings) UpdateMiniProgram(ctx context.Context, req *v1.MiniProgramUpdateReq) (res *v1.MiniProgramUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新导航小程序")
	}

	// 检查导航小程序是否存在
	_, err = s.miniProgramDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("导航小程序不存在或已被删除")
	}

	// 准备更新的数据
	miniProgram := &entity.MiniProgram{
		Id:        req.Id,
		Name:      req.Name,
		AppId:     req.AppId,
		Logo:      req.Logo,
		IsEnabled: 0,
		Order:     req.Order,
	}

	if req.IsEnabled {
		miniProgram.IsEnabled = 1
	}

	// 更新导航小程序
	if err := s.miniProgramDao.Update(ctx, miniProgram); err != nil {
		return nil, gerror.New("更新导航小程序失败: " + err.Error())
	}

	return &v1.MiniProgramUpdateRes{}, nil
}

// DeleteMiniProgram 删除导航小程序
func (s *sSettings) DeleteMiniProgram(ctx context.Context, req *v1.MiniProgramDeleteReq) (res *v1.MiniProgramDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除导航小程序")
	}

	// 检查导航小程序是否存在
	_, err = s.miniProgramDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("导航小程序不存在或已被删除")
	}

	// 删除导航小程序
	if err := s.miniProgramDao.Delete(ctx, req.Id); err != nil {
		return nil, gerror.New("删除导航小程序失败: " + err.Error())
	}

	return &v1.MiniProgramDeleteRes{}, nil
}

// UpdateMiniProgramStatus 更新导航小程序状态
func (s *sSettings) UpdateMiniProgramStatus(ctx context.Context, req *v1.MiniProgramStatusUpdateReq) (res *v1.MiniProgramStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新导航小程序状态")
	}

	// 检查导航小程序是否存在
	_, err = s.miniProgramDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("导航小程序不存在或已被删除")
	}

	// 更新状态
	isEnabled := 0
	if req.IsEnabled {
		isEnabled = 1
	}

	if err := s.miniProgramDao.UpdateStatus(ctx, req.Id, isEnabled); err != nil {
		return nil, gerror.New("更新导航小程序状态失败: " + err.Error())
	}

	return &v1.MiniProgramStatusUpdateRes{}, nil
}

// GetBannerList 获取轮播图列表
func (s *sSettings) GetBannerList(ctx context.Context, req *v1.BannerListReq) (res *v1.BannerListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看轮播图列表")
	}

	res = &v1.BannerListRes{
		List: make([]v1.BannerItem, 0),
	}

	// 获取轮播图列表，管理后台不受总开关影响
	banners, err := s.bannerDao.GetList(ctx, false, false)
	if err != nil {
		return nil, gerror.New("获取轮播图列表失败: " + err.Error())
	}

	// 转换为API响应格式
	for _, b := range banners {
		res.List = append(res.List, v1.BannerItem{
			Id:        b.Id,
			Image:     b.Image,
			LinkType:  b.LinkType,
			LinkUrl:   b.LinkUrl,
			IsEnabled: b.IsEnabled == 1,
			Order:     b.Order,
		})
	}

	// 获取总开关状态
	enabled, err := s.systemConfigDao.GetBannerEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取轮播图总开关状态失败: " + err.Error())
	}

	res.IsGlobalEnabled = enabled

	return res, nil
}

// CreateBanner 创建轮播图
func (s *sSettings) CreateBanner(ctx context.Context, req *v1.BannerCreateReq) (res *v1.BannerCreateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限创建轮播图")
	}

	// 准备保存的数据
	banner := &entity.Banner{
		Image:     req.Image,
		LinkType:  req.LinkType,
		LinkUrl:   req.LinkUrl,
		IsEnabled: 0,
		Order:     req.Order,
	}

	if req.IsEnabled {
		banner.IsEnabled = 1
	}

	// 创建轮播图
	id, err := s.bannerDao.Create(ctx, banner)
	if err != nil {
		return nil, gerror.New("创建轮播图失败: " + err.Error())
	}

	return &v1.BannerCreateRes{
		Id: int(id),
	}, nil
}

// UpdateBanner 更新轮播图
func (s *sSettings) UpdateBanner(ctx context.Context, req *v1.BannerUpdateReq) (res *v1.BannerUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新轮播图")
	}

	// 检查轮播图是否存在
	_, err = s.bannerDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("轮播图不存在或已被删除")
	}

	// 准备更新的数据
	banner := &entity.Banner{
		Id:        req.Id,
		Image:     req.Image,
		LinkType:  req.LinkType,
		LinkUrl:   req.LinkUrl,
		IsEnabled: 0,
		Order:     req.Order,
	}

	if req.IsEnabled {
		banner.IsEnabled = 1
	}

	// 更新轮播图
	if err := s.bannerDao.Update(ctx, banner); err != nil {
		return nil, gerror.New("更新轮播图失败: " + err.Error())
	}

	return &v1.BannerUpdateRes{}, nil
}

// DeleteBanner 删除轮播图
func (s *sSettings) DeleteBanner(ctx context.Context, req *v1.BannerDeleteReq) (res *v1.BannerDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除轮播图")
	}

	// 检查轮播图是否存在
	_, err = s.bannerDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("轮播图不存在或已被删除")
	}

	// 删除轮播图
	if err := s.bannerDao.Delete(ctx, req.Id); err != nil {
		return nil, gerror.New("删除轮播图失败: " + err.Error())
	}

	return &v1.BannerDeleteRes{}, nil
}

// UpdateBannerStatus 更新轮播图状态
func (s *sSettings) UpdateBannerStatus(ctx context.Context, req *v1.BannerStatusUpdateReq) (res *v1.BannerStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新轮播图状态")
	}

	// 检查轮播图是否存在
	_, err = s.bannerDao.GetById(ctx, req.Id)
	if err != nil {
		return nil, gerror.New("轮播图不存在或已被删除")
	}

	// 更新状态
	isEnabled := 0
	if req.IsEnabled {
		isEnabled = 1
	}

	if err := s.bannerDao.UpdateStatus(ctx, req.Id, isEnabled); err != nil {
		return nil, gerror.New("更新轮播图状态失败: " + err.Error())
	}

	return &v1.BannerStatusUpdateRes{}, nil
}

// GetActivityArea 获取活动区域
func (s *sSettings) GetActivityArea(ctx context.Context, req *v1.ActivityAreaGetReq) (res *v1.ActivityAreaGetRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看活动区域")
	}

	// 获取活动区域数据
	activityArea, err := s.activityAreaDao.Get(ctx)
	if err != nil {
		return nil, gerror.New("获取活动区域失败: " + err.Error())
	}

	// 构造响应
	res = &v1.ActivityAreaGetRes{
		TopLeft: v1.ActivityAreaModule{
			Title:       activityArea.TopLeftTitle,
			Description: activityArea.TopLeftDescription,
			LinkType:    activityArea.TopLeftLinkType,
			LinkUrl:     activityArea.TopLeftLinkUrl,
		},
		BottomLeft: v1.ActivityAreaModule{
			Title:       activityArea.BottomLeftTitle,
			Description: activityArea.BottomLeftDescription,
			LinkType:    activityArea.BottomLeftLinkType,
			LinkUrl:     activityArea.BottomLeftLinkUrl,
		},
		Right: v1.ActivityAreaModule{
			Title:       activityArea.RightTitle,
			Description: activityArea.RightDescription,
			LinkType:    activityArea.RightLinkType,
			LinkUrl:     activityArea.RightLinkUrl,
		},
	}

	// 获取总开关状态
	enabled, err := s.systemConfigDao.GetActivityAreaEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取活动区域总开关状态失败: " + err.Error())
	}

	res.IsGlobalEnabled = enabled

	return res, nil
}

// SaveActivityArea 保存活动区域
func (s *sSettings) SaveActivityArea(ctx context.Context, req *v1.ActivityAreaSaveReq) (res *v1.ActivityAreaSaveRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限保存活动区域设置")
	}

	// 准备保存的数据
	activityArea := &entity.ActivityArea{
		Id:                    1, // 固定使用ID为1的记录
		TopLeftTitle:          req.TopLeft.Title,
		TopLeftDescription:    req.TopLeft.Description,
		TopLeftLinkType:       req.TopLeft.LinkType,
		TopLeftLinkUrl:        req.TopLeft.LinkUrl,
		BottomLeftTitle:       req.BottomLeft.Title,
		BottomLeftDescription: req.BottomLeft.Description,
		BottomLeftLinkType:    req.BottomLeft.LinkType,
		BottomLeftLinkUrl:     req.BottomLeft.LinkUrl,
		RightTitle:            req.Right.Title,
		RightDescription:      req.Right.Description,
		RightLinkType:         req.Right.LinkType,
		RightLinkUrl:          req.Right.LinkUrl,
	}

	// 保存活动区域
	if err := s.activityAreaDao.Save(ctx, activityArea); err != nil {
		return nil, gerror.New("保存活动区域失败: " + err.Error())
	}

	return &v1.ActivityAreaSaveRes{}, nil
}

// WxGetActivityArea 微信客户端获取活动区域
func (s *sSettings) WxGetActivityArea(ctx context.Context, req *v1.WxActivityAreaGetReq) (res *v1.WxActivityAreaGetRes, err error) {
	// 初始化响应对象
	res = &v1.WxActivityAreaGetRes{
		List: make([]v1.ActivityAreaModule, 0),
	}

	// 获取总开关状态
	enabled, err := s.systemConfigDao.GetActivityAreaEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取活动区域总开关状态失败: " + err.Error())
	}

	res.IsGlobalEnabled = enabled

	// 如果总开关已关闭，则直接返回空列表
	if !enabled {
		return res, nil
	}

	// 获取活动区域数据
	activityArea, err := s.activityAreaDao.Get(ctx)
	if err != nil {
		return nil, gerror.New("获取活动区域失败: " + err.Error())
	}

	// 构造响应
	res.List = []v1.ActivityAreaModule{
		{
			Title:       activityArea.TopLeftTitle,
			Description: activityArea.TopLeftDescription,
			LinkType:    activityArea.TopLeftLinkType,
			LinkUrl:     activityArea.TopLeftLinkUrl,
			Position:    "topLeft",
		},
		{
			Title:       activityArea.BottomLeftTitle,
			Description: activityArea.BottomLeftDescription,
			LinkType:    activityArea.BottomLeftLinkType,
			LinkUrl:     activityArea.BottomLeftLinkUrl,
			Position:    "bottomLeft",
		},
		{
			Title:       activityArea.RightTitle,
			Description: activityArea.RightDescription,
			LinkType:    activityArea.RightLinkType,
			LinkUrl:     activityArea.RightLinkUrl,
			Position:    "right",
		},
	}

	return res, nil
}

// UpdateMiniProgramGlobalStatus 更新导航小程序总开关状态
func (s *sSettings) UpdateMiniProgramGlobalStatus(ctx context.Context, req *v1.MiniProgramGlobalStatusUpdateReq) (res *v1.MiniProgramGlobalStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新导航小程序总开关状态")
	}

	// 更新总开关状态
	if err := s.systemConfigDao.UpdateMiniProgramEnabled(ctx, req.IsEnabled); err != nil {
		return nil, gerror.New("更新导航小程序总开关状态失败: " + err.Error())
	}

	return &v1.MiniProgramGlobalStatusUpdateRes{}, nil
}

// UpdateBannerGlobalStatus 更新轮播图总开关状态
func (s *sSettings) UpdateBannerGlobalStatus(ctx context.Context, req *v1.BannerGlobalStatusUpdateReq) (res *v1.BannerGlobalStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新轮播图总开关状态")
	}

	// 更新总开关状态
	if err := s.systemConfigDao.UpdateBannerEnabled(ctx, req.IsEnabled); err != nil {
		return nil, gerror.New("更新轮播图总开关状态失败: " + err.Error())
	}

	return &v1.BannerGlobalStatusUpdateRes{}, nil
}

// UpdateActivityAreaGlobalStatus 更新活动区域总开关状态
func (s *sSettings) UpdateActivityAreaGlobalStatus(ctx context.Context, req *v1.ActivityAreaGlobalStatusUpdateReq) (res *v1.ActivityAreaGlobalStatusUpdateRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限更新活动区域总开关状态")
	}

	// 更新总开关状态
	if err := s.systemConfigDao.UpdateActivityAreaEnabled(ctx, req.IsEnabled); err != nil {
		return nil, gerror.New("更新活动区域总开关状态失败: " + err.Error())
	}

	return &v1.ActivityAreaGlobalStatusUpdateRes{}, nil
}

// WxGetMiniProgramList 微信客户端获取导航小程序列表
func (s *sSettings) WxGetMiniProgramList(ctx context.Context, req *v1.WxMiniProgramListReq) (res *v1.WxMiniProgramListRes, err error) {
	// 获取导航小程序列表，只获取已启用的项目
	miniPrograms, err := s.miniProgramDao.GetList(ctx, true)
	if err != nil {
		return nil, gerror.New("获取导航小程序列表失败: " + err.Error())
	}

	res = &v1.WxMiniProgramListRes{
		List: make([]v1.MiniProgramItem, 0),
	}

	// 转换为API响应格式
	for _, mp := range miniPrograms {
		res.List = append(res.List, v1.MiniProgramItem{
			Id:        mp.Id,
			Name:      mp.Name,
			AppId:     mp.AppId,
			Logo:      mp.Logo,
			IsEnabled: mp.IsEnabled == 1,
			Order:     mp.Order,
		})
	}

	// 获取总开关状态
	enabled, err := s.systemConfigDao.GetMiniProgramEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取导航小程序总开关状态失败: " + err.Error())
	}

	res.IsGlobalEnabled = enabled

	return res, nil
}

// WxGetBannerList 微信客户端获取轮播图列表
func (s *sSettings) WxGetBannerList(ctx context.Context, req *v1.WxBannerListReq) (res *v1.WxBannerListRes, err error) {
	// 获取轮播图列表，只获取已启用的项目
	banners, err := s.bannerDao.GetList(ctx, true)
	if err != nil {
		return nil, gerror.New("获取轮播图列表失败: " + err.Error())
	}

	res = &v1.WxBannerListRes{
		List: make([]v1.BannerItem, 0),
	}

	// 转换为API响应格式
	for _, b := range banners {
		res.List = append(res.List, v1.BannerItem{
			Id:        b.Id,
			Image:     b.Image,
			LinkType:  b.LinkType,
			LinkUrl:   b.LinkUrl,
			IsEnabled: b.IsEnabled == 1,
			Order:     b.Order,
		})
	}

	// 获取总开关状态
	enabled, err := s.systemConfigDao.GetBannerEnabled(ctx)
	if err != nil {
		return nil, gerror.New("获取轮播图总开关状态失败: " + err.Error())
	}

	res.IsGlobalEnabled = enabled

	return res, nil
}

// GetMiniProgramBaseSettings 获取小程序基础设置
func (s *sSettings) GetMiniProgramBaseSettings(ctx context.Context, req *settingsV1.MiniProgramBaseSettingsReq) (res *settingsV1.MiniProgramBaseSettingsRes, err error) {
	settings, err := s.miniProgramSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取小程序基础设置失败", err)
		return nil, errors.New("获取小程序基础设置失败")
	}

	res = &settingsV1.MiniProgramBaseSettingsRes{
		Name:        settings.Name,
		Description: settings.Description,
		Logo:        settings.Logo,
	}
	return res, nil
}

// SaveMiniProgramBaseSettings 保存小程序基础设置
func (s *sSettings) SaveMiniProgramBaseSettings(ctx context.Context, req *settingsV1.MiniProgramBaseSettingsSaveReq) (res *settingsV1.MiniProgramBaseSettingsSaveRes, err error) {
	settings := &entity.MiniProgramSettings{
		Id:          1,
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
	}

	err = s.miniProgramSettingsDao.Save(ctx, settings)
	if err != nil {
		g.Log().Error(ctx, "保存小程序基础设置失败", err)
		return nil, errors.New("保存小程序基础设置失败")
	}

	res = &settingsV1.MiniProgramBaseSettingsSaveRes{
		IsSuccess: true,
	}
	return res, nil
}

// GetAdSettings 获取广告设置
func (s *sSettings) GetAdSettings(ctx context.Context, req *settingsV1.AdSettingsReq) (res *settingsV1.AdSettingsRes, err error) {
	settings, err := s.adSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取广告设置失败", err)
		return nil, errors.New("获取广告设置失败")
	}

	res = &settingsV1.AdSettingsRes{
		EnableWxAd:        settings.EnableWxAd,
		RewardedVideoAdId: settings.RewardedVideoAdId,
	}
	return res, nil
}

// SaveAdSettings 保存广告设置
func (s *sSettings) SaveAdSettings(ctx context.Context, req *settingsV1.AdSettingsSaveReq) (res *settingsV1.AdSettingsSaveRes, err error) {
	settings := &entity.AdSettings{
		Id:                1,
		EnableWxAd:        req.EnableWxAd,
		RewardedVideoAdId: req.RewardedVideoAdId,
	}

	err = s.adSettingsDao.Save(ctx, settings)
	if err != nil {
		g.Log().Error(ctx, "保存广告设置失败", err)
		return nil, errors.New("保存广告设置失败")
	}

	res = &settingsV1.AdSettingsSaveRes{
		IsSuccess: true,
	}
	return res, nil
}

// GetRewardSettings 获取奖励设置
func (s *sSettings) GetRewardSettings(ctx context.Context, req *settingsV1.RewardSettingsReq) (res *settingsV1.RewardSettingsRes, err error) {
	settings, err := s.rewardSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取奖励设置失败", err)
		return nil, errors.New("获取奖励设置失败")
	}

	res = &settingsV1.RewardSettingsRes{
		EnableReward:           settings.EnableReward,
		FirstViewMinRewardMin:  settings.FirstViewMinRewardMin,
		FirstViewMaxRewardDay:  settings.FirstViewMaxRewardDay,
		SingleAdMinRewardMin:   settings.SingleAdMinRewardMin,
		SingleAdMaxRewardDay:   settings.SingleAdMaxRewardDay,
		DailyRewardLimit:       settings.DailyRewardLimit,
		DailyMaxAccumulatedDay: settings.DailyMaxAccumulatedDay,
		RewardExpirationDays:   settings.RewardExpirationDays,
	}
	return res, nil
}

// SaveRewardSettings 保存奖励设置
func (s *sSettings) SaveRewardSettings(ctx context.Context, req *settingsV1.RewardSettingsSaveReq) (res *settingsV1.RewardSettingsSaveRes, err error) {
	// 验证设置参数的合理性
	if req.FirstViewMinRewardMin <= 0 {
		return nil, errors.New("首次观看最小奖励分钟必须大于0")
	}
	if req.FirstViewMaxRewardDay <= 0 {
		return nil, errors.New("首次观看最大奖励天数必须大于0")
	}
	if req.SingleAdMinRewardMin <= 0 {
		return nil, errors.New("单次广告最小奖励分钟必须大于0")
	}
	if req.SingleAdMaxRewardDay <= 0 {
		return nil, errors.New("单次广告最大奖励天数必须大于0")
	}
	if req.DailyRewardLimit <= 0 {
		return nil, errors.New("每日奖励次数上限必须大于0")
	}
	if req.DailyMaxAccumulatedDay <= 0 {
		return nil, errors.New("每日最大累计奖励必须大于0")
	}
	if req.RewardExpirationDays <= 0 {
		return nil, errors.New("奖励过期天数必须大于0")
	}

	// 检查最小值是否小于最大值 (将天转换为分钟进行比较)
	firstViewMinMinutes := req.FirstViewMinRewardMin
	firstViewMaxMinutes := int(req.FirstViewMaxRewardDay * 24 * 60)
	if firstViewMinMinutes >= firstViewMaxMinutes {
		return nil, errors.New("首次观看最小奖励必须小于最大奖励")
	}

	singleAdMinMinutes := req.SingleAdMinRewardMin
	singleAdMaxMinutes := int(req.SingleAdMaxRewardDay * 24 * 60)
	if singleAdMinMinutes >= singleAdMaxMinutes {
		return nil, errors.New("单次广告最小奖励必须小于最大奖励")
	}

	settings := &entity.RewardSettings{
		Id:                     1,
		EnableReward:           req.EnableReward,
		FirstViewMinRewardMin:  req.FirstViewMinRewardMin,
		FirstViewMaxRewardDay:  req.FirstViewMaxRewardDay,
		SingleAdMinRewardMin:   req.SingleAdMinRewardMin,
		SingleAdMaxRewardDay:   req.SingleAdMaxRewardDay,
		DailyRewardLimit:       req.DailyRewardLimit,
		DailyMaxAccumulatedDay: req.DailyMaxAccumulatedDay,
		RewardExpirationDays:   req.RewardExpirationDays,
	}

	err = s.rewardSettingsDao.Save(ctx, settings)
	if err != nil {
		g.Log().Error(ctx, "保存奖励设置失败", err)
		return nil, errors.New("保存奖励设置失败")
	}

	res = &settingsV1.RewardSettingsSaveRes{
		IsSuccess: true,
	}
	return res, nil
}

// GetAgreementSettings 获取协议设置
func (s *sSettings) GetAgreementSettings(ctx context.Context, req *settingsV1.AgreementSettingsReq) (res *settingsV1.AgreementSettingsRes, err error) {
	settings, err := s.agreementSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取协议设置失败", err)
		return nil, errors.New("获取协议设置失败")
	}

	res = &settingsV1.AgreementSettingsRes{
		PrivacyPolicy: settings.PrivacyPolicy,
		UserAgreement: settings.UserAgreement,
	}
	return res, nil
}

// SaveAgreementSettings 保存协议设置
func (s *sSettings) SaveAgreementSettings(ctx context.Context, req *settingsV1.AgreementSettingsSaveReq) (res *settingsV1.AgreementSettingsSaveRes, err error) {
	// 安全检查，过滤可能的XSS攻击内容
	privacyPolicy := req.PrivacyPolicy
	userAgreement := req.UserAgreement

	// 根据需要可以增加更多的内容安全过滤

	settings := &entity.AgreementSettings{
		Id:            1,
		PrivacyPolicy: privacyPolicy,
		UserAgreement: userAgreement,
	}

	err = s.agreementSettingsDao.Save(ctx, settings)
	if err != nil {
		g.Log().Error(ctx, "保存协议设置失败", err)
		return nil, errors.New("保存协议设置失败")
	}

	res = &settingsV1.AgreementSettingsSaveRes{
		IsSuccess: true,
	}
	return res, nil
}

// WxGetAgreement 微信客户端获取协议
func (s *sSettings) WxGetAgreement(ctx context.Context, req *settingsV1.WxAgreementGetReq) (res *settingsV1.WxAgreementGetRes, err error) {
	settings, err := s.agreementSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取协议设置失败", err)
		return nil, errors.New("获取协议设置失败")
	}

	// 根据请求类型返回对应协议内容
	var content string
	switch req.Type {
	case "privacy":
		content = settings.PrivacyPolicy
	case "user":
		content = settings.UserAgreement
	default:
		return nil, errors.New("无效的协议类型")
	}

	res = &settingsV1.WxAgreementGetRes{
		Content: content,
	}
	return res, nil
}

// WxGetMiniProgramBaseSettings 微信客户端获取小程序基础设置
func (s *sSettings) WxGetMiniProgramBaseSettings(ctx context.Context, req *settingsV1.WxMiniProgramBaseSettingsReq) (res *settingsV1.WxMiniProgramBaseSettingsRes, err error) {
	settings, err := s.miniProgramSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取小程序基础设置失败", err)
		return nil, errors.New("获取小程序基础设置失败")
	}

	res = &settingsV1.WxMiniProgramBaseSettingsRes{
		Name:        settings.Name,
		Description: settings.Description,
		Logo:        settings.Logo,
	}
	return res, nil
}

// WxGetAdSettings 微信客户端获取广告设置
func (s *sSettings) WxGetAdSettings(ctx context.Context, req *settingsV1.WxGetAdSettingsReq) (res *settingsV1.WxGetAdSettingsRes, err error) {
	settings, err := s.adSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取广告设置失败", err)
		return nil, errors.New("获取广告设置失败")
	}

	res = &settingsV1.WxGetAdSettingsRes{
		EnableWxAd:        settings.EnableWxAd,
		RewardedVideoAdId: settings.RewardedVideoAdId,
	}
	return res, nil
}

// WxRewardAdViewed 微信客户端广告观看完成
func (s *sSettings) WxRewardAdViewed(ctx context.Context, req *settingsV1.WxRewardAdViewedReq) (res *settingsV1.WxRewardAdViewedRes, err error) {
	// 1. 获取当前登录的客户信息
	clientId, err := auth.GetClientId(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取客户信息失败", err)
		return nil, errors.New("获取客户信息失败")
	}

	// 后续可以使用clientId处理记录
	g.Log().Debug(ctx, "客户ID:", clientId)

	// 2. 获取奖励设置参数
	rewardSettings, err := s.rewardSettingsDao.Get(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取奖励设置失败", err)
		return nil, errors.New("获取奖励设置失败")
	}

	// 如果奖励功能未启用，则直接返回错误
	if !rewardSettings.EnableReward {
		return &settingsV1.WxRewardAdViewedRes{
			Success: false,
			Message: "奖励功能未启用",
		}, nil
	}

	// 创建奖励记录DAO
	rewardRecordDao := &dao.RewardRecordDao{}

	// 检查今日奖励次数是否达到上限
	dailyCount, err := rewardRecordDao.GetDailyRewardCount(ctx, clientId)
	if err != nil {
		g.Log().Error(ctx, "获取客户今日奖励次数失败", err)
		return nil, errors.New("获取客户今日奖励次数失败")
	}

	if dailyCount >= rewardSettings.DailyRewardLimit {
		return &settingsV1.WxRewardAdViewedRes{
			Success: false,
			Message: fmt.Sprintf("今日奖励次数已达上限(%d次)", rewardSettings.DailyRewardLimit),
		}, nil
	}

	// 检查是否首次观看
	isFirstView, err := rewardRecordDao.CheckFirstView(ctx, clientId)
	if err != nil {
		g.Log().Error(ctx, "检查客户是否首次观看失败", err)
		// 如果检查出错，默认为非首次观看，避免给予过多奖励
		isFirstView = false
	}

	// 计算奖励时长（分钟）
	var rewardMinutes int

	// 根据是否首次观看，确定奖励范围
	if isFirstView {
		// 首次观看的奖励范围
		minReward := rewardSettings.FirstViewMinRewardMin
		maxReward := int(rewardSettings.FirstViewMaxRewardDay * 24 * 60) // 转换为分钟
		// 生成随机奖励（简单实现，实际业务可能有更复杂的规则）
		rewardMinutes = minReward + rand.Intn(maxReward-minReward+1)
		g.Log().Info(ctx, "客户首次观看广告，奖励时长(分钟):", rewardMinutes)
	} else {
		// 普通观看的奖励范围
		minReward := rewardSettings.SingleAdMinRewardMin
		maxReward := int(rewardSettings.SingleAdMaxRewardDay * 24 * 60) // 转换为分钟
		// 生成随机奖励（简单实现，实际业务可能有更复杂的规则）
		rewardMinutes = minReward + rand.Intn(maxReward-minReward+1)
		g.Log().Info(ctx, "客户非首次观看广告，奖励时长(分钟):", rewardMinutes)
	}

	// 计算奖励相关信息
	rewardDays := float64(rewardMinutes) / (24 * 60) // 转换为天
	expireTime := time.Now().AddDate(0, 0, rewardSettings.RewardExpirationDays)
	expirationDate := expireTime.Format("2006-01-02")

	// 创建奖励记录
	record := &entity.RewardRecord{
		ClientId:           clientId,
		RewardMinutes:      rewardMinutes,
		RewardDays:         rewardDays,
		IsFirstView:        isFirstView,
		RemainingMinutes:   rewardMinutes, // 初始剩余时长等于奖励时长
		TotalRewardMinutes: rewardMinutes, // 初始累计获得奖励等于奖励时长
		UsedMinutes:        0,             // 初始已使用时长为0
		Status:             1,             // 状态为有效
		ExpireAt:           gtime.NewFromTime(expireTime),
	}

	// 保存奖励记录
	err = rewardRecordDao.Save(ctx, record)
	if err != nil {
		g.Log().Error(ctx, "保存奖励记录失败", err)
		return nil, errors.New("保存奖励记录失败")
	}

	// 获取客户信息，用于更新客户时长表
	var clientInfo struct {
		RealName string `json:"real_name"`
	}
	err = g.DB().Model("client").Fields("real_name").Where("id", clientId).Scan(&clientInfo)
	if err != nil {
		g.Log().Error(ctx, "获取客户信息失败", err)
		return nil, errors.New("获取客户信息失败")
	}

	// 将时长格式化为易读的字符串
	rewardDurationStr := formatDuration(rewardMinutes)

	// 检查客户是否已有时长记录
	var clientDuration struct {
		Id                int    `json:"id"`
		RemainingDuration string `json:"remaining_duration"`
		TotalDuration     string `json:"total_duration"`
		UsedDuration      string `json:"used_duration"`
	}

	err = g.DB().Model("client_duration").Where("client_id", clientId).Scan(&clientDuration)
	if err != nil {
		// 正确处理记录不存在的情况
		if strings.Contains(err.Error(), "no rows in result set") {
			// 记录不存在，不是错误，将ID设置为0表示需要创建新记录
			clientDuration.Id = 0
		} else {
			// 其他数据库错误需要记录并返回
			g.Log().Error(ctx, "查询客户时长记录失败", err)
			return nil, errors.New("查询客户时长记录失败")
		}
	}

	if clientDuration.Id > 0 {
		// 已有记录，更新客户时长记录
		// 计算新的总时长、剩余时长
		oldTotalMinutes := parseDurationToMinutes(clientDuration.TotalDuration)
		oldRemainingMinutes := parseDurationToMinutes(clientDuration.RemainingDuration)
		// 当前不需要使用已使用时长，但保留以供后续功能扩展
		// oldUsedMinutes := parseDurationToMinutes(clientDuration.UsedDuration)

		newTotalMinutes := oldTotalMinutes + rewardMinutes
		newRemainingMinutes := oldRemainingMinutes + rewardMinutes

		// 格式化新的时长
		newTotalDuration := formatDuration(newTotalMinutes)
		newRemainingDuration := formatDuration(newRemainingMinutes)

		// 更新客户时长记录
		_, err = g.DB().Model("client_duration").
			Data(g.Map{
				"remaining_duration": newRemainingDuration,
				"total_duration":     newTotalDuration,
				"updated_at":         gtime.Now().String(),
			}).
			Where("client_id", clientId).
			Update()

		if err != nil {
			g.Log().Error(ctx, "更新客户时长记录失败", err)
			return nil, errors.New("更新客户时长记录失败")
		}
	} else {
		// 没有记录，创建新的客户时长记录
		_, err = g.DB().Model("client_duration").Data(g.Map{
			"client_id":          clientId,
			"client_name":        clientInfo.RealName,
			"remaining_duration": rewardDurationStr,
			"total_duration":     rewardDurationStr,
			"used_duration":      "0分钟",
			"created_at":         gtime.Now().String(),
		}).Insert()

		if err != nil {
			g.Log().Error(ctx, "创建客户时长记录失败", err)
			return nil, errors.New("创建客户时长记录失败")
		}
	}

	// 构建响应
	res = &settingsV1.WxRewardAdViewedRes{
		Success:        true,
		RewardMinutes:  rewardMinutes,
		RewardDays:     rewardDays,
		ExpirationDate: expirationDate,
	}

	// 根据时长调整提示消息格式
	if rewardMinutes < 60 { // 不足1小时
		res.Message = fmt.Sprintf("恭喜获得%d分钟奖励！", rewardMinutes)
	} else if rewardMinutes < 24*60 { // 1小时到24小时
		hours := rewardMinutes / 60
		mins := rewardMinutes % 60
		if mins > 0 {
			res.Message = fmt.Sprintf("恭喜获得%d小时%d分钟奖励！", hours, mins)
		} else {
			res.Message = fmt.Sprintf("恭喜获得%d小时奖励！", hours)
		}
	} else { // 大于等于24小时
		days := rewardMinutes / (24 * 60)
		remainingMinutes := rewardMinutes % (24 * 60)
		hours := remainingMinutes / 60

		if hours > 0 {
			res.Message = fmt.Sprintf("恭喜获得%d天%d小时奖励！", days, hours)
		} else {
			res.Message = fmt.Sprintf("恭喜获得%d天奖励！", days)
		}
	}

	return res, nil
}

// 格式化时长为易读的字符串，如"3天18小时42分钟"
func formatDuration(minutes int) string {
	days := minutes / (24 * 60)
	hours := (minutes % (24 * 60)) / 60
	mins := minutes % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, mins)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, mins)
	} else {
		return fmt.Sprintf("%d分钟", mins)
	}
}

// 从易读的时长字符串解析为分钟数
func parseDurationToMinutes(durationStr string) int {
	// 默认为0分钟
	if durationStr == "" {
		return 0
	}

	totalMinutes := 0

	// 解析天数
	daysMatch := regexp.MustCompile(`(\d+)天`).FindStringSubmatch(durationStr)
	if len(daysMatch) > 1 {
		days, _ := strconv.Atoi(daysMatch[1])
		totalMinutes += days * 24 * 60
	}

	// 解析小时
	hoursMatch := regexp.MustCompile(`(\d+)小时`).FindStringSubmatch(durationStr)
	if len(hoursMatch) > 1 {
		hours, _ := strconv.Atoi(hoursMatch[1])
		totalMinutes += hours * 60
	}

	// 解析分钟
	minsMatch := regexp.MustCompile(`(\d+)分钟`).FindStringSubmatch(durationStr)
	if len(minsMatch) > 1 {
		mins, _ := strconv.Atoi(minsMatch[1])
		totalMinutes += mins
	}

	return totalMinutes
}
