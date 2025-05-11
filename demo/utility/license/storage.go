// Package license 提供授权验证相关工具
package license

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"

	"demo/internal/model"
)

var (
	// 授权配置文件路径
	configFilePath = "manifest/config/license.json"
	// 配置文件操作锁
	configLock sync.RWMutex
	// 授权状态缓存
	licenseStatus = make(map[string]bool)
	// 授权状态缓存锁
	statusLock sync.RWMutex
	// 授权到期时间缓存
	licenseExpireTime = make(map[string]time.Time)
	// 授权到期时间缓存锁
	expireTimeLock sync.RWMutex
)

// init 初始化函数
func init() {
	// 启动时预加载授权记录
	go preloadLicenseRecords()
}

// preloadLicenseRecords 预加载授权记录
func preloadLicenseRecords() {
	ctx := context.Background()
	glog.Info(ctx, "正在预加载授权记录...")

	// 获取配置文件的绝对路径
	absPath, err := GetAbsConfigPath()
	if err != nil {
		glog.Error(ctx, "获取授权配置文件路径失败:", err)
		return
	}

	// 如果文件不存在，则跳过预加载
	if !gfile.Exists(absPath) {
		glog.Info(ctx, "授权配置文件不存在，跳过预加载")
		return
	}

	// 读取授权配置
	config, err := ReadLicenseConfig(absPath)
	if err != nil {
		glog.Error(ctx, "读取授权配置失败:", err)
		return
	}

	// 如果没有记录，则跳过预加载
	if len(config.Records) == 0 {
		glog.Info(ctx, "授权配置中没有记录，跳过预加载")
		return
	}

	// 获取唯一的记录
	record := config.Records[0]

	// 只处理有效的授权记录
	if record.Valid {
		// 调用外部验证接口重新验证授权状态
		glog.Info(ctx, "正在重新验证授权记录: 域名", record.Domain)
		valid, err := directVerifyLicense(ctx, record.ClientIP, record.Domain, record.LicenseCode)

		if err != nil {
			glog.Error(ctx, "重新验证授权记录失败:", err)
			// 验证失败时，仍然使用本地记录，但设置较短的有效期
			expiryTime := time.Now().Add(1 * time.Hour)
			domain := record.Domain
			SetLicenseStatus(domain, true)
			SetLicenseExpireTime(domain, expiryTime)
			glog.Warning(ctx, "外部验证失败，使用本地记录(短期有效): 域名", domain, "到期时间", expiryTime.Format("2006-01-02 15:04:05"))
		} else if valid {
			// 验证成功，更新授权状态和到期时间
			domain := record.Domain
			expiryTime := time.Now().Add(time.Duration(DefaultCacheTTLHours) * time.Hour)
			SetLicenseStatus(domain, true)
			SetLicenseExpireTime(domain, expiryTime)
			glog.Info(ctx, "外部验证成功，更新授权记录: 域名", domain, "到期时间", expiryTime.Format("2006-01-02 15:04:05"))

			// 更新本地记录
			newRecord := &model.LicenseRecord{
				ClientIP:    record.ClientIP,
				Domain:      record.Domain,
				LicenseCode: record.LicenseCode,
				Valid:       true,
				CreateTime:  time.Now(),
			}

			if err := SaveLicenseRecord(ctx, newRecord); err != nil {
				glog.Error(ctx, "更新授权记录失败:", err)
			}
		} else {
			// 验证失败，更新状态和本地记录
			domain := record.Domain
			glog.Warning(ctx, "外部验证失败，授权已无效: 域名", domain)
			ClearLicenseStatus(domain)

			// 更新本地记录为无效
			newRecord := &model.LicenseRecord{
				ClientIP:    record.ClientIP,
				Domain:      record.Domain,
				LicenseCode: record.LicenseCode,
				Valid:       false,
				CreateTime:  time.Now(),
			}

			if err := SaveLicenseRecord(ctx, newRecord); err != nil {
				glog.Error(ctx, "更新授权记录失败:", err)
			}
		}
	} else {
		glog.Info(ctx, "本地记录无效，跳过验证: 域名", record.Domain)
	}

	glog.Info(ctx, "授权记录预加载完成")
}

// SaveLicenseRecord 保存授权记录到配置文件
func SaveLicenseRecord(ctx context.Context, record *model.LicenseRecord) error {
	configLock.Lock()
	defer configLock.Unlock()

	// 获取配置文件的绝对路径
	absPath, err := GetAbsConfigPath()
	if err != nil {
		return err
	}

	// 读取现有配置
	config, err := ReadLicenseConfig(absPath)
	if err != nil {
		if !gfile.Exists(absPath) {
			// 如果文件不存在，则创建一个新的配置对象
			config = &model.LicenseConfig{
				Records:        make([]model.LicenseRecord, 0),
				LastUpdateTime: time.Now().Format(time.RFC3339),
			}
		} else {
			return err
		}
	}

	// 更新配置 - 不再追加到记录列表，而是替换为只包含最新的一条记录
	config.Records = []model.LicenseRecord{*record}
	config.LastUpdateTime = time.Now().Format(time.RFC3339)

	// 写入配置文件
	return writeLicenseConfig(absPath, config)
}

