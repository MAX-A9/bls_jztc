package storage

import (
	"context"

	v1 "demo/api/storage/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
)

// controllerClient 客户端存储控制器
type controllerClient struct{}

// WxUploadImage 微信客户端-上传图片
func (c *controllerClient) WxUploadImage(ctx context.Context, req *v1.WxUploadImageReq) (res *v1.WxUploadImageRes, err error) {
	// 权限检查
	_, err = auth.GetClientInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}

	return service.Storage().WxUploadImage(ctx, req)
}

// Client 客户端存储控制器实例
var Client = &controllerClient{}
