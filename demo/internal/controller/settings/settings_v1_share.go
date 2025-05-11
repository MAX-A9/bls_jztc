package settings

import (
	"context"
	v1 "demo/api/settings/v1"
	"demo/internal/service"
)

// GetShareSettings 获取分享设置
func (c *ControllerV1) GetShareSettings(ctx context.Context, req *v1.ShareSettingsReq) (res *v1.ShareSettingsRes, err error) {
	return service.Settings().GetShareSettings(ctx, req)
}

// SaveShareSettings 保存分享设置
func (c *ControllerV1) SaveShareSettings(ctx context.Context, req *v1.SaveShareSettingsReq) (res *v1.SaveShareSettingsRes, err error) {
	return service.Settings().SaveShareSettings(ctx, req)
}

// WxGetShareSettings 微信客户端获取分享设置
func (c *ControllerV1) WxGetShareSettings(ctx context.Context, req *v1.WxShareSettingsReq) (res *v1.WxShareSettingsRes, err error) {
	return service.Settings().WxGetShareSettings(ctx, req)
}
