// Package service 提供服务接口定义
package service

import (
	"context"
)

// ILicense 授权验证服务接口
type ILicense interface {
	// VerifyLicense 调用外部授权验证接口
	VerifyLicense(ctx context.Context, clientIP, domain, licenseCode string) (bool, error)
}
