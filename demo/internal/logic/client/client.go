package client

import (
	"context"

	v1 "demo/api/client/v1"
	"demo/internal/service"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
)

type sClient struct{}

func New() service.ClientService {
	return &sClient{}
}

// List 获取客户列表
func (s *sClient) List(ctx context.Context, req *v1.ClientListReq) (res *v1.ClientListRes, err error) {
	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.List(ctx, req)
}

// Create 创建客户
func (s *sClient) Create(ctx context.Context, req *v1.ClientCreateReq) (res *v1.ClientCreateRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.Create(ctx, req)
}

// Update 更新客户
func (s *sClient) Update(ctx context.Context, req *v1.ClientUpdateReq) (res *v1.ClientUpdateRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.Update(ctx, req)
}

// Delete 删除客户
func (s *sClient) Delete(ctx context.Context, req *v1.ClientDeleteReq) (res *v1.ClientDeleteRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.Delete(ctx, req)
}

// WxappLogin 微信小程序登录
func (s *sClient) WxappLogin(ctx context.Context, req *v1.WxappLoginReq) (res *v1.WxappLoginRes, err error) {
	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.WxappLogin(ctx, req)
}

// GetWxappConfig 获取微信小程序配置
func (s *sClient) GetWxappConfig(ctx context.Context, req *v1.WxappConfigGetReq) (res *v1.WxappConfigGetRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.GetWxappConfig(ctx, req)
}

// SaveWxappConfig 保存微信小程序配置
func (s *sClient) SaveWxappConfig(ctx context.Context, req *v1.WxappConfigSaveReq) (res *v1.WxappConfigSaveRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.SaveWxappConfig(ctx, req)
}

// Info 获取客户信息
func (s *sClient) Info(ctx context.Context, req *v1.ClientInfoReq) (res *v1.ClientInfoRes, err error) {
	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.Info(ctx, req)
}

// UpdateProfile 更新客户个人信息
func (s *sClient) UpdateProfile(ctx context.Context, req *v1.ClientUpdateProfileReq) (res *v1.ClientUpdateProfileRes, err error) {
	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.UpdateProfile(ctx, req)
}

// DurationList 获取客户时长列表
func (s *sClient) DurationList(ctx context.Context, req *v1.ClientDurationListReq) (res *v1.ClientDurationListRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.DurationList(ctx, req)
}

// DurationDetail 获取客户时长详情
func (s *sClient) DurationDetail(ctx context.Context, req *v1.ClientDurationDetailReq) (res *v1.ClientDurationDetailRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.DurationDetail(ctx, req)
}

// DurationCreate 创建客户时长
func (s *sClient) DurationCreate(ctx context.Context, req *v1.ClientDurationCreateReq) (res *v1.ClientDurationCreateRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.DurationCreate(ctx, req)
}

// DurationUpdate 更新客户时长
func (s *sClient) DurationUpdate(ctx context.Context, req *v1.ClientDurationUpdateReq) (res *v1.ClientDurationUpdateRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.DurationUpdate(ctx, req)
}

// DurationDelete 删除客户时长
func (s *sClient) DurationDelete(ctx context.Context, req *v1.ClientDurationDeleteReq) (res *v1.ClientDurationDeleteRes, err error) {
	// 如果没有权限，则返回错误
	if !s.hasAdminAccess(ctx) {
		return nil, gerror.New("无访问权限")
	}

	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.DurationDelete(ctx, req)
}

// WxGetRemainingDuration 获取客户端用户剩余时长
func (s *sClient) WxGetRemainingDuration(ctx context.Context, req *v1.WxClientRemainingDurationReq) (res *v1.WxClientRemainingDurationRes, err error) {
	// 实例化API调用对象
	impl := &v1.ControllerImpl{}
	return impl.WxGetRemainingDuration(ctx, req)
}

// hasAdminAccess 检查是否有管理员权限
func (s *sClient) hasAdminAccess(ctx context.Context) bool {
	// 从context中获取用户角色，检查是否为管理员
	role := ctx.Value(auth.CtxKeyRole)
	if role == nil {
		return false
	}
	return role.(string) == "admin"
}
