// Package license 提供授权验证控制器
package license

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"

	"demo/api/license"
	licenseLogic "demo/internal/logic/license"
	"demo/internal/model"
	licenseUtil "demo/utility/license"
)

// Controller 授权验证控制器
type Controller struct{}

// NewController 创建授权验证控制器
func NewController() *Controller {
	return &Controller{}
}

// Register 注册路由
func (c *Controller) Register(group *ghttp.RouterGroup) {
	group.Bind(
		c.Verify,
		c.Status,
	)
}

// Verify 验证授权
func (c *Controller) Verify(ctx context.Context, req *license.VerifyReq) (res *license.VerifyRes, err error) {
	// 如果客户端IP为空，自动获取服务器IP
	if gstr.Trim(req.ClientIP) == "" {
		req.ClientIP = licenseUtil.GetServerIP()
		glog.Debug(ctx, "使用服务器IP作为客户端IP:", req.ClientIP)
	}

	// 记录请求日志
	glog.Debug(ctx, "收到授权验证请求，客户端IP:", req.ClientIP, "域名:", req.Domain)

	// 初始化响应
	res = &license.VerifyRes{
		Valid: false,
	}

	// 标准化处理域名
	domain := licenseUtil.ValidateDomain(req.Domain)

	// 调用业务逻辑验证授权
	licenseService := licenseLogic.New()
	valid, err := licenseService.VerifyLicense(ctx, req.ClientIP, domain, req.LicenseCode)
	if err != nil {
		glog.Error(ctx, "授权验证过程中发生错误:", err)
		return res, gerror.NewCode(gcode.CodeInternalError, "授权验证服务暂时不可用")
	}

	// 设置验证结果
	res.Valid = valid

	// 记录验证结果日志
	if valid {
		glog.Info(ctx, "授权验证成功，域名:", domain)
	} else {
		glog.Warning(ctx, "授权验证失败，域名:", domain)
	}

	// 创建授权记录对象
	record := &model.LicenseRecord{
		ClientIP:    req.ClientIP,
		Domain:      domain,
		LicenseCode: req.LicenseCode,
		Valid:       valid,
		CreateTime:  time.Now(),
	}

	// 保存授权记录（无论成功或失败）
	if err := licenseUtil.SaveLicenseRecord(ctx, record); err != nil {
		glog.Error(ctx, "保存授权记录失败:", err)
		// 不影响返回结果，只记录日志
	}

	// 根据授权验证结果设置或清除域名授权状态缓存
	if valid {
		licenseUtil.SetLicenseStatus(domain, true)
		glog.Info(ctx, "更新授权状态缓存，域名:", domain)
	} else {
		// 清除可能存在的缓存状态
		licenseUtil.ClearLicenseStatus(domain)
		glog.Warning(ctx, "清除授权状态缓存，域名:", domain)
	}

	return res, nil
}

// Status 获取授权状态
func (c *Controller) Status(ctx context.Context, req *license.StatusReq) (res *license.StatusRes, err error) {
	// 记录请求日志
	glog.Debug(ctx, "收到获取授权状态请求，域名:", req.Domain)

	// 初始化响应
	res = &license.StatusRes{
		Valid: false,
	}

	// 标准化处理域名
	domain := licenseUtil.ValidateDomain(req.Domain)

	// 检查域名是否为本地环境
	if licenseUtil.IsLocalhost(domain) {
		res.Valid = true
		return res, nil
	}

	// 获取域名授权状态
	valid := licenseUtil.GetLicenseStatus(domain)

	// 如果授权无效，尝试直接从历史记录中读取状态
	if !valid {
		glog.Debug(ctx, "尝试从历史记录中读取授权状态，域名:", domain)

		// 获取配置文件的绝对路径
		absPath, err := licenseUtil.GetAbsConfigPath()
		if err == nil && gfile.Exists(absPath) {
			// 读取授权配置
			config, err := licenseUtil.ReadLicenseConfig(absPath)
			if err == nil && config != nil && len(config.Records) > 0 {
				// 获取唯一的记录
				record := config.Records[0]
				// 只有域名匹配且记录有效时才继续检查
				if record.Domain == domain && record.Valid {
					// 检查是否在有效期内
					expiryTime := record.CreateTime.Add(time.Duration(licenseUtil.DefaultCacheTTLHours) * time.Hour)
					isValid := expiryTime.After(time.Now())

					if isValid {
						// 更新内存中的授权状态和到期时间
						licenseUtil.SetLicenseStatus(domain, true)
						licenseUtil.SetLicenseExpireTime(domain, expiryTime)
						valid = true
						glog.Info(ctx, "从历史记录恢复授权状态成功，域名:", domain)
					}
				}
			}
		}
	}

	res.Valid = valid

	return res, nil
}
