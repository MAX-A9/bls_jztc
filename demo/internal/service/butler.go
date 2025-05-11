// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"demo/api/butler/v1"
)

// Butler 管家服务接口定义
type Butler interface {
	// SaveImage 保存管家图片
	SaveImage(ctx context.Context, req *v1.SaveImageReq) (res *v1.SaveImageRes, err error)
	// GetImage 获取管家图片
	GetImage(ctx context.Context, req *v1.GetImageReq) (res *v1.GetImageRes, err error)
}

var (
	localButler Butler
)

// SetButler 设置管家服务实例
func SetButler(s Butler) {
	g.Log().Debug(context.Background(), "service.SetButler")
	localButler = s
}

// Butler 获取管家服务实例
func ButlerService() Butler {
	return localButler
} 