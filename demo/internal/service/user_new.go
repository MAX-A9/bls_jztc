package service

import (
	"sync"
)

var (
	userServiceInstance UserService
	userOnce            sync.Once
)

// User 获取用户服务实例
func User() UserService {
	if userServiceInstance == nil {
		panic("用户服务未初始化")
	}
	return userServiceInstance
}

// SetUser 设置用户服务实例
func SetUser(service UserService) {
	userOnce.Do(func() {
		userServiceInstance = service
	})
}
