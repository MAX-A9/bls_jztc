package backend

import (
	"context"

	v1 "demo/api/backend/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
)

// ControllerV1 后端控制器V1版本
type ControllerV1 struct{}

// FileUpload 上传文件
func (c *ControllerV1) FileUpload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限上传文件")
	}

	return service.Backend().FileUpload(ctx, req)
}

// FileList 获取文件列表
func (c *ControllerV1) FileList(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看文件列表")
	}

	return service.Backend().FileList(ctx, req)
}

// FileDetail 获取文件详情
func (c *ControllerV1) FileDetail(ctx context.Context, req *v1.FileDetailReq) (res *v1.FileDetailRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限查看文件详情")
	}

	return service.Backend().FileDetail(ctx, req)
}

// FileDelete 删除文件
func (c *ControllerV1) FileDelete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限删除文件")
	}

	return service.Backend().FileDelete(ctx, req)
}

// FileBatchDelete 批量删除文件
func (c *ControllerV1) FileBatchDelete(ctx context.Context, req *v1.FileBatchDeleteReq) (res *v1.FileBatchDeleteRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限批量删除文件")
	}

	return service.Backend().FileBatchDelete(ctx, req)
}

// FileUpdatePublic 更新文件公开状态
func (c *ControllerV1) FileUpdatePublic(ctx context.Context, req *v1.FileUpdatePublicReq) (res *v1.FileUpdatePublicRes, err error) {
	// 权限检查
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.New("未登录或登录已过期")
	}
	if role != "admin" {
		return nil, gerror.New("您没有权限修改文件状态")
	}

	return service.Backend().FileUpdatePublic(ctx, req)
}

// NewV1 创建后端控制器V1版本实例
func NewV1() *ControllerV1 {
	return &ControllerV1{}
}
