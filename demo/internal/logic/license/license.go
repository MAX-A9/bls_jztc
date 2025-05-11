// Package license 提供授权验证相关的业务逻辑
package license

import (
	"context"
	"demo/utility/license"
)

// LicenseService 授权验证服务
type LicenseService struct {
	client *license.LicenseClient
}

// New 创建授权验证服务
func New() *LicenseService {
	// 创建授权验证客户端，使用默认地址
	return &LicenseService{
		client: license.NewClient("", 0), // 使用默认参数
	}
}

// VerifyLicense 验证授权
func (s *LicenseService) VerifyLicense(ctx context.Context, clientIP, domain, licenseCode string) (bool, error) {
	// 调用授权验证客户端
	return s.client.VerifyLicense(ctx, clientIP, domain, licenseCode)
}
