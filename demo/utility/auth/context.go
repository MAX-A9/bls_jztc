package auth

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

// 上下文键类型
type ContextKey string

// 上下文键常量
const (
	// 用户ID键
	CtxKeyUserId ContextKey = "userId"
	// 用户名键
	CtxKeyUsername ContextKey = "username"
	// 角色键
	CtxKeyRole ContextKey = "role"
	// 客户ID键
	CtxKeyClientId ContextKey = "clientId"
	// 令牌类型键
	CtxKeyTokenType ContextKey = "tokenType"
	// 请求来源键
	CtxKeyRequestFrom ContextKey = "requestFrom"
	// 微信OpenID键
	CtxKeyOpenID ContextKey = "openId"
)

// SetUserInfo 设置用户信息到上下文
func SetUserInfo(ctx context.Context, userId int, username string, role string) context.Context {
	ctx = context.WithValue(ctx, CtxKeyUserId, userId)
	ctx = context.WithValue(ctx, CtxKeyUsername, username)
	ctx = context.WithValue(ctx, CtxKeyRole, role)
	ctx = context.WithValue(ctx, CtxKeyTokenType, TokenTypeAdmin)
	return ctx
}

// SetClientInfo 设置客户信息到上下文
func SetClientInfo(ctx context.Context, clientId int) context.Context {
	ctx = context.WithValue(ctx, CtxKeyClientId, clientId)
	ctx = context.WithValue(ctx, CtxKeyTokenType, TokenTypeClient)
	return ctx
}

// GetUserInfo 从上下文获取用户信息
func GetUserInfo(ctx context.Context) (userId int, username string, role string, err error) {
	// 获取用户ID
	userIdValue := ctx.Value(CtxKeyUserId)
	if userIdValue == nil {
		return 0, "", "", gerror.New("未登录或登录已过期")
	}

	// 获取令牌类型
	tokenType := ctx.Value(CtxKeyTokenType)
	if tokenType == nil || tokenType.(string) != TokenTypeAdmin {
		return 0, "", "", gerror.New("非管理员令牌")
	}

	userId = userIdValue.(int)

	// 获取用户名
	usernameValue := ctx.Value(CtxKeyUsername)
	if usernameValue != nil {
		username = usernameValue.(string)
	}

	// 获取角色
	roleValue := ctx.Value(CtxKeyRole)
	if roleValue != nil {
		role = roleValue.(string)
	} else {
		role = "admin" // 默认为管理员角色
	}

	return userId, username, role, nil
}

// GetClientInfo 从上下文获取客户信息
func GetClientInfo(ctx context.Context) (clientId int, err error) {
	// 获取客户ID
	clientIdValue := ctx.Value(CtxKeyClientId)
	if clientIdValue == nil {
		return 0, gerror.New("未登录或登录已过期")
	}

	// 获取令牌类型
	tokenType := ctx.Value(CtxKeyTokenType)
	if tokenType == nil || tokenType.(string) != TokenTypeClient {
		return 0, gerror.New("非客户令牌")
	}

	clientId = clientIdValue.(int)
	return clientId, nil
}

// GetCurrentUser 获取当前登录用户（兼容原有方法）
func GetCurrentUser(ctx context.Context) (userId int, username string, role string, err error) {
	return GetUserInfo(ctx)
}

// GetRequestFrom 获取请求来源
func GetRequestFrom(ctx context.Context) string {
	// 从请求头或上下文中获取请求来源
	// 如果是微信小程序，返回"wxapp"，其他来源可以返回相应的标识
	// 这里简单实现，实际应用中需要根据具体情况判断请求来源
	reqFrom := ctx.Value(CtxKeyRequestFrom)
	if reqFrom != nil {
		return reqFrom.(string)
	}

	// 默认返回空字符串
	return ""
}

// GetOpenID 获取微信OpenID
func GetOpenID(ctx context.Context) string {
	// 从上下文中获取OpenID
	openIDValue := ctx.Value(CtxKeyOpenID)
	if openIDValue != nil {
		return openIDValue.(string)
	}

	// 默认返回空字符串
	return ""
}

// SetOpenID 设置微信OpenID到上下文
func SetOpenID(ctx context.Context, openID string) context.Context {
	return context.WithValue(ctx, CtxKeyOpenID, openID)
}

// SetRequestFrom 设置请求来源到上下文
func SetRequestFrom(ctx context.Context, requestFrom string) context.Context {
	return context.WithValue(ctx, CtxKeyRequestFrom, requestFrom)
}
