package auth

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware(r *ghttp.Request) {
	// 获取请求头中的Authorization Token
	authHeader := r.GetHeader("Authorization")
	if authHeader == "" {
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "未授权访问",
		})
		r.Exit()
		return
	}

	// 处理Bearer Token格式
	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 验证管理员令牌
	userId, username, role, err := VerifyAdminToken(r.GetCtx(), tokenString)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": err.Error(),
		})
		r.Exit()
		return
	}

	// 验证用户角色是否为管理员
	if role != "admin" {
		r.Response.WriteJson(g.Map{
			"code":    403,
			"message": "您没有权限访问此资源，需要管理员权限",
		})
		r.Exit()
		return
	}

	// 将用户信息存储到请求上下文中
	ctx := SetUserInfo(r.GetCtx(), userId, username, role)
	r.SetCtx(ctx)

	// 继续处理请求
	r.Middleware.Next()
}

// ClientAuthMiddleware 客户认证中间件
func ClientAuthMiddleware(r *ghttp.Request) {
	// 获取请求头中的Authorization Token
	authHeader := r.GetHeader("Authorization")
	if authHeader == "" {
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "未授权访问",
		})
		r.Exit()
		return
	}

	// 处理Bearer Token格式
	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 验证客户令牌
	clientId, err := VerifyClientToken(r.GetCtx(), tokenString)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": err.Error(),
		})
		r.Exit()
		return
	}

	// 将客户ID存储到上下文中
	ctx := SetClientInfo(r.GetCtx(), clientId)

	// 检查是否来自微信小程序的请求
	userAgent := r.Header.Get("User-Agent")
	isWxApp := r.GetHeader("X-Source") == "wxapp" || strings.Contains(strings.ToLower(userAgent), "miniprogram")

	if isWxApp {
		// 如果是微信小程序请求，设置请求来源为wxapp
		ctx = SetRequestFrom(ctx, "wxapp")

		// 从数据库查询客户的openid
		var client struct {
			OpenId string `json:"open_id"`
		}
		err = g.DB().Model("client").
			Fields("open_id").
			Where("id", clientId).
			Scan(&client)

		if err == nil && client.OpenId != "" {
			// 如果查询成功并且openid不为空，则设置到上下文中
			ctx = SetOpenID(ctx, client.OpenId)
			g.Log().Debug(ctx, "设置微信OpenID到上下文", "clientId", clientId, "openId", client.OpenId)
		} else {
			g.Log().Warning(ctx, "无法获取客户openId", "clientId", clientId, "error", err)
		}
	}

	r.SetCtx(ctx)

	// 设置旧版兼容变量
	r.SetCtxVar("client_id", clientId)

	// 继续处理请求
	r.Middleware.Next()
}
