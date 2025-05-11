package settings

import (
	"context"

	v1 "demo/api/settings/v1"
	"demo/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// GetShareSettings 获取分享设置
func (s *sSettings) GetShareSettings(ctx context.Context, req *v1.ShareSettingsReq) (res *v1.ShareSettingsRes, err error) {
	res = &v1.ShareSettingsRes{}

	// 使用ShareSettingsDao获取分享设置
	settings, err := s.shareSettingsDao.Get(ctx)
	if err != nil {
		// 当分享设置不存在时，返回错误信息
		return nil, gerror.New("分享设置数据不存在，请先在管理后台配置")
	}

	// 将获取到的分享设置转换为响应格式
	res.DefaultShareText = settings.DefaultShareText
	res.DefaultShareImage = settings.DefaultShareImage
	res.ContentShareText = settings.ContentShareText
	res.ContentShareImage = settings.ContentShareImage
	res.HomeShareText = settings.HomeShareText
	res.HomeShareImage = settings.HomeShareImage

	return res, nil
}

// SaveShareSettings 保存分享设置
func (s *sSettings) SaveShareSettings(ctx context.Context, req *v1.SaveShareSettingsReq) (res *v1.SaveShareSettingsRes, err error) {
	res = &v1.SaveShareSettingsRes{Success: false}

	// 创建分享设置实体
	settings := &entity.ShareSettings{
		Id:                1,
		DefaultShareText:  req.DefaultShareText,
		DefaultShareImage: req.DefaultShareImage,
		ContentShareText:  req.ContentShareText,
		ContentShareImage: req.ContentShareImage,
		HomeShareText:     req.HomeShareText,
		HomeShareImage:    req.HomeShareImage,
	}

	// 保存分享设置
	if err := s.shareSettingsDao.Save(ctx, settings); err != nil {
		g.Log().Error(ctx, "保存分享设置失败", err)
		return nil, gerror.New("保存分享设置失败")
	}

	res.Success = true
	return res, nil
}

// WxGetShareSettings 微信客户端获取分享设置
func (s *sSettings) WxGetShareSettings(ctx context.Context, req *v1.WxShareSettingsReq) (res *v1.WxShareSettingsRes, err error) {
	res = &v1.WxShareSettingsRes{}

	// 复用管理端的获取分享设置方法
	settingsRes, err := s.GetShareSettings(ctx, &v1.ShareSettingsReq{})
	if err != nil {
		return nil, err
	}

	// 将获取到的分享设置转换为微信客户端响应格式
	res.DefaultShareText = settingsRes.DefaultShareText
	res.DefaultShareImage = settingsRes.DefaultShareImage
	res.ContentShareText = settingsRes.ContentShareText
	res.ContentShareImage = settingsRes.ContentShareImage
	res.HomeShareText = settingsRes.HomeShareText
	res.HomeShareImage = settingsRes.HomeShareImage

	return res, nil
}
