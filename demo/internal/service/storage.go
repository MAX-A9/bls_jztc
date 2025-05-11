package service

import (
	"context"
	"sync"

	v1 "demo/api/storage/v1"
)

// StorageService 存储服务接口
type StorageService interface {
	// GetConfig 获取存储配置
	GetConfig(ctx context.Context, req *v1.StorageConfigReq) (res *v1.StorageConfigRes, err error)

	// SaveConfig 保存存储配置
	SaveConfig(ctx context.Context, req *v1.SaveStorageConfigReq) (res *v1.SaveStorageConfigRes, err error)

	// WxUploadImage 微信客户端-上传图片
	WxUploadImage(ctx context.Context, req *v1.WxUploadImageReq) (res *v1.WxUploadImageRes, err error)
}

var (
	storageServiceInstance StorageService
	storageOnce            sync.Once
)

// Storage 获取存储服务实例
func Storage() StorageService {
	if storageServiceInstance == nil {
		panic("存储服务未初始化")
	}
	return storageServiceInstance
}

// SetStorage 设置存储服务实例
func SetStorage(service StorageService) {
	storageOnce.Do(func() {
		storageServiceInstance = service
	})
}
