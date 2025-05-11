package service

import (
	"context"
	v1 "demo/api/content/v1"
	settingsV1 "demo/api/settings/v1"
)

// ISettings 城市系统基础设置服务接口
type ISettings interface {
	// GetMiniProgramList 获取导航小程序列表
	GetMiniProgramList(ctx context.Context, req *v1.MiniProgramListReq) (res *v1.MiniProgramListRes, err error)

	// CreateMiniProgram 创建导航小程序
	CreateMiniProgram(ctx context.Context, req *v1.MiniProgramCreateReq) (res *v1.MiniProgramCreateRes, err error)

	// UpdateMiniProgram 更新导航小程序
	UpdateMiniProgram(ctx context.Context, req *v1.MiniProgramUpdateReq) (res *v1.MiniProgramUpdateRes, err error)

	// DeleteMiniProgram 删除导航小程序
	DeleteMiniProgram(ctx context.Context, req *v1.MiniProgramDeleteReq) (res *v1.MiniProgramDeleteRes, err error)

	// UpdateMiniProgramStatus 更新导航小程序状态
	UpdateMiniProgramStatus(ctx context.Context, req *v1.MiniProgramStatusUpdateReq) (res *v1.MiniProgramStatusUpdateRes, err error)

	// UpdateMiniProgramGlobalStatus 更新导航小程序总开关状态
	UpdateMiniProgramGlobalStatus(ctx context.Context, req *v1.MiniProgramGlobalStatusUpdateReq) (res *v1.MiniProgramGlobalStatusUpdateRes, err error)

	// GetMiniProgramBaseSettings 获取小程序基础设置
	GetMiniProgramBaseSettings(ctx context.Context, req *settingsV1.MiniProgramBaseSettingsReq) (res *settingsV1.MiniProgramBaseSettingsRes, err error)

	// SaveMiniProgramBaseSettings 保存小程序基础设置
	SaveMiniProgramBaseSettings(ctx context.Context, req *settingsV1.MiniProgramBaseSettingsSaveReq) (res *settingsV1.MiniProgramBaseSettingsSaveRes, err error)

	// GetAdSettings 获取广告设置
	GetAdSettings(ctx context.Context, req *settingsV1.AdSettingsReq) (res *settingsV1.AdSettingsRes, err error)

	// SaveAdSettings 保存广告设置
	SaveAdSettings(ctx context.Context, req *settingsV1.AdSettingsSaveReq) (res *settingsV1.AdSettingsSaveRes, err error)

	// GetRewardSettings 获取奖励设置
	GetRewardSettings(ctx context.Context, req *settingsV1.RewardSettingsReq) (res *settingsV1.RewardSettingsRes, err error)

	// SaveRewardSettings 保存奖励设置
	SaveRewardSettings(ctx context.Context, req *settingsV1.RewardSettingsSaveReq) (res *settingsV1.RewardSettingsSaveRes, err error)

	// GetAgreementSettings 获取协议设置
	GetAgreementSettings(ctx context.Context, req *settingsV1.AgreementSettingsReq) (res *settingsV1.AgreementSettingsRes, err error)

	// SaveAgreementSettings 保存协议设置
	SaveAgreementSettings(ctx context.Context, req *settingsV1.AgreementSettingsSaveReq) (res *settingsV1.AgreementSettingsSaveRes, err error)

	// WxGetAgreement 微信客户端获取协议
	WxGetAgreement(ctx context.Context, req *settingsV1.WxAgreementGetReq) (res *settingsV1.WxAgreementGetRes, err error)

	// GetBannerList 获取轮播图列表
	GetBannerList(ctx context.Context, req *v1.BannerListReq) (res *v1.BannerListRes, err error)

	// CreateBanner 创建轮播图
	CreateBanner(ctx context.Context, req *v1.BannerCreateReq) (res *v1.BannerCreateRes, err error)

	// UpdateBanner 更新轮播图
	UpdateBanner(ctx context.Context, req *v1.BannerUpdateReq) (res *v1.BannerUpdateRes, err error)

	// DeleteBanner 删除轮播图
	DeleteBanner(ctx context.Context, req *v1.BannerDeleteReq) (res *v1.BannerDeleteRes, err error)

	// UpdateBannerStatus 更新轮播图状态
	UpdateBannerStatus(ctx context.Context, req *v1.BannerStatusUpdateReq) (res *v1.BannerStatusUpdateRes, err error)

	// GetActivityArea 获取活动区域
	GetActivityArea(ctx context.Context, req *v1.ActivityAreaGetReq) (res *v1.ActivityAreaGetRes, err error)

	// SaveActivityArea 保存活动区域
	SaveActivityArea(ctx context.Context, req *v1.ActivityAreaSaveReq) (res *v1.ActivityAreaSaveRes, err error)

	// WxGetActivityArea 微信客户端获取活动区域
	WxGetActivityArea(ctx context.Context, req *v1.WxActivityAreaGetReq) (res *v1.WxActivityAreaGetRes, err error)

	// UpdateBannerGlobalStatus 更新轮播图总开关状态
	UpdateBannerGlobalStatus(ctx context.Context, req *v1.BannerGlobalStatusUpdateReq) (res *v1.BannerGlobalStatusUpdateRes, err error)

	// UpdateActivityAreaGlobalStatus 更新活动区域总开关状态
	UpdateActivityAreaGlobalStatus(ctx context.Context, req *v1.ActivityAreaGlobalStatusUpdateReq) (res *v1.ActivityAreaGlobalStatusUpdateRes, err error)

	// WxGetMiniProgramList 微信客户端获取导航小程序列表
	WxGetMiniProgramList(ctx context.Context, req *v1.WxMiniProgramListReq) (res *v1.WxMiniProgramListRes, err error)

	// WxGetBannerList 微信客户端获取轮播图列表
	WxGetBannerList(ctx context.Context, req *v1.WxBannerListReq) (res *v1.WxBannerListRes, err error)

	// GetInnerBannerList 获取内页轮播图列表
	GetInnerBannerList(ctx context.Context, req *v1.InnerBannerListReq) (res *v1.InnerBannerListRes, err error)

	// CreateInnerBanner 创建内页轮播图
	CreateInnerBanner(ctx context.Context, req *v1.InnerBannerCreateReq) (res *v1.InnerBannerCreateRes, err error)

	// UpdateInnerBanner 更新内页轮播图
	UpdateInnerBanner(ctx context.Context, req *v1.InnerBannerUpdateReq) (res *v1.InnerBannerUpdateRes, err error)

	// DeleteInnerBanner 删除内页轮播图
	DeleteInnerBanner(ctx context.Context, req *v1.InnerBannerDeleteReq) (res *v1.InnerBannerDeleteRes, err error)

	// UpdateInnerBannerStatus 更新内页轮播图状态
	UpdateInnerBannerStatus(ctx context.Context, req *v1.InnerBannerStatusUpdateReq) (res *v1.InnerBannerStatusUpdateRes, err error)

	// UpdateHomeInnerBannerGlobalStatus 更新首页内页轮播图总开关状态
	UpdateHomeInnerBannerGlobalStatus(ctx context.Context, req *v1.HomeInnerBannerGlobalStatusUpdateReq) (res *v1.HomeInnerBannerGlobalStatusUpdateRes, err error)

	// UpdateIdleInnerBannerGlobalStatus 更新闲置页内页轮播图总开关状态
	UpdateIdleInnerBannerGlobalStatus(ctx context.Context, req *v1.IdleInnerBannerGlobalStatusUpdateReq) (res *v1.IdleInnerBannerGlobalStatusUpdateRes, err error)

	// WxGetInnerBannerList 微信客户端获取内页轮播图列表
	WxGetInnerBannerList(ctx context.Context, req *v1.WxInnerBannerListReq) (res *v1.WxInnerBannerListRes, err error)

	// WxGetMiniProgramBaseSettings 微信客户端获取小程序基础设置
	WxGetMiniProgramBaseSettings(ctx context.Context, req *settingsV1.WxMiniProgramBaseSettingsReq) (res *settingsV1.WxMiniProgramBaseSettingsRes, err error)

	// WxGetAdSettings 微信客户端获取广告设置
	WxGetAdSettings(ctx context.Context, req *settingsV1.WxGetAdSettingsReq) (res *settingsV1.WxGetAdSettingsRes, err error)

	// WxRewardAdViewed 微信客户端广告观看完成
	WxRewardAdViewed(ctx context.Context, req *settingsV1.WxRewardAdViewedReq) (res *settingsV1.WxRewardAdViewedRes, err error)

	// GetShareSettings 获取分享设置
	GetShareSettings(ctx context.Context, req *settingsV1.ShareSettingsReq) (res *settingsV1.ShareSettingsRes, err error)

	// SaveShareSettings 保存分享设置
	SaveShareSettings(ctx context.Context, req *settingsV1.SaveShareSettingsReq) (res *settingsV1.SaveShareSettingsRes, err error)

	// WxGetShareSettings 微信客户端获取分享设置
	WxGetShareSettings(ctx context.Context, req *settingsV1.WxShareSettingsReq) (res *settingsV1.WxShareSettingsRes, err error)
}
