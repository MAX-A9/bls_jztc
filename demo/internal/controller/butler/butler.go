package butler

import (
	"context"

	v1 "demo/api/butler/v1"
	"demo/internal/service"
)

var (
	Butler = cButler{}
)

type cButler struct{}

// Get 获取专属管家信息
func (c *cButler) Get(ctx context.Context, req *v1.ButlerReq) (res *v1.ButlerRes, err error) {
	// 初始化响应
	res = &v1.ButlerRes{}

	// 检查服务是否初始化
	if service.ButlerService() == nil {
		return res, nil
	}

	// 获取专属管家图片信息
	imageRes, err := service.ButlerService().GetImage(ctx, &v1.GetImageReq{})
	if err != nil {
		return res, nil // 出错时返回空结果，不中断请求
	}

	// 设置返回值（即使URL为空也可以安全返回）
	res.ImageUrl = imageRes.ImageUrl

	return res, nil
}

// SaveImage 保存专属管家图片
func (c *cButler) SaveImage(ctx context.Context, req *v1.SaveImageReq) (res *v1.SaveImageRes, err error) {
	return service.ButlerService().SaveImage(ctx, req)
}

// GetImage 客户端获取专属管家图片
func (c *cButler) GetImage(ctx context.Context, req *v1.GetImageReq) (res *v1.GetImageRes, err error) {
	// 初始化响应
	res = &v1.GetImageRes{}

	// 检查服务是否初始化
	if service.ButlerService() == nil {
		return res, nil
	}

	// 调用服务
	serviceRes, err := service.ButlerService().GetImage(ctx, req)
	if err != nil {
		return res, nil // 出错时返回空结果，不中断请求
	}

	// 拷贝数据
	if serviceRes != nil {
		res.ImageUrl = serviceRes.ImageUrl
	}

	return res, nil
}
