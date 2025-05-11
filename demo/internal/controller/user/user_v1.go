package user

import (
	"context"

	v1 "demo/api/user/v1"
	"demo/internal/logic/user"
)

// controllerV1 用户控制器V1版本
type controllerV1 struct{}

// List 获取用户列表
func (c *controllerV1) List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	return user.New().List(ctx, req)
}

// Create 创建用户
func (c *controllerV1) Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	return user.New().Create(ctx, req)
}

// Update 更新用户
func (c *controllerV1) Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	return user.New().Update(ctx, req)
}

// Delete 删除用户
func (c *controllerV1) Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	return user.New().Delete(ctx, req)
}
