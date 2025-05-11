// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"demo/api/user/v1"
)

type IUserV1 interface {
	// 用户列表
	List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error)
	// 创建用户
	Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error)
	// 更新用户
	Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error)
	// 删除用户
	Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error)
	// 获取验证码
	GetCaptcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error)
	// 用户登录
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	// 退出登录
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
	// 获取用户个人信息
	GetUserInfo(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error)
} 