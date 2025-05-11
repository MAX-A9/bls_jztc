// Package model 提供数据模型定义
package model

import (
	"time"
)

// LicenseRecord 授权记录模型
type LicenseRecord struct {
	ClientIP    string    `json:"client_ip"`    // 客户端IP
	Domain      string    `json:"domain"`       // 域名
	LicenseCode string    `json:"license_code"` // 授权码
	Valid       bool      `json:"valid"`        // 是否有效
	CreateTime  time.Time `json:"create_time"`  // 创建时间
}

// LicenseConfig 授权配置模型
type LicenseConfig struct {
	Records        []LicenseRecord `json:"records"`        // 授权记录列表
	LastUpdateTime string          `json:"lastUpdateTime"` // 最后更新时间
}
