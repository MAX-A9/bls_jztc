package auth

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/golang-jwt/jwt/v5"
)

// 配置结构
type JWTConfig struct {
	// JWT密钥
	Secret string
	// 令牌过期时间(秒)
	Expire int
	// 刷新令牌过期时间(秒)
	RefreshExpire int
	// 签名算法
	Algorithm string
	// 发行人
	Issuer string
}

// 令牌类型
const (
	TokenTypeAdmin  = "admin"  // 管理员令牌
	TokenTypeClient = "client" // 客户令牌
)

// 定义黑名单缓存，用于存储已注销的令牌
var tokenBlacklist = gcache.New()

// 自定义声明结构
type CustomClaims struct {
	// 用户或客户ID
	ID int `json:"id,omitempty"`
	// 用户名
	Username string `json:"username,omitempty"`
	// 令牌类型
	TokenType string `json:"token_type,omitempty"`
	// 角色
	Role string `json:"role,omitempty"`
	// 其他数据
	Data map[string]interface{} `json:"data,omitempty"`
	// JWT标准字段
	jwt.RegisteredClaims
}

// GetConfig 获取JWT配置
func GetConfig(ctx context.Context) *JWTConfig {
	config := &JWTConfig{
		Secret:        "jztc_secret_key", // 默认密钥，与原系统保持一致
		Expire:        604800,            // 7天
		RefreshExpire: 2592000,           // 30天
		Algorithm:     "HS256",
		Issuer:        "demo-server",
	}

	// 从配置文件读取配置
	if g.Cfg().Available(ctx) {
		cfgSecret := g.Cfg().MustGet(ctx, "jwt.secret", "")
		if cfgSecret.String() != "" {
			config.Secret = cfgSecret.String()
		}

		// 其他配置项读取
		config.Expire = g.Cfg().MustGet(ctx, "jwt.expire", config.Expire).Int()
		config.RefreshExpire = g.Cfg().MustGet(ctx, "jwt.refreshExpire", config.RefreshExpire).Int()
		config.Algorithm = g.Cfg().MustGet(ctx, "jwt.algorithm", config.Algorithm).String()
		config.Issuer = g.Cfg().MustGet(ctx, "jwt.issuer", config.Issuer).String()
	}

	return config
}

// 获取签名方法
func getSigningMethod(algorithm string) jwt.SigningMethod {
	switch algorithm {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

// CreateToken 创建令牌
// id: 用户或客户ID
// username: 用户名
// tokenType: 令牌类型(admin/client)
// role: 用户角色(可选)
// data: 其他数据(可选)
func CreateToken(ctx context.Context, id int, username string, tokenType string, role string, data map[string]interface{}) (token string, expire int, err error) {
	// 获取配置
	config := GetConfig(ctx)
	expire = config.Expire

	// 创建自定义声明
	now := time.Now()
	claims := &CustomClaims{
		ID:        id,
		Username:  username,
		TokenType: tokenType,
		Role:      role,
		Data:      data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.Issuer,
		},
	}

	// 创建令牌
	jwtToken := jwt.NewWithClaims(getSigningMethod(config.Algorithm), claims)

	// 签名令牌
	token, err = jwtToken.SignedString([]byte(config.Secret))
	if err != nil {
		return "", 0, gerror.New("生成令牌失败: " + err.Error())
	}

	return token, expire, nil
}

// CreateAdminToken 创建管理员令牌
func CreateAdminToken(ctx context.Context, userId int, username string) (token string, expire int, err error) {
	return CreateToken(ctx, userId, username, TokenTypeAdmin, "admin", nil)
}

// CreateClientToken 创建客户令牌
func CreateClientToken(ctx context.Context, clientId int, username string) (token string, expire int, err error) {
	return CreateToken(ctx, clientId, username, TokenTypeClient, "", nil)
}

// ParseToken 解析令牌
func ParseToken(ctx context.Context, tokenString string) (*CustomClaims, error) {
	if tokenString == "" {
		return nil, gerror.New("令牌不能为空")
	}

	// 检查令牌是否在黑名单中
	isBlacklisted, err := tokenBlacklist.Get(ctx, tokenString)
	if err == nil && isBlacklisted != nil && isBlacklisted.Bool() {
		return nil, gerror.New("令牌已被注销")
	}

	// 获取配置
	config := GetConfig(ctx)

	// 调试日志
	g.Log().Debug(ctx, "解析JWT令牌", "密钥", config.Secret)

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})

	if err != nil {
		g.Log().Error(ctx, "令牌解析失败", "错误", err.Error())
		return nil, gerror.New("令牌无效: " + err.Error())
	}

	if !token.Valid {
		return nil, gerror.New("令牌已失效")
	}

	// 获取声明
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, gerror.New("令牌结构错误")
	}

	return claims, nil
}

// VerifyAdminToken 验证管理员令牌
func VerifyAdminToken(ctx context.Context, tokenString string) (userId int, username string, role string, err error) {
	// 解析令牌
	claims, err := ParseToken(ctx, tokenString)
	if err != nil {
		return 0, "", "", err
	}

	// 验证令牌类型
	if claims.TokenType != TokenTypeAdmin {
		return 0, "", "", gerror.New("非管理员令牌")
	}

	// 验证ID
	if claims.ID <= 0 {
		return 0, "", "", gerror.New("令牌中无用户ID")
	}

	// 返回数据
	userId = claims.ID
	username = claims.Username
	role = claims.Role
	if role == "" {
		role = "admin"
	}

	g.Log().Debug(ctx, "管理员令牌验证成功", "userId", userId, "username", username, "role", role)
	return userId, username, role, nil
}

// VerifyClientToken 验证客户令牌
func VerifyClientToken(ctx context.Context, tokenString string) (clientId int, err error) {
	// 解析令牌
	claims, err := ParseToken(ctx, tokenString)
	if err != nil {
		return 0, err
	}

	// 验证令牌类型
	if claims.TokenType != TokenTypeClient {
		return 0, gerror.New("非客户令牌")
	}

	// 验证ID
	if claims.ID <= 0 {
		return 0, gerror.New("令牌中无客户ID")
	}

	// 返回客户ID
	clientId = claims.ID
	g.Log().Debug(ctx, "客户令牌验证成功", "clientId", clientId)
	return clientId, nil
}

// RevokeToken 将令牌加入黑名单
func RevokeToken(ctx context.Context, tokenString string) error {
	if tokenString == "" {
		return gerror.New("令牌不能为空")
	}

	// 解析令牌以获取过期时间
	claims, err := ParseToken(ctx, tokenString)
	if err != nil {
		return err
	}

	// 计算剩余有效时间
	expirationTime := claims.ExpiresAt.Time
	now := time.Now()

	// 如果令牌已经过期，无需加入黑名单
	if expirationTime.Before(now) {
		return nil
	}

	// 计算剩余有效期（秒）
	remainingTime := int64(expirationTime.Sub(now).Seconds())
	if remainingTime <= 0 {
		remainingTime = 1 // 至少1秒
	}

	// 将令牌加入黑名单，过期时间设置为令牌的原始过期时间
	err = tokenBlacklist.Set(ctx, tokenString, true, time.Duration(remainingTime)*time.Second)
	if err != nil {
		return gerror.New("将令牌加入黑名单失败: " + err.Error())
	}

	g.Log().Debug(ctx, "令牌已加入黑名单", "令牌", tokenString, "过期时间", remainingTime)
	return nil
}
