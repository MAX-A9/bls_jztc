package storage

import (
	"context"

	v1 "demo/api/storage/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// ControllerV1 存储控制器V1版本
type ControllerV1 struct {
	Client *controllerClient
}

// GetConfig 获取存储配置
func (c *ControllerV1) GetConfig(ctx context.Context, req *v1.StorageConfigReq) (res *v1.StorageConfigRes, err error) {
	// 检查权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil || role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限查看存储配置", nil))
	}

	return service.Storage().GetConfig(ctx, req)
}

// SaveConfig 保存存储配置
func (c *ControllerV1) SaveConfig(ctx context.Context, req *v1.SaveStorageConfigReq) (res *v1.SaveStorageConfigRes, err error) {
	// 权限验证：只有管理员才能修改配置
	userId, username, role, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(gcode.New(403, "未登录或登录已过期", nil))
	}
	if role != "admin" {
		return nil, gerror.NewCode(gcode.New(403, "您没有权限修改存储配置", nil))
	}

	// 调用服务实现保存配置
	res, err = service.Storage().SaveConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	// 记录操作日志
	service.Log().Record(ctx, userId, username, "存储设置", "修改存储配置", 1, "")

	return res, nil
}
