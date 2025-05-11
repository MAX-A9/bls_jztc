package settings

import (
	"context"
	v1 "demo/api/content/v1"
	settingsV1 "demo/api/settings/v1"
	"demo/internal/service"
)

// Controller 城市系统基础设置控制器
type Controller struct{}

// V1 创建V1版本控制器
func (c *Controller) V1() *ControllerV1 {
	return &ControllerV1{}
}

// ControllerV1 V1版本城市系统基础设置控制器
type ControllerV1 struct{}

// GetMiniProgramList 获取导航小程序列表
func (c *ControllerV1) GetMiniProgramList(ctx context.Context, req *v1.MiniProgramListReq) (res *v1.MiniProgramListRes, err error) {
	return service.Settings().GetMiniProgramList(ctx, req)
}

// CreateMiniProgram 创建导航小程序
func (c *ControllerV1) CreateMiniProgram(ctx context.Context, req *v1.MiniProgramCreateReq) (res *v1.MiniProgramCreateRes, err error) {
	return service.Settings().CreateMiniProgram(ctx, req)
}

// UpdateMiniProgram 更新导航小程序
func (c *ControllerV1) UpdateMiniProgram(ctx context.Context, req *v1.MiniProgramUpdateReq) (res *v1.MiniProgramUpdateRes, err error) {
	return service.Settings().UpdateMiniProgram(ctx, req)
}

// DeleteMiniProgram 删除导航小程序
func (c *ControllerV1) DeleteMiniProgram(ctx context.Context, req *v1.MiniProgramDeleteReq) (res *v1.MiniProgramDeleteRes, err error) {
	return service.Settings().DeleteMiniProgram(ctx, req)
}

// UpdateMiniProgramStatus 更新导航小程序状态
func (c *ControllerV1) UpdateMiniProgramStatus(ctx context.Context, req *v1.MiniProgramStatusUpdateReq) (res *v1.MiniProgramStatusUpdateRes, err error) {
	return service.Settings().UpdateMiniProgramStatus(ctx, req)
}

// UpdateMiniProgramGlobalStatus 更新导航小程序总开关状态
func (c *ControllerV1) UpdateMiniProgramGlobalStatus(ctx context.Context, req *v1.MiniProgramGlobalStatusUpdateReq) (res *v1.MiniProgramGlobalStatusUpdateRes, err error) {
	return service.Settings().UpdateMiniProgramGlobalStatus(ctx, req)
}

// GetMiniProgramBaseSettings 获取小程序基础设置
func (c *ControllerV1) GetMiniProgramBaseSettings(ctx context.Context, req *settingsV1.MiniProgramBaseSettingsReq) (res *settingsV1.MiniProgramBaseSettingsRes, err error) {
	return service.Settings().GetMiniProgramBaseSettings(ctx, req)
}

// SaveMiniProgramBaseSettings 保存小程序基础设置
func (c *ControllerV1) SaveMiniProgramBaseSettings(ctx context.Context, req *settingsV1.MiniProgramBaseSettingsSaveReq) (res *settingsV1.MiniProgramBaseSettingsSaveRes, err error) {
	return service.Settings().SaveMiniProgramBaseSettings(ctx, req)
}

// GetAdSettings 获取广告设置
func (c *ControllerV1) GetAdSettings(ctx context.Context, req *settingsV1.AdSettingsReq) (res *settingsV1.AdSettingsRes, err error) {
	return service.Settings().GetAdSettings(ctx, req)
}

// SaveAdSettings 保存广告设置
func (c *ControllerV1) SaveAdSettings(ctx context.Context, req *settingsV1.AdSettingsSaveReq) (res *settingsV1.AdSettingsSaveRes, err error) {
	return service.Settings().SaveAdSettings(ctx, req)
}

// GetRewardSettings 获取奖励设置
func (c *ControllerV1) GetRewardSettings(ctx context.Context, req *settingsV1.RewardSettingsReq) (res *settingsV1.RewardSettingsRes, err error) {
	return service.Settings().GetRewardSettings(ctx, req)
}

// SaveRewardSettings 保存奖励设置
func (c *ControllerV1) SaveRewardSettings(ctx context.Context, req *settingsV1.RewardSettingsSaveReq) (res *settingsV1.RewardSettingsSaveRes, err error) {
	return service.Settings().SaveRewardSettings(ctx, req)
}

