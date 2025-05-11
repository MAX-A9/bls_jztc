package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// StorageConfigReq 获取存储配置请求
type StorageConfigReq struct {
	g.Meta `path:"/config" method:"get" tags:"存储设置" summary:"获取存储配置" security:"Bearer" description:"获取存储配置，需要管理员权限"`
}

// StorageConfigRes 获取存储配置响应
type StorageConfigRes struct {
	g.Meta `mime:"application/json" example:"json"`
	// OSS配置
	AccessKeyId     string `json:"accessKeyId" dc:"AccessKey ID"`    // AccessKey ID
	AccessKeySecret string `json:"accessKeySecret" dc:"AccessKey密钥"` // AccessKey Secret
	Endpoint        string `json:"endpoint" dc:"OSS端点"`              // OSS端点
	Bucket          string `json:"bucket" dc:"Bucket名称"`             // Bucket名称
	Region          string `json:"region" dc:"区域"`                   // 区域
	Directory       string `json:"directory" dc:"存储目录"`              // 存储目录
	PublicAccess    bool   `json:"publicAccess" dc:"是否公开访问"`         // 是否公开访问
}

// SaveStorageConfigReq 保存存储配置请求
type SaveStorageConfigReq struct {
	g.Meta `path:"/config" method:"post" tags:"存储设置" summary:"保存存储配置" security:"Bearer" description:"保存存储配置，需要管理员权限"`
	// OSS配置
	AccessKeyId     string `v:"required#请输入AccessKey ID" json:"accessKeyId" dc:"AccessKey ID"`        // AccessKey ID
	AccessKeySecret string `v:"required#请输入AccessKey Secret" json:"accessKeySecret" dc:"AccessKey密钥"` // AccessKey Secret
	Endpoint        string `v:"required#请输入OSS端点" json:"endpoint" dc:"OSS端点"`                         // OSS端点
	Bucket          string `v:"required#请输入Bucket名称" json:"bucket" dc:"Bucket名称"`                     // Bucket名称
	Region          string `v:"required#请选择区域" json:"region" dc:"区域"`                                 // 区域
	Directory       string `json:"directory" dc:"存储目录"`                                               // 存储目录
	PublicAccess    bool   `json:"publicAccess" dc:"是否公开访问"`                                          // 是否公开访问
}

// SaveStorageConfigRes 保存存储配置响应
type SaveStorageConfigRes struct {
	g.Meta `mime:"application/json" example:"json"`
	// 操作结果
	Success bool   `json:"success" dc:"是否成功"` // 是否成功
	Message string `json:"message" dc:"消息"`   // 消息
}
