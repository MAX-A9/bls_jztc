package model

// StorageConfig 存储配置模型
type StorageConfig struct {
	// OSS配置
	AccessKeyId     string `json:"accessKeyId"`     // AccessKey ID
	AccessKeySecret string `json:"accessKeySecret"` // AccessKey Secret
	Endpoint        string `json:"endpoint"`        // OSS端点
	Bucket          string `json:"bucket"`          // Bucket名称
	Region          string `json:"region"`          // 区域
	Directory       string `json:"directory"`       // 存储目录
	PublicAccess    bool   `json:"publicAccess"`    // 是否公开访问
}