// GetAgreementSettings 获取协议设置
func (c *ControllerV1) GetAgreementSettings(ctx context.Context, req *settingsV1.AgreementSettingsReq) (res *settingsV1.AgreementSettingsRes, err error) {
	return service.Settings().GetAgreementSettings(ctx, req)
}

// SaveAgreementSettings 保存协议设置
func (c *ControllerV1) SaveAgreementSettings(ctx context.Context, req *settingsV1.AgreementSettingsSaveReq) (res *settingsV1.AgreementSettingsSaveRes, err error) {
	return service.Settings().SaveAgreementSettings(ctx, req)
}

// WxGetAgreement 微信客户端获取协议
func (c *ControllerV1) WxGetAgreement(ctx context.Context, req *settingsV1.WxAgreementGetReq) (res *settingsV1.WxAgreementGetRes, err error) {
	return service.Settings().WxGetAgreement(ctx, req)
}

// GetBannerList 获取轮播图列表
func (c *ControllerV1) GetBannerList(ctx context.Context, req *v1.BannerListReq) (res *v1.BannerListRes, err error) {
	return service.Settings().GetBannerList(ctx, req)
}

// CreateBanner 创建轮播图
func (c *ControllerV1) CreateBanner(ctx context.Context, req *v1.BannerCreateReq) (res *v1.BannerCreateRes, err error) {
	return service.Settings().CreateBanner(ctx, req)
}

// UpdateBanner 更新轮播图
func (c *ControllerV1) UpdateBanner(ctx context.Context, req *v1.BannerUpdateReq) (res *v1.BannerUpdateRes, err error) {
	return service.Settings().UpdateBanner(ctx, req)
}

// DeleteBanner 删除轮播图
func (c *ControllerV1) DeleteBanner(ctx context.Context, req *v1.BannerDeleteReq) (res *v1.BannerDeleteRes, err error) {
	return service.Settings().DeleteBanner(ctx, req)
}

// UpdateBannerStatus 更新轮播图状态
func (c *ControllerV1) UpdateBannerStatus(ctx context.Context, req *v1.BannerStatusUpdateReq) (res *v1.BannerStatusUpdateRes, err error) {
	return service.Settings().UpdateBannerStatus(ctx, req)
}

// UpdateBannerGlobalStatus 更新轮播图总开关状态
func (c *ControllerV1) UpdateBannerGlobalStatus(ctx context.Context, req *v1.BannerGlobalStatusUpdateReq) (res *v1.BannerGlobalStatusUpdateRes, err error) {
	return service.Settings().UpdateBannerGlobalStatus(ctx, req)
}

// GetActivityArea 获取活动区域
func (c *ControllerV1) GetActivityArea(ctx context.Context, req *v1.ActivityAreaGetReq) (res *v1.ActivityAreaGetRes, err error) {
	return service.Settings().GetActivityArea(ctx, req)
}

// SaveActivityArea 保存活动区域
func (c *ControllerV1) SaveActivityArea(ctx context.Context, req *v1.ActivityAreaSaveReq) (res *v1.ActivityAreaSaveRes, err error) {
	return service.Settings().SaveActivityArea(ctx, req)
}

// WxGetActivityArea 微信客户端获取活动区域
func (c *ControllerV1) WxGetActivityArea(ctx context.Context, req *v1.WxActivityAreaGetReq) (res *v1.WxActivityAreaGetRes, err error) {
	return service.Settings().WxGetActivityArea(ctx, req)
}

// UpdateActivityAreaGlobalStatus 更新活动区域总开关状态
func (c *ControllerV1) UpdateActivityAreaGlobalStatus(ctx context.Context, req *v1.ActivityAreaGlobalStatusUpdateReq) (res *v1.ActivityAreaGlobalStatusUpdateRes, err error) {
	return service.Settings().UpdateActivityAreaGlobalStatus(ctx, req)
}

// WxGetMiniProgramList 微信客户端获取导航小程序列表
func (c *ControllerV1) WxGetMiniProgramList(ctx context.Context, req *v1.WxMiniProgramListReq) (res *v1.WxMiniProgramListRes, err error) {
	return service.Settings().WxGetMiniProgramList(ctx, req)
}

