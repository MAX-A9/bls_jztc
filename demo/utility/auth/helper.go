package auth

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ExtractTokenFromRequest 从请求中提取令牌
func ExtractTokenFromRequest(r *ghttp.Request) (string, error) {
	// 获取请求头中的Authorization Token
	authHeader := r.GetHeader("Authorization")
	if authHeader == "" {
		return "", gerror.New("未提供授权令牌")
	}

	// 处理Bearer Token格式
	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	return tokenString, nil
}

// ExtractTokenFromHeader 从HTTP头中提取令牌
func ExtractTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", gerror.New("未提供授权令牌")
	}

	// 处理Bearer Token格式
	tokenString := header
	if strings.HasPrefix(header, "Bearer ") {
		tokenString = strings.TrimPrefix(header, "Bearer ")
	}

	return tokenString, nil
}

// GetClientId 从上下文中获取客户ID
func GetClientId(ctx context.Context) (int, error) {
	clientId, err := GetClientInfo(ctx)
	return clientId, err
}

// GetUserId 从上下文中获取用户ID
func GetUserId(ctx context.Context) (int, error) {
	userId, _, _, err := GetUserInfo(ctx)
	return userId, err
}
