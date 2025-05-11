package v1

import (
	"context"
	"encoding/json"
	"path/filepath"

	"demo/utility/auth"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
)

// wxappConfigModel 微信小程序配置结构
type wxappConfigModel struct {
	AppId     string `json:"appId"`     // 小程序AppID
	AppSecret string `json:"appSecret"` // 小程序AppSecret
	Enabled   bool   `json:"enabled"`   // 是否启用
}

// GetWxappConfig 获取微信小程序配置
func (c *ControllerImpl) GetWxappConfig(ctx context.Context, req *WxappConfigGetReq) (res *WxappConfigGetRes, err error) {
	// 检查是否有管理员权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil || role != "admin" {
		return nil, gerror.New("您没有权限访问此接口")
	}

	// 获取配置
	config, err := loadWxappConfig()
	if err != nil {
		return nil, err
	}

	// 返回配置
	res = &WxappConfigGetRes{
		AppId:     config.AppId,
		AppSecret: config.AppSecret,
		Enabled:   config.Enabled,
	}
	return res, nil
}

// SaveWxappConfig 保存微信小程序配置
func (c *ControllerImpl) SaveWxappConfig(ctx context.Context, req *WxappConfigSaveReq) (res *WxappConfigSaveRes, err error) {
	// 检查是否有管理员权限
	_, _, role, err := auth.GetUserInfo(ctx)
	if err != nil || role != "admin" {
		return nil, gerror.New("您没有权限访问此接口")
	}

	// 准备配置
	config := &wxappConfigModel{
		AppId:     req.AppId,
		AppSecret: req.AppSecret,
		Enabled:   req.Enabled,
	}

	// 保存配置
	err = saveWxappConfig(config)
	if err != nil {
		return nil, err
	}

	// 返回响应
	res = &WxappConfigSaveRes{}
	return res, nil
}

// loadWxappConfig 加载微信小程序配置
func loadWxappConfig() (*wxappConfigModel, error) {
	configPath := "manifest/config/wxapp.json"

	// 检查配置文件是否存在
	if !gfile.Exists(configPath) {
		// 创建默认配置
		defaultConfig := &wxappConfigModel{
			AppId:     "",
			AppSecret: "",
			Enabled:   false,
		}

		// 保存默认配置
		configData, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return nil, gerror.New("创建默认配置失败: " + err.Error())
		}

		err = gfile.PutBytes(configPath, configData)
		if err != nil {
			return nil, gerror.New("保存默认配置失败: " + err.Error())
		}

		return defaultConfig, nil
	}

	// 读取配置文件
	configData := gfile.GetBytes(configPath)
	if len(configData) == 0 {
		return nil, gerror.New("读取微信小程序配置失败: 文件为空")
	}

	// 解析配置
	config := &wxappConfigModel{}
	if err := json.Unmarshal(configData, config); err != nil {
		return nil, gerror.New("解析微信小程序配置失败: " + err.Error())
	}

	return config, nil
}

// saveWxappConfig 保存微信小程序配置
func saveWxappConfig(config *wxappConfigModel) error {
	configPath := "manifest/config/wxapp.json"

	// 确保目录存在
	configDir := filepath.Dir(configPath)
	if !gfile.Exists(configDir) {
		if err := gfile.Mkdir(configDir); err != nil {
			return gerror.New("创建配置目录失败: " + err.Error())
		}
	}

	// 序列化配置
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return gerror.New("序列化配置失败: " + err.Error())
	}

	// 保存配置
	err = gfile.PutBytes(configPath, configData)
	if err != nil {
		return gerror.New("保存配置失败: " + err.Error())
	}

	return nil
}
