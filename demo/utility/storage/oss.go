package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"demo/internal/model"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gogf/gf/v2/util/guid"
)

const (
	// 配置文件路径
	configFilePath = "manifest/config/storage.json"
	// 客户端过期时间(分钟)
	clientExpireMinutes = 30
)

var (
	// 客户端实例
	ossClientInstance *OSSClient
	// 同步锁
	ossClientMutex sync.RWMutex
	// 客户端创建时间
	clientCreatedAt time.Time
)

// OSSBucket 定义OSS桶接口
type OSSBucket interface {
	PutObject(objectKey string, reader interface{}) error
	SignURL(objectKey string, method string, expireSeconds int64) (string, error)
	DeleteObject(objectKey string) error
}

// OSSSimpleBucket 简单的OSS桶实现
type OSSSimpleBucket struct{}

// PutObject 上传对象
func (b *OSSSimpleBucket) PutObject(objectKey string, reader interface{}) error {
	// 临时实现，返回错误让用户知道需要安装OSS SDK
	return errors.New("请先安装阿里云OSS SDK: go get github.com/aliyun/aliyun-oss-go-sdk/oss")
}

// SignURL 生成签名URL
func (b *OSSSimpleBucket) SignURL(objectKey string, method string, expireSeconds int64) (string, error) {
	// 临时实现，返回错误让用户知道需要安装OSS SDK
	return "", errors.New("请先安装阿里云OSS SDK: go get github.com/aliyun/aliyun-oss-go-sdk/oss")
}

// DeleteObject 删除对象
func (b *OSSSimpleBucket) DeleteObject(objectKey string) error {
	// 临时实现，返回错误让用户知道需要安装OSS SDK
	return errors.New("请先安装阿里云OSS SDK: go get github.com/aliyun/aliyun-oss-go-sdk/oss")
}

// OSSClient OSS客户端封装
type OSSClient struct {
	Client     *oss.Client          // OSS客户端
	Bucket     *oss.Bucket          // Bucket对象
	Config     *model.StorageConfig // 配置信息
	PublicURL  string               // 公共访问URL前缀
	PrivateURL string               // 私有访问URL前缀
}

// GetOSSClient 获取OSS客户端实例(单例模式)
func GetOSSClient(ctx context.Context) (*OSSClient, error) {
	ossClientMutex.RLock()
	// 如果客户端已存在且未过期，直接返回
	if ossClientInstance != nil && time.Since(clientCreatedAt).Minutes() < clientExpireMinutes {
		defer ossClientMutex.RUnlock()
		return ossClientInstance, nil
	}
	ossClientMutex.RUnlock()

	// 获取写锁，重新创建客户端
	ossClientMutex.Lock()
	defer ossClientMutex.Unlock()

	// 双重检查，防止在获取锁期间其他协程已创建
	if ossClientInstance != nil && time.Since(clientCreatedAt).Minutes() < clientExpireMinutes {
		return ossClientInstance, nil
	}

	// 读取配置
	config, err := readStorageConfig(ctx)
	if err != nil {
		return nil, errors.New("读取存储配置失败: " + err.Error())
	}

	// 验证配置完整性
	if config.AccessKeyId == "" || config.AccessKeySecret == "" || config.Endpoint == "" || config.Bucket == "" {
		return nil, errors.New("存储配置不完整，请先配置")
	}

	// 创建OSS客户端
	client, err := oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		glog.Error(ctx, "创建OSS客户端失败", err)
		return nil, errors.New("创建OSS客户端失败: " + err.Error())
	}

	// 获取Bucket对象
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		glog.Error(ctx, "获取OSS Bucket失败", err)
		return nil, errors.New("获取OSS Bucket失败: " + err.Error())
	}

	// 构建URL前缀
	publicURL := ""
	privateURL := ""

	// 构建公共访问URL
	if config.Endpoint != "" && config.Bucket != "" {
		if config.PublicAccess {
			publicURL = "https://" + config.Bucket + "." + config.Endpoint
			privateURL = publicURL
		} else {
			privateURL = "https://" + config.Bucket + "." + config.Endpoint
		}
	}

	// 创建客户端实例
	ossClientInstance = &OSSClient{
		Client:     client,
		Bucket:     bucket,
		Config:     config,
		PublicURL:  publicURL,
		PrivateURL: privateURL,
	}

	// 更新创建时间
	clientCreatedAt = time.Now()

	glog.Info(ctx, "OSS客户端初始化成功")
	return ossClientInstance, nil
}

