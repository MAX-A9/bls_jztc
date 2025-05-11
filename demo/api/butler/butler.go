package butler

import (
	"context"

	v1 "demo/api/butler/v1"
)

// 管理端接口
type ButlerController interface {
	// 保存专属管家图片
	SaveImage(ctx context.Context, req *v1.SaveImageReq) (res *v1.SaveImageRes, err error)
}

// 客户端接口
type ClientButlerController interface {
	// 获取专属管家图片
	GetImage(ctx context.Context, req *v1.GetImageReq) (res *v1.GetImageRes, err error)
}