// WxGetBannerList 微信客户端获取轮播图列表
func (c *ControllerV1) WxGetBannerList(ctx context.Context, req *v1.WxBannerListReq) (res *v1.WxBannerListRes, err error) {
	return service.Settings().WxGetBannerList(ctx, req)
}

// GetInnerBannerList 获取内页轮播图列表
func (c *ControllerV1) GetInnerBannerList(ctx context.Context, req *v1.InnerBannerListReq) (res *v1.InnerBannerListRes, err error) {
	return service.Settings().GetInnerBannerList(ctx, req)
}

// CreateInnerBanner 创建内页轮播图
func (c *ControllerV1) CreateInnerBanner(ctx context.Context, req *v1.InnerBannerCreateReq) (res *v1.InnerBannerCreateRes, err error) {
	return service.Settings().CreateInnerBanner(ctx, req)
}

// UpdateInnerBanner 更新内页轮播图
func (c *ControllerV1) UpdateInnerBanner(ctx context.Context, req *v1.InnerBannerUpdateReq) (res *v1.InnerBannerUpdateRes, err error) {
	return service.Settings().UpdateInnerBanner(ctx, req)
}

// DeleteInnerBanner 删除内页轮播图
func (c *ControllerV1) DeleteInnerBanner(ctx context.Context, req *v1.InnerBannerDeleteReq) (res *v1.InnerBannerDeleteRes, err error) {
	return service.Settings().DeleteInnerBanner(ctx, req)
}

// UpdateInnerBannerStatus 更新内页轮播图状态
func (c *ControllerV1) UpdateInnerBannerStatus(ctx context.Context, req *v1.InnerBannerStatusUpdateReq) (res *v1.InnerBannerStatusUpdateRes, err error) {
	return service.Settings().UpdateInnerBannerStatus(ctx, req)
}

// UpdateHomeInnerBannerGlobalStatus 更新首页内页轮播图总开关状态
func (c *ControllerV1) UpdateHomeInnerBannerGlobalStatus(ctx context.Context, req *v1.HomeInnerBannerGlobalStatusUpdateReq) (res *v1.HomeInnerBannerGlobalStatusUpdateRes, err error) {
	return service.Settings().UpdateHomeInnerBannerGlobalStatus(ctx, req)
}

// UpdateIdleInnerBannerGlobalStatus 更新闲置页内页轮播图总开关状态
func (c *ControllerV1) UpdateIdleInnerBannerGlobalStatus(ctx context.Context, req *v1.IdleInnerBannerGlobalStatusUpdateReq) (res *v1.IdleInnerBannerGlobalStatusUpdateRes, err error) {
	return service.Settings().UpdateIdleInnerBannerGlobalStatus(ctx, req)
}

// WxGetInnerBannerList 微信客户端获取内页轮播图列表
func (c *ControllerV1) WxGetInnerBannerList(ctx context.Context, req *v1.WxInnerBannerListReq) (res *v1.WxInnerBannerListRes, err error) {
	return service.Settings().WxGetInnerBannerList(ctx, req)
}

// WxGetMiniProgramBaseSettings 微信客户端获取小程序基础设置
func (c *ControllerV1) WxGetMiniProgramBaseSettings(ctx context.Context, req *settingsV1.WxMiniProgramBaseSettingsReq) (res *settingsV1.WxMiniProgramBaseSettingsRes, err error) {
	return service.Settings().WxGetMiniProgramBaseSettings(ctx, req)
}

// WxGetAdSettings 微信客户端获取广告设置
func (c *ControllerV1) WxGetAdSettings(ctx context.Context, req *settingsV1.WxGetAdSettingsReq) (res *settingsV1.WxGetAdSettingsRes, err error) {
	return service.Settings().WxGetAdSettings(ctx, req)
}

// WxRewardAdViewed 微信客户端广告观看完成
func (c *ControllerV1) WxRewardAdViewed(ctx context.Context, req *settingsV1.WxRewardAdViewedReq) (res *settingsV1.WxRewardAdViewedRes, err error) {
	return service.Settings().WxRewardAdViewed(ctx, req)
}

// New 创建城市系统基础设置控制器实例
func New() *Controller {
	return &Controller{}
}
