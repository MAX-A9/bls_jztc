package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CaptchaReq 验证码请求
type CaptchaReq struct {
	g.Meta `path:"/captcha" method:"get" summary:"获取验证码" tags:"用户" description:"获取登录验证码，返回验证码ID和Base64编码的图片"`
}

// CaptchaRes 验证码响应
type CaptchaRes struct {
	// 验证码ID
	Id string `json:"id" dc:"验证码唯一标识，用于后续登录验证"`

	// 验证码Base64字符串
	Base64 string `json:"base64" dc:"验证码图片的Base64编码字符串"`

	// 过期时间
	ExpiredAt string `json:"expiredAt" dc:"验证码过期时间"`
}

// LoginReq 登录请求
type LoginReq struct {
	g.Meta `path:"/login" method:"post" summary:"用户登录" tags:"用户" description:"用户登录接口，登录成功后返回JWT令牌"`

	// 用户名
	Username string `v:"required|length:4,30#用户名不能为空|用户名长度应当在:min到:max之间" json:"username" dc:"用户名，长度4-30位"`

	// 密码
	Password string `v:"required|length:6,30#密码不能为空|密码长度应当在:min到:max之间" json:"password" dc:"密码，长度6-30位"`

	// 验证码ID
	CaptchaId string `v:"required#验证码ID不能为空" json:"captchaId" dc:"验证码ID，通过获取验证码接口返回"`

	// 验证码
	CaptchaCode string `v:"required#验证码不能为空" json:"captchaCode" dc:"验证码内容"`
}

// LoginRes 登录响应
type LoginRes struct {
	// 用户ID
	UserId int `json:"userId" dc:"用户ID"`

	// 用户昵称
	Nickname string `json:"nickname" dc:"用户昵称"`

	// 用户角色
	Role string `json:"role" dc:"用户角色 admin:管理员 user:普通用户"`

	// 令牌
	Token string `json:"token" dc:"JWT令牌，后续请求需要在Authorization头中携带此令牌"`

	// 过期时间(秒)
	ExpireIn int `json:"expireIn" dc:"令牌过期时间(秒)"`
}

// LogoutReq 退出登录请求
type LogoutReq struct {
	g.Meta `path:"/logout" method:"post" summary:"退出登录" tags:"用户" security:"Bearer" description:"用户退出登录接口，使当前令牌失效"`
}

// LogoutRes 退出登录响应
type LogoutRes struct {
	g.Meta  `mime:"application/json" example:"json"`
	Success bool `json:"success" dc:"退出登录是否成功"`
}
