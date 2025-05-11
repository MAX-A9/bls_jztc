// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package backend

import (
	"context"

	"demo/api/backend/v1"
)

// IBackendV1 后端接口V1版本
type IBackendV1 interface {
	// 文件管理
	FileUpload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error)
	FileList(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error)
	FileDetail(ctx context.Context, req *v1.FileDetailReq) (res *v1.FileDetailRes, err error)
	FileDelete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error)
	FileBatchDelete(ctx context.Context, req *v1.FileBatchDeleteReq) (res *v1.FileBatchDeleteRes, err error)
	FileUpdatePublic(ctx context.Context, req *v1.FileUpdatePublicReq) (res *v1.FileUpdatePublicRes, err error)
} 