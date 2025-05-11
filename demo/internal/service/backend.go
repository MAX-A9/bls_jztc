package service

import (
	"context"
	"sync"

	v1 "demo/api/backend/v1"
)

// BackendService 后端服务接口
type BackendService interface {
	// 文件管理
	FileUpload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error)
	FileList(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error)
	FileDetail(ctx context.Context, req *v1.FileDetailReq) (res *v1.FileDetailRes, err error)
	FileDelete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error)
	FileBatchDelete(ctx context.Context, req *v1.FileBatchDeleteReq) (res *v1.FileBatchDeleteRes, err error)
	FileUpdatePublic(ctx context.Context, req *v1.FileUpdatePublicReq) (res *v1.FileUpdatePublicRes, err error)
}

var (
	// 后端服务单例
	backendInstance BackendService
	// 后端服务初始化锁
	backendOnce sync.Once
)

// Backend 获取后端服务实例
func Backend() BackendService {
	if backendInstance == nil {
		panic("后端服务未初始化")
	}
	return backendInstance
}

// SetBackend 设置后端服务实例
func SetBackend(service BackendService) {
	backendOnce.Do(func() {
		backendInstance = service
	})
}
