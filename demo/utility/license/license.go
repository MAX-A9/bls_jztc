// Package license 提供授权验证相关工具
package license

import (
	"context"
	"demo/internal/model"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
)

// 授权验证缓存
var licenseCache = gcache.New()

// LicenseClient 授权验证客户端
type LicenseClient struct {
	// 授权验证接口地址
	apiURL string
	// 缓存有效期
	cacheTTL time.Duration
	// HTTP客户端
	client *gclient.Client
}

// NewClient 创建授权验证客户端
func NewClient(apiURL string, cacheTTLHours int) *LicenseClient {
	if apiURL == "" {
		apiURL = DefaultLicenseVerifyURL
	}
	if cacheTTLHours <= 0 {
		cacheTTLHours = DefaultCacheTTLHours
	}

	return &LicenseClient{
		apiURL:   apiURL,
		cacheTTL: time.Duration(cacheTTLHours) * time.Hour,
		client:   g.Client(),
	}
}

// VerifyLicense 验证授权
func (s *LicenseClient) VerifyLicense(ctx context.Context, clientIP, domain, licenseCode string) (bool, error) {
	// 构建缓存键
	cacheKey := "license:" + domain

	// 尝试从缓存中获取
	if v, err := licenseCache.Get(ctx, cacheKey); err == nil && !v.IsNil() {
		return v.Bool(), nil
	}

	// 缓存过期，尝试从license.json文件读取历史记录
	var lastRecord *model.LicenseRecord = nil
	absPath, err := GetAbsConfigPath()
	if err == nil {
		config, err := ReadLicenseConfig(absPath)
		if err == nil && config != nil && len(config.Records) > 0 {
			// 获取唯一的记录
			record := config.Records[0]
			// 只有域名匹配且记录有效时才使用
			if record.Domain == domain && record.Valid {
				lastRecord = &record
			}
		}
	}

	// 如果找到历史记录，使用历史记录中的参数请求外部验证接口
	if lastRecord != nil {
		glog.Info(ctx, "使用历史授权记录参数验证，域名:", domain)
		// 使用历史记录中的参数
		clientIP = lastRecord.ClientIP
		licenseCode = lastRecord.LicenseCode
	}

	// 构建请求参数
	params := g.Map{
		"client_ip":    clientIP,
		"domain":       domain,
		"license_code": licenseCode,
	}

	// 调用外部验证接口
	response, err := s.client.Post(ctx, s.apiURL, params)
	if err != nil {
		glog.Error(ctx, "授权验证接口请求失败:", err)
		return false, err
	}
	defer response.Close()

	// 解析响应结果
	respContent := response.ReadAllString()

	respData, err := gjson.DecodeToJson(respContent)
	if err != nil {
		glog.Error(ctx, "解析授权验证接口响应失败:", err)
		return false, err
	}

	// 检查响应状态码
	code := respData.Get("code").Int()
	if code != 0 {
		errMsg := respData.Get("message").String()
		glog.Warning(ctx, "授权验证失败:", errMsg)
		return false, nil
	}

	// 获取授权结果
	valid := respData.Get("data.valid").Bool()

	// 缓存验证结果，无论成功还是失败
	if valid {
		glog.Info(ctx, "授权验证成功，域名:", domain, "将结果缓存", s.cacheTTL.Hours(), "小时")
		_ = licenseCache.Set(ctx, cacheKey, true, s.cacheTTL)
	} else {
		glog.Warning(ctx, "授权验证失败，域名:", domain, "记录已保存至配置文件")
	}

	return valid, nil
}
