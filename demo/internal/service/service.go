package service

import (
	"sync"
)

var (
	// Services 存储所有已注册的服务
	services = make(map[string]interface{})
	// 互斥锁，用于保护 services 映射
	mu sync.RWMutex
)

// 注册商品服务
func RegisterProduct(service IProductService) {
	mu.Lock()
	defer mu.Unlock()
	services["product"] = service
}

// 获取商品服务实例
func Product() IProductService {
	mu.RLock()
	defer mu.RUnlock()
	if svc, ok := services["product"]; ok && svc != nil {
		return svc.(IProductService)
	}
	panic("ProductService implementation not found")
}
