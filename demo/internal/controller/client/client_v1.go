package client

import (
	"context"

	v1 "demo/api/client/v1"
)

// List 获取客户列表
func (c *ControllerV1) List(ctx context.Context, req *v1.ClientListReq) (res *v1.ClientListRes, err error) {
	return c.service.List(ctx, req)
}

// Create 创建客户
func (c *ControllerV1) Create(ctx context.Context, req *v1.ClientCreateReq) (res *v1.ClientCreateRes, err error) {
	return c.service.Create(ctx, req)
}

// Update 更新客户
func (c *ControllerV1) Update(ctx context.Context, req *v1.ClientUpdateReq) (res *v1.ClientUpdateRes, err error) {
	return c.service.Update(ctx, req)
}

// Delete 删除客户
func (c *ControllerV1) Delete(ctx context.Context, req *v1.ClientDeleteReq) (res *v1.ClientDeleteRes, err error) {
	return c.service.Delete(ctx, req)
}

// WxappLogin 微信小程序登录
func (c *ControllerV1) WxappLogin(ctx context.Context, req *v1.WxappLoginReq) (res *v1.WxappLoginRes, err error) {
	return c.service.WxappLogin(ctx, req)
}

// Info 获取客户信息
func (c *ControllerV1) Info(ctx context.Context, req *v1.ClientInfoReq) (res *v1.ClientInfoRes, err error) {
	return c.service.Info(ctx, req)
}

// UpdateProfile 更新客户个人信息
func (c *ControllerV1) UpdateProfile(ctx context.Context, req *v1.ClientUpdateProfileReq) (res *v1.ClientUpdateProfileRes, err error) {
	return c.service.UpdateProfile(ctx, req)
}

// DurationList 获取客户时长列表
func (c *ControllerV1) DurationList(ctx context.Context, req *v1.ClientDurationListReq) (res *v1.ClientDurationListRes, err error) {
	return c.service.DurationList(ctx, req)
}

// DurationDetail 获取客户时长详情
func (c *ControllerV1) DurationDetail(ctx context.Context, req *v1.ClientDurationDetailReq) (res *v1.ClientDurationDetailRes, err error) {
	return c.service.DurationDetail(ctx, req)
}

// DurationCreate 创建客户时长
func (c *ControllerV1) DurationCreate(ctx context.Context, req *v1.ClientDurationCreateReq) (res *v1.ClientDurationCreateRes, err error) {
	return c.service.DurationCreate(ctx, req)
}

// DurationUpdate 更新客户时长
func (c *ControllerV1) DurationUpdate(ctx context.Context, req *v1.ClientDurationUpdateReq) (res *v1.ClientDurationUpdateRes, err error) {
	return c.service.DurationUpdate(ctx, req)
}

// DurationDelete 删除客户时长
func (c *ControllerV1) DurationDelete(ctx context.Context, req *v1.ClientDurationDeleteReq) (res *v1.ClientDurationDeleteRes, err error) {
	return c.service.DurationDelete(ctx, req)
}

// WxGetRemainingDuration 获取客户端用户剩余时长
func (c *ControllerV1) WxGetRemainingDuration(ctx context.Context, req *v1.WxClientRemainingDurationReq) (res *v1.WxClientRemainingDurationRes, err error) {
	return c.service.WxGetRemainingDuration(ctx, req)
}
