package service

import (
	"demo/internal/logic/settings"
)

// 单例模式封装
var localSettings = settings.New()

// Settings 获取城市系统基础设置服务接口实例
func Settings() ISettings {
	return localSettings
}
