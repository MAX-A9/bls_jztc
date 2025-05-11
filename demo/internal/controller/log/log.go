package log

import (
	"context"

	v1 "demo/api/log/v1"
	"demo/internal/service"
)

// Controller 操作日志控制器
type Controller struct{}

// V1 创建V1版本操作日志控制器
func (c *Controller) V1() *controllerV1 {
	return &controllerV1{}
}

// New 创建一个操作日志控制器实例
func New() *Controller {
	return &Controller{}
}

// controllerV1 操作日志控制器V1版本
type controllerV1 struct{}

// List 获取操作日志列表
func (c *controllerV1) List(ctx context.Context, req *v1.OperationLogListReq) (res *v1.OperationLogListRes, err error) {
	return service.Log().List(ctx, req)
}

// Delete 删除操作日志
func (c *controllerV1) Delete(ctx context.Context, req *v1.OperationLogDeleteReq) (res *v1.OperationLogDeleteRes, err error) {
	return service.Log().Delete(ctx, req)
}

// Export 导出操作日志
func (c *controllerV1) Export(ctx context.Context, req *v1.OperationLogExportReq) (res *v1.OperationLogExportRes, err error) {
	return service.Log().Export(ctx, req)
}
