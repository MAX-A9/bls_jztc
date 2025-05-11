package user

import (
	"context"
	v1 "demo/api/user/v1"
	"demo/internal/service"
)

// 获取验证码
func (c *controllerV1) GetCaptcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error) {
	return service.User().GetCaptcha(ctx, req)
}

// 用户登录
func (c *controllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	return service.User().Login(ctx, req)
}

// 用户退出登录
func (c *controllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	return service.User().Logout(ctx, req)
}

// 获取用户个人信息
func (c *controllerV1) GetUserInfo(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
	return service.User().GetUserInfo(ctx, req)
}
