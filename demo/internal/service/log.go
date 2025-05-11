package service

import (
	"context"

	v1 "demo/api/log/v1"
)

// LogService 操作日志服务接口
type LogService interface {
	// List 获取操作日志列表
	List(ctx context.Context, req *v1.OperationLogListReq) (res *v1.OperationLogListRes, err error)

	// Delete 删除操作日志
	Delete(ctx context.Context, req *v1.OperationLogDeleteReq) (res *v1.OperationLogDeleteRes, err error)

	// Export 导出操作日志
	Export(ctx context.Context, req *v1.OperationLogExportReq) (res *v1.OperationLogExportRes, err error)

	// Record 记录操作日志
	Record(ctx context.Context, userId int, username, module, action string, result int, details string) error
}

var (
	localLog LogService
)

// Log 获取日志服务实例
func Log() LogService {
	if localLog == nil {
		panic("Log service not initialized")
	}
	return localLog
}

// SetLog 设置日志服务实例
func SetLog(s LogService) {
	localLog = s
}
