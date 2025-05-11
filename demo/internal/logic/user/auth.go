package user

import (
	"context"
	"strings"
	"time"

	v1 "demo/api/user/v1"
	"demo/internal/consts"
	"demo/internal/model/do"
	"demo/utility/auth"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/mojocn/base64Captcha"
)

// 全局缓存对象
var cache = gcache.New()

// 验证码存储器
var captchaStore = base64Captcha.DefaultMemStore

// createCaptcha 创建验证码
func createCaptcha() (id, base64Str string, err error) {
	// 从配置文件获取验证码配置
	width := g.Cfg().MustGet(context.Background(), "captcha.width", 120).Int()
	height := g.Cfg().MustGet(context.Background(), "captcha.height", 40).Int()
	length := g.Cfg().MustGet(context.Background(), "captcha.length", 4).Int()

	// 配置验证码驱动
	driver := base64Captcha.NewDriverDigit(
		height, // 高度
		width,  // 宽度
		length, // 长度
		0.7,    // 验证码数字的最大倾斜角度
		80,     // 干扰点数量
	)

	// 创建验证码并获取相关信息
	captcha := base64Captcha.NewCaptcha(driver, captchaStore)
	// 生成验证码图片
	var b64s string
	id, b64s, _, err = captcha.Generate()
	if err != nil {
		return "", "", err
	}

	return id, b64s, nil
}

// GetCaptcha 获取验证码
func (s *userImpl) GetCaptcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error) {
	res = &v1.CaptchaRes{}

	// 生成验证码
	id, base64Str, err := createCaptcha()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeCaptchaGenerateFailed, err.Error())
	}

	// 设置响应
	res.Id = id
	res.Base64 = base64Str

	// 获取验证码过期时间（秒）
	expire := g.Cfg().MustGet(ctx, "captcha.expire", 300).Int()
	res.ExpiredAt = gtime.Now().Add(time.Duration(expire) * time.Second).String()

	return res, nil
}

// VerifyCaptcha 验证验证码
func (s *userImpl) VerifyCaptcha(ctx context.Context, id string, code string) (match bool, err error) {
	if id == "" || code == "" {
		return false, gerror.NewCode(consts.CodeCaptchaInvalid)
	}

	// 验证码不区分大小写
	match = captchaStore.Verify(id, code, true)
	if !match {
		return false, gerror.NewCode(consts.CodeCaptchaVerifyFailed)
	}

	return true, nil
}

// Login 用户登录
func (s *userImpl) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	res = &v1.LoginRes{}

	// 验证验证码
	match, err := s.VerifyCaptcha(ctx, req.CaptchaId, req.CaptchaCode)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, gerror.NewCode(consts.CodeCaptchaVerifyFailed)
	}

	// 获取用户
	user, err := s.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeLoginFailed, err.Error())
	}
	if user == nil {
		return nil, gerror.NewCode(consts.CodeUserNotExists)
	}

	// 验证用户状态
	if user.Status == 0 {
		return nil, gerror.NewCode(consts.CodeUserForbidden)
	}

	// 验证密码
	encryptedPassword, err := gmd5.EncryptString(req.Password)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeLoginFailed, "密码加密失败")
	}
	if user.Password != encryptedPassword {
		return nil, gerror.NewCode(consts.CodeUserPasswordError)
	}

	// 生成token (使用utility/jwt包中的CreateToken)
	token, expireIn, err := auth.CreateAdminToken(ctx, user.Id, user.Username)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeTokenGenerateFailed, err.Error())
	}

	// 更新最后登录信息
	ip := g.RequestFromCtx(ctx).GetClientIp()
	_, err = s.userDao.Update(ctx, &do.UserDO{
		LastLoginIp:   ip,
		LastLoginTime: gtime.Now(),
	}, user.Id)
	if err != nil {
		g.Log().Error(ctx, "更新用户最后登录信息失败", err)
	}

	// 设置响应
	res.Token = token
	res.UserId = user.Id
	res.Nickname = user.Nickname
	res.Role = user.Role
	res.ExpireIn = expireIn

	return res, nil
}

// GenerateToken 生成JWT令牌
func (s *userImpl) GenerateToken(ctx context.Context, userId int, username string) (token string, expireIn int, err error) {
	// 调用auth包的CreateAdminToken方法
	return auth.CreateAdminToken(ctx, userId, username)
}

// Logout 用户退出登录
func (s *userImpl) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	res = &v1.LogoutRes{
		Success: false,
	}

	// 获取请求中的Authorization头
	authHeader := g.RequestFromCtx(ctx).GetHeader("Authorization")
	if authHeader == "" {
		return res, gerror.NewCode(consts.CodeUnauthorized, "未提供认证信息")
	}

	// 处理Bearer Token格式
	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 将令牌加入黑名单
	err = auth.RevokeToken(ctx, tokenString)
	if err != nil {
		g.Log().Error(ctx, "退出登录失败", err)
		return res, gerror.NewCode(consts.CodeLogoutFailed, err.Error())
	}

	res.Success = true
	return res, nil
}

// GetUserInfo 获取用户个人信息
func (s *userImpl) GetUserInfo(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
	res = &v1.UserInfoRes{}

	// 从上下文中获取用户信息
	userId, _, _, err := auth.GetUserInfo(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "未登录或登录已过期")
	}

	// 获取用户信息
	user, err := s.GetUserById(ctx, userId)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeUserNotExists, err.Error())
	}
	if user == nil {
		return nil, gerror.NewCode(consts.CodeUserNotExists)
	}

	// 转换为API响应格式
	res.Id = user.Id
	res.Username = user.Username
	res.Nickname = user.Nickname
	res.Role = user.Role
	res.Status = user.Status

	if user.LastLoginTime != nil {
		res.LastLogin = user.LastLoginTime.String()
	}

	return res, nil
}
