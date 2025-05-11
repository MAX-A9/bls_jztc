package service

import (
	"context"
	v1 "demo/api/user/v1"
	"demo/internal/model/entity"
)

// UserService 用户服务接口
type UserService interface {
	// List 获取用户列表
	List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error)

	// Create 创建用户
	Create(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error)

	// Update 更新用户
	Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error)

	// Delete 删除用户
	Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error)

	// GetUserById 通过ID获取用户
	GetUserById(ctx context.Context, id int) (user *entity.User, err error)

	// GetUserByUsername 通过用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (user *entity.User, err error)

	// GetCaptcha 获取验证码
	GetCaptcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error)

	// VerifyCaptcha 验证验证码
	VerifyCaptcha(ctx context.Context, id string, code string) (match bool, err error)

	// Login 用户登录
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)

	// Logout 用户退出登录
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)

	// GetUserInfo 获取用户个人信息
	GetUserInfo(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error)

	// GenerateToken 生成JWT令牌
	GenerateToken(ctx context.Context, userId int, username string) (token string, expireIn int, err error)
}