// GetLicenseStatus 获取域名的授权状态
func GetLicenseStatus(domain string) bool {
	// 首先从缓存中获取
	statusLock.RLock()
	if valid, exists := licenseStatus[domain]; exists {
		statusLock.RUnlock()
		return valid
	}
	statusLock.RUnlock()

	// 如果缓存中没有，直接从历史记录中检查
	ctx := context.Background()

	// 获取配置文件的绝对路径
	absPath, err := GetAbsConfigPath()
	if err != nil {
		glog.Error(ctx, "获取授权配置文件路径失败:", err)
		return false
	}

	// 如果文件不存在，则直接返回false
	if !gfile.Exists(absPath) {
		return false
	}

	// 读取授权配置
	config, err := ReadLicenseConfig(absPath)
	if err != nil {
		glog.Error(ctx, "读取授权配置失败:", err)
		return false
	}

	// 检查是否有记录
	if len(config.Records) == 0 {
		return false
	}

	// 获取唯一的记录
	latestRecord := config.Records[0]

	// 只有记录的域名匹配且记录有效时才继续检查
	if latestRecord.Domain == domain && latestRecord.Valid {
		// 检查记录是否在有效期内
		expiryTime := latestRecord.CreateTime.Add(time.Duration(DefaultCacheTTLHours) * time.Hour)
		isValid := expiryTime.After(time.Now())

		if isValid {
			// 更新内存中的授权状态和到期时间
			SetLicenseStatus(domain, true)
			SetLicenseExpireTime(domain, expiryTime)
			glog.Info(ctx, "从历史记录读取授权状态: 域名", domain, "到期时间", expiryTime.Format("2006-01-02 15:04:05"))
			return true
		} else {
			glog.Warning(ctx, "授权已过期: 域名", domain, "过期时间", expiryTime.Format("2006-01-02 15:04:05"))
		}
	}

	return false
}

// SetLicenseStatus 设置域名的授权状态
func SetLicenseStatus(domain string, valid bool) {
	statusLock.Lock()
	defer statusLock.Unlock()

	licenseStatus[domain] = valid

	// 设置授权到期时间（默认24小时）
	if valid {
		SetLicenseExpireTime(domain, time.Now().Add(DefaultCacheTTLHours*time.Hour))
	} else {
		ClearLicenseExpireTime(domain)
	}
}

// ClearLicenseStatus 清除域名的授权状态
func ClearLicenseStatus(domain string) {
	statusLock.Lock()
	defer statusLock.Unlock()

	delete(licenseStatus, domain)
	ClearLicenseExpireTime(domain)
}

// GetLicenseExpireTime 获取域名授权到期时间
func GetLicenseExpireTime(ctx context.Context, domain string) (time.Time, error) {
	expireTimeLock.RLock()
	defer expireTimeLock.RUnlock()

	if expireTime, exists := licenseExpireTime[domain]; exists {
		return expireTime, nil
	}
	return time.Time{}, errors.New("未找到域名授权到期时间")
}

// SetLicenseExpireTime 设置域名授权到期时间
func SetLicenseExpireTime(domain string, expireTime time.Time) {
	expireTimeLock.Lock()
	defer expireTimeLock.Unlock()

	licenseExpireTime[domain] = expireTime
}

// ClearLicenseExpireTime 清除域名授权到期时间
func ClearLicenseExpireTime(domain string) {
	expireTimeLock.Lock()
	defer expireTimeLock.Unlock()

	delete(licenseExpireTime, domain)
}

// ReadLicenseConfig 读取许可证配置文件
func ReadLicenseConfig(path string) (*model.LicenseConfig, error) {
	// 如果文件不存在，则返回错误
	if !gfile.Exists(path) {
		return nil, errors.New("授权配置文件不存在")
	}

	// 读取文件内容
	content := gfile.GetBytes(path)
	if len(content) == 0 {
		glog.Error(context.Background(), "读取授权配置文件失败或文件为空")
		return nil, errors.New("授权配置文件为空")
	}

	// 解析JSON内容
	config := &model.LicenseConfig{}
	if err := gjson.DecodeTo(content, config); err != nil {
		glog.Error(context.Background(), "解析授权配置文件失败:", err)
		return nil, err
	}

	return config, nil
}

// 写入许可证配置文件
func writeLicenseConfig(path string, config *model.LicenseConfig) error {
	// 转换为JSON
	content, err := gjson.Marshal(config)
	if err != nil {
		glog.Error(context.Background(), "序列化授权配置失败:", err)
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := gfile.Mkdir(dir); err != nil {
		glog.Error(context.Background(), "创建授权配置目录失败:", err)
		return err
	}

	// 写入文件
	if err := gfile.PutBytes(path, content); err != nil {
		glog.Error(context.Background(), "写入授权配置文件失败:", err)
		return err
	}

	return nil
}

// GetAbsConfigPath 获取配置文件的绝对路径
func GetAbsConfigPath() (string, error) {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		glog.Error(context.Background(), "获取工作目录失败:", err)
		return "", err
	}

	// 构建绝对路径
	absPath := filepath.Join(workDir, configFilePath)
	return absPath, nil
}

// directVerifyLicense 直接调用外部验证接口验证授权，不读取本地配置，避免递归调用
func directVerifyLicense(ctx context.Context, clientIP, domain, licenseCode string) (bool, error) {
	// 创建验证客户端
	client := g.Client()
	apiURL := DefaultLicenseVerifyURL

	// 构建请求参数
	params := g.Map{
		"client_ip":    clientIP,
		"domain":       domain,
		"license_code": licenseCode,
	}

	// 调用外部验证接口
	response, err := client.Post(ctx, apiURL, params)
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
	return valid, nil
}