// MustGetOSSClient 获取OSS客户端，如果失败则panic
func MustGetOSSClient(ctx context.Context) *OSSClient {
	client, err := GetOSSClient(ctx)
	if err != nil {
		panic("获取OSS客户端失败: " + err.Error())
	}
	return client
}

// ResetOSSClient 重置OSS客户端
func ResetOSSClient() {
	ossClientMutex.Lock()
	defer ossClientMutex.Unlock()
	ossClientInstance = nil
}

// readStorageConfig 读取存储配置
func readStorageConfig(ctx context.Context) (*model.StorageConfig, error) {
	// 检查配置文件是否存在
	if !gfile.Exists(configFilePath) {
		return nil, errors.New("存储配置文件不存在")
	}

	// 读取配置文件内容
	content := gfile.GetContents(configFilePath)
	if content == "" {
		glog.Error(ctx, "读取配置文件失败")
		return nil, errors.New("读取配置文件失败")
	}

	// 解析JSON内容
	config := &model.StorageConfig{}
	if err := gjson.DecodeTo(content, config); err != nil {
		glog.Error(ctx, "解析配置文件失败", err)
		return nil, errors.New("解析配置文件失败")
	}

	return config, nil
}

// GetObjectURL 获取对象的URL
func (c *OSSClient) GetObjectURL(objectKey string, private bool, expireSeconds int64) (string, error) {
	// 判断是否为公开访问
	if !private && c.Config.PublicAccess {
		// 公开访问
		if c.PublicURL != "" {
			return c.PublicURL + "/" + objectKey, nil
		}
		return "", errors.New("未配置公共访问URL")
	}

	// 私有访问，生成签名URL
	if expireSeconds <= 0 {
		expireSeconds = 3600 // 默认1小时
	}

	// 生成签名URL
	signedURL, err := c.Bucket.SignURL(objectKey, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", errors.New("生成签名URL失败: " + err.Error())
	}

	return signedURL, nil
}

// GetFullObjectKey 获取完整的对象键
func (c *OSSClient) GetFullObjectKey(objectKey string) string {
	// 拼接目录和对象键
	if c.Config.Directory != "" {
		// 确保目录以/结尾
		directory := c.Config.Directory
		if directory[len(directory)-1:] != "/" {
			directory += "/"
		}
		// 确保对象键不以/开头
		if objectKey != "" && objectKey[0:1] == "/" {
			objectKey = objectKey[1:]
		}
		return directory + objectKey
	}
	// 无目录，直接返回对象键
	return objectKey
}

// GenerateStoragePath 生成存储路径
// pathType: 路径类型，如"wx_avatar", "wx_image", "admin_file"等
// extName: 文件扩展名，包含前导点，如".jpg"
// 返回最终的对象键（不含目录前缀，会通过GetFullObjectKey添加）
func (c *OSSClient) GenerateStoragePath(pathType string, extName string) string {
	var subPath string

	switch pathType {
	case "wx_avatar":
		// 微信客户端用户头像
		subPath = fmt.Sprintf("wx/avatar/%s/%s",
			time.Now().Format("2006/01/02"),
			guid.S())
	case "wx_image":
		// 微信客户端上传图片
		subPath = fmt.Sprintf("wx/uploads/images/%s/%s",
			time.Now().Format("20060102"),
			grand.S(16))
	case "admin_image":
		// 管理后台上传图片
		subPath = fmt.Sprintf("uploads/images/%s/%s",
			time.Now().Format("20060102"),
			grand.S(16))
	case "admin_document":
		// 管理后台上传文档
		subPath = fmt.Sprintf("uploads/documents/%s/%s",
			time.Now().Format("20060102"),
			grand.S(16))
	case "admin_video":
		// 管理后台上传视频
		subPath = fmt.Sprintf("uploads/videos/%s/%s",
			time.Now().Format("20060102"),
			grand.S(16))
	case "admin_audio":
		// 管理后台上传音频
		subPath = fmt.Sprintf("uploads/audios/%s/%s",
			time.Now().Format("20060102"),
			grand.S(16))
	default:
		// 默认路径
		subPath = fmt.Sprintf("uploads/others/%s/%s",
			time.Now().Format("20060102"),
			grand.S(16))
	}

	return subPath + extName
}
