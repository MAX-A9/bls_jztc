package storage

import (
	"context"
	"os"

	v1 "demo/api/storage/v1"
	"demo/internal/model"
	"demo/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

// 配置文件路径
const (
	configFilePath = "manifest/config/storage.json"
)

// 存储服务实现
type storageImpl struct{}

// New 创建存储服务实例
func New() service.StorageService {
	return &storageImpl{}
}

// GetConfig 获取存储配置
func (s *storageImpl) GetConfig(ctx context.Context, req *v1.StorageConfigReq) (res *v1.StorageConfigRes, err error) {
	config, err := s.readConfig(ctx)
	if err != nil {
		return nil, err
	}

	// 创建响应对象
	res = &v1.StorageConfigRes{
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: s.maskSecret(config.AccessKeySecret), // 掩码处理密钥
		Endpoint:        config.Endpoint,
		Bucket:          config.Bucket,
		Region:          config.Region,
		Directory:       config.Directory,
		PublicAccess:    config.PublicAccess,
	}

	return res, nil
}

// SaveConfig 保存存储配置
func (s *storageImpl) SaveConfig(ctx context.Context, req *v1.SaveStorageConfigReq) (res *v1.SaveStorageConfigRes, err error) {
	// 获取现有配置
	oldConfig, err := s.readConfig(ctx)
	if err != nil {
		// 如果读取失败但是不是因为文件不存在，则返回错误
		if !os.IsNotExist(err) {
			g.Log().Error(ctx, "读取存储配置失败", err)
			return nil, gerror.New("读取现有配置失败")
		}
	}

	// 创建新的配置对象
	config := &model.StorageConfig{
		AccessKeyId:  req.AccessKeyId,
		Endpoint:     req.Endpoint,
		Bucket:       req.Bucket,
		Region:       req.Region,
		Directory:    req.Directory,
		PublicAccess: req.PublicAccess,
	}

	// 处理AccessKeySecret
	// 如果是掩码（没有修改），则使用旧配置中的值
	if req.AccessKeySecret == "********" && oldConfig != nil {
		config.AccessKeySecret = oldConfig.AccessKeySecret
	} else {
		config.AccessKeySecret = req.AccessKeySecret
	}

	// 创建配置目录（如果不存在）
	if err = gfile.Mkdir(gfile.Dir(configFilePath)); err != nil {
		g.Log().Error(ctx, "创建配置目录失败", err)
		return nil, gerror.New("创建配置目录失败")
	}

	// 将配置对象转换为JSON并保存到文件
	if err = gfile.PutContents(configFilePath, gjson.MustEncodeString(config)); err != nil {
		g.Log().Error(ctx, "保存存储配置失败", err)
		return nil, gerror.New("保存配置失败")
	}

	// 创建响应对象
	res = &v1.SaveStorageConfigRes{
		Success: true,
		Message: "配置保存成功",
	}

	// 记录日志
	g.Log().Info(ctx, "存储配置已更新", gtime.Now())

	return res, nil
}

// readConfig 读取存储配置
func (s *storageImpl) readConfig(ctx context.Context) (config *model.StorageConfig, err error) {
	// 检查配置文件是否存在
	if !gfile.Exists(configFilePath) {
		// 文件不存在时返回默认配置
		return &model.StorageConfig{
			AccessKeyId:     "",
			AccessKeySecret: "",
			Endpoint:        "",
			Bucket:          "",
			Region:          "",
			Directory:       "",
			PublicAccess:    false,
		}, nil
	}

	// 读取配置文件内容
	content := gfile.GetContents(configFilePath)
	if content == "" {
		g.Log().Error(ctx, "读取配置文件失败")
		return nil, gerror.New("读取配置文件失败")
	}

	// 解析JSON内容
	config = &model.StorageConfig{}
	if err = gjson.DecodeTo(content, config); err != nil {
		g.Log().Error(ctx, "解析配置文件失败", err)
		return nil, gerror.New("解析配置文件失败")
	}

	return config, nil
}

// maskSecret 对密钥进行掩码处理
func (s *storageImpl) maskSecret(secret string) string {
	if secret == "" {
		return ""
	}
	return "********"
}
